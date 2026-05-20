package service

import (
	"errors"
	"time"

	"github.com/go-kipi/worldcup-2026/internal/models"
	"gorm.io/gorm"
)

type PredictionService struct {
	db *gorm.DB
}

func NewPredictionService(db *gorm.DB) *PredictionService {
	return &PredictionService{db: db}
}

func (s *PredictionService) GetUserPredictions(userID uint) ([]models.MatchPrediction, []models.KnockoutPrediction, error) {
	var matchPreds []models.MatchPrediction
	if err := s.db.Where("user_id = ?", userID).Find(&matchPreds).Error; err != nil {
		return nil, nil, err
	}

	var koPreds []models.KnockoutPrediction
	if err := s.db.Where("user_id = ?", userID).Find(&koPreds).Error; err != nil {
		return nil, nil, err
	}

	return matchPreds, koPreds, nil
}

type MatchPredictionInput struct {
	MatchID   uint `json:"match_id" binding:"required"`
	HomeScore int  `json:"home_score"`
	AwayScore int  `json:"away_score"`
}

type KnockoutPredictionInput struct {
	SlotID    uint `json:"slot_id" binding:"required"`
	HomeScore int  `json:"home_score"`
	AwayScore int  `json:"away_score"`
}

func (s *PredictionService) SavePredictions(userID uint, matchPreds []MatchPredictionInput, koPreds []KnockoutPredictionInput) error {
	var settings models.AppSetting
	if err := s.db.First(&settings, 1).Error; err != nil {
		return errors.New("system settings not found")
	}

	if time.Now().After(settings.LockAt) {
		return errors.New("predictions are closed (past cutoff)")
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		// Insert new predictions
		for _, mp := range matchPreds {
			pred := models.MatchPrediction{
				UserID:    userID,
				MatchID:   mp.MatchID,
				HomeScore: mp.HomeScore,
				AwayScore: mp.AwayScore,
			}
			// Upsert logic for MatchPrediction
			if err := tx.Where(models.MatchPrediction{UserID: userID, MatchID: mp.MatchID}).Assign(pred).FirstOrCreate(&pred).Error; err != nil {
				return err
			}
		}

		for _, kp := range koPreds {
			h, a := kp.HomeScore, kp.AwayScore
			pred := models.KnockoutPrediction{
				UserID:    userID,
				SlotID:    kp.SlotID,
				HomeScore: &h,
				AwayScore: &a,
			}
			// Upsert logic for KnockoutPrediction
			if err := tx.Where(models.KnockoutPrediction{UserID: userID, SlotID: kp.SlotID}).Assign(pred).FirstOrCreate(&pred).Error; err != nil {
				return err
			}
		}

		return nil
	})
}
