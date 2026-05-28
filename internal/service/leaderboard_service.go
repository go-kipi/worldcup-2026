package service

import (
	"context"
	"github.com/go-kipi/worldcup-2026/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type LeaderboardService struct {
	db *mongo.Database
}

func NewLeaderboardService(db *mongo.Database) *LeaderboardService {
	return &LeaderboardService{db: db}
}

type LeaderboardData struct {
	Users               []models.User               `json:"users"`
	Matches             []models.Match              `json:"matches"`
	KnockoutSlots       []models.KnockoutSlot       `json:"knockout_slots"`
	MatchPredictions    []models.MatchPrediction    `json:"match_predictions"`
	KnockoutPredictions []models.KnockoutPrediction `json:"knockout_predictions"`
}

func (s *LeaderboardService) GetLeaderboard() (*LeaderboardData, error) {
	ctx := context.Background()
	var data LeaderboardData

	cursor, err := s.db.Collection("users").Find(ctx, bson.M{})
	if err == nil {
		cursor.All(ctx, &data.Users)
	}

	cursor, err = s.db.Collection("matches").Find(ctx, bson.M{})
	if err == nil {
		cursor.All(ctx, &data.Matches)
	}

	cursor, err = s.db.Collection("knockout_slots").Find(ctx, bson.M{})
	if err == nil {
		cursor.All(ctx, &data.KnockoutSlots)
	}

	cursor, err = s.db.Collection("match_predictions").Find(ctx, bson.M{})
	if err == nil {
		cursor.All(ctx, &data.MatchPredictions)
	}

	cursor, err = s.db.Collection("knockout_predictions").Find(ctx, bson.M{})
	if err == nil {
		cursor.All(ctx, &data.KnockoutPredictions)
	}

	return &data, nil
}
