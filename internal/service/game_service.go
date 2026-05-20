package service

import (
	"github.com/go-kipi/worldcup-2026/internal/models"
	"gorm.io/gorm"
)

type GameService struct {
	db *gorm.DB
}

func NewGameService(db *gorm.DB) *GameService {
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

func (s *GameService) GetAppState(userID uint) (*AppState, error) {
	var state AppState

	if err := s.db.Find(&state.Teams).Error; err != nil {
		return nil, err
	}
	if err := s.db.Order("kickoff ASC").Find(&state.Matches).Error; err != nil {
		return nil, err
	}
	if err := s.db.Order("id ASC").Find(&state.KnockoutSlots).Error; err != nil {
		return nil, err
	}
	if err := s.db.First(&state.Settings, 1).Error; err != nil {
		return nil, err
	}

	if userID > 0 {
		s.db.Where("user_id = ?", userID).Find(&state.MatchPredictions)
		s.db.Where("user_id = ?", userID).Find(&state.KnockoutPredictions)
	}

	return &state, nil
}

func (s *GameService) UpdateMatch(matchID uint, homeScore *int, awayScore *int, finished bool) error {
	return s.db.Model(&models.Match{}).Where("id = ?", matchID).Updates(map[string]interface{}{
		"home_score": homeScore,
		"away_score": awayScore,
		"finished":   finished,
	}).Error
}

func (s *GameService) UpdateKnockoutSlot(slotID uint, patch map[string]interface{}) error {
	return s.db.Model(&models.KnockoutSlot{}).Where("id = ?", slotID).Updates(patch).Error
}

func (s *GameService) UpdateSettings(lockAt string) error {
	return s.db.Model(&models.AppSetting{}).Where("id = 1").Update("lock_at", lockAt).Error
}
