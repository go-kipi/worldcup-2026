package db

import (
	"log"

	"github.com/go-kipi/worldcup-2026/internal/config"
	"github.com/go-kipi/worldcup-2026/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewDatabase(cfg *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(cfg.DBPath), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Migrate the schema
	err = db.AutoMigrate(
		&models.User{},
		&models.OTP{},
		&models.Team{},
		&models.Match{},
		&models.KnockoutSlot{},
		&models.MatchPrediction{},
		&models.KnockoutPrediction{},
		&models.AppSetting{},
	)
	if err != nil {
		log.Printf("Failed to migrate database: %v", err)
		return nil, err
	}

	return db, nil
}
