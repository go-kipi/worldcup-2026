package service

import (
	"github.com/go-kipi/worldcup-2026/internal/models"
	"gorm.io/gorm"
)

type LeaderboardService struct {
	db *gorm.DB
}

func NewLeaderboardService(db *gorm.DB) *LeaderboardService {
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
	var data LeaderboardData
	if err := s.db.Find(&data.Users).Error; err != nil {
		return nil, err
	}
	if err := s.db.Find(&data.Matches).Error; err != nil {
		return nil, err
	}
	if err := s.db.Find(&data.KnockoutSlots).Error; err != nil {
		return nil, err
	}
	if err := s.db.Find(&data.MatchPredictions).Error; err != nil {
		return nil, err
	}
	if err := s.db.Find(&data.KnockoutPredictions).Error; err != nil {
		return nil, err
	}
	return &data, nil
}
