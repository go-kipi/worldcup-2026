package service

import (
	"context"
	"github.com/go-kipi/worldcup-2026/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type GameService struct {
	db *mongo.Database
}

func NewGameService(db *mongo.Database) *GameService {
	return &GameService{db: db}
}

type AppState struct {
	Teams               []models.Team               `json:"teams"`
	Matches             []models.Match              `json:"matches"`
	KnockoutSlots       []models.KnockoutSlot       `json:"knockout_slots"`
	Settings            models.AppSetting           `json:"settings"`
	MatchPredictions    []models.MatchPrediction    `json:"match_predictions"`
	KnockoutPredictions []models.KnockoutPrediction `json:"knockout_predictions"`
}

func (s *GameService) GetAppState(userID string) (*AppState, error) {
	ctx := context.Background()
	var state AppState

	// Teams
	cursor, err := s.db.Collection("teams").Find(ctx, bson.M{})
	if err == nil {
		cursor.All(ctx, &state.Teams)
	}

	// Matches
	opts := options.Find().SetSort(bson.M{"kickoff": 1})
	cursor, err = s.db.Collection("matches").Find(ctx, bson.M{}, opts)
	if err == nil {
		cursor.All(ctx, &state.Matches)
	}

	// KnockoutSlots
	opts = options.Find().SetSort(bson.M{"_id": 1})
	cursor, err = s.db.Collection("knockout_slots").Find(ctx, bson.M{}, opts)
	if err == nil {
		cursor.All(ctx, &state.KnockoutSlots)
	}

	// Settings
	s.db.Collection("app_settings").FindOne(ctx, bson.M{"_id": "1"}).Decode(&state.Settings)

	// Predictions
	if userID != "" {
		cursor, err = s.db.Collection("match_predictions").Find(ctx, bson.M{"user_id": userID})
		if err == nil {
			cursor.All(ctx, &state.MatchPredictions)
		}

		cursor, err = s.db.Collection("knockout_predictions").Find(ctx, bson.M{"user_id": userID})
		if err == nil {
			cursor.All(ctx, &state.KnockoutPredictions)
		}
	}

	return &state, nil
}

func (s *GameService) UpdateMatch(matchID string, homeScore int, awayScore int, finished bool) error {
	_, err := s.db.Collection("matches").UpdateOne(context.Background(), bson.M{"_id": matchID}, bson.M{
		"$set": bson.M{
			"home_score": homeScore,
			"away_score": awayScore,
			"finished":   finished,
		},
	})
	return err
}

func (s *GameService) UpdateKnockoutSlot(slotID string, patch map[string]interface{}) error {
	_, err := s.db.Collection("knockout_slots").UpdateOne(context.Background(), bson.M{"_id": slotID}, bson.M{
		"$set": patch,
	})
	return err
}

func (s *GameService) UpdateSettings(lockAt string) error {
	_, err := s.db.Collection("app_settings").UpdateOne(context.Background(), bson.M{"_id": "1"}, bson.M{
		"$set": bson.M{"lock_at": lockAt},
	})
	return err
}
