package service

import (
	"context"
	"errors"
	"time"

	"github.com/go-kipi/worldcup-2026/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type PredictionService struct {
	db *mongo.Database
}

func NewPredictionService(db *mongo.Database) *PredictionService {
	return &PredictionService{db: db}
}

func (s *PredictionService) GetUserPredictions(userID string) ([]models.MatchPrediction, []models.KnockoutPrediction, error) {
	ctx := context.Background()
	var matchPreds []models.MatchPrediction
	cursor, err := s.db.Collection("match_predictions").Find(ctx, bson.M{"user_id": userID})
	if err == nil {
		cursor.All(ctx, &matchPreds)
	}

	var koPreds []models.KnockoutPrediction
	cursor, err = s.db.Collection("knockout_predictions").Find(ctx, bson.M{"user_id": userID})
	if err == nil {
		cursor.All(ctx, &koPreds)
	}

	return matchPreds, koPreds, nil
}

type MatchPredictionInput struct {
	MatchID   string `json:"match_id" binding:"required"`
	HomeScore int    `json:"home_score"`
	AwayScore int    `json:"away_score"`
}

type KnockoutPredictionInput struct {
	SlotID    string `json:"slot_id" binding:"required"`
	HomeScore int    `json:"home_score"`
	AwayScore int    `json:"away_score"`
}

func (s *PredictionService) SavePredictions(userID string, matchPreds []MatchPredictionInput, koPreds []KnockoutPredictionInput) error {
	ctx := context.Background()
	var settings models.AppSetting
	err := s.db.Collection("app_settings").FindOne(ctx, bson.M{"_id": "1"}).Decode(&settings)
	if err != nil {
		return errors.New("system settings not found")
	}

	if time.Now().After(settings.LockAt) {
		return errors.New("predictions are closed (past cutoff)")
	}

	// upsert match predictions
	for _, mp := range matchPreds {
		filter := bson.M{"user_id": userID, "match_id": mp.MatchID}
		update := bson.M{
			"$set": bson.M{
				"home_score": mp.HomeScore,
				"away_score": mp.AwayScore,
				"updated_at": time.Now(),
			},
		}
		opts := options.UpdateOne().SetUpsert(true)
		_, err := s.db.Collection("match_predictions").UpdateOne(ctx, filter, update, opts)
		if err != nil {
			return err
		}
	}

	// upsert knockout predictions
	for _, kp := range koPreds {
		filter := bson.M{"user_id": userID, "slot_id": kp.SlotID}
		update := bson.M{
			"$set": bson.M{
				"home_score": kp.HomeScore,
				"away_score": kp.AwayScore,
				"updated_at": time.Now(),
			},
		}
		opts := options.UpdateOne().SetUpsert(true)
		_, err := s.db.Collection("knockout_predictions").UpdateOne(ctx, filter, update, opts)
		if err != nil {
			return err
		}
	}

	return nil
}
