package models

import (
	"time"
)

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Email     string    `gorm:"uniqueIndex;not null" json:"email"`
	Name      string    `json:"name"`
	Role      string    `gorm:"default:'player'" json:"role"` // 'player' or 'admin'
	Points    int       `gorm:"default:0" json:"points"`      // cached leaderboard score
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type OTP struct {
	ID        uint   `gorm:"primaryKey"`
	Email     string `gorm:"index"`
	Code      string
	ExpiresAt time.Time
	CreatedAt time.Time
}

type Team struct {
	ID        string `gorm:"primaryKey;type:varchar(5)" json:"id"` // 3-letter code, e.g. 'BRA'
	Name      string `json:"name"`
	GroupCode string `json:"group_code"`
	FlagEmoji string `json:"flag_emoji"`
}

type Match struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	GroupCode   string    `json:"group_code"`
	MatchNumber *int      `json:"match_number"`
	Kickoff     time.Time `json:"kickoff"`
	HomeTeam    string    `json:"home_team"`
	AwayTeam    string    `json:"away_team"`
	HomeScore   *int      `json:"home_score"`
	AwayScore   *int      `json:"away_score"`
	Finished    bool      `gorm:"default:false" json:"finished"`
}

type KnockoutSlot struct {
	ID             uint       `gorm:"primaryKey" json:"id"`
	Round          string     `json:"round"`
	SlotCode       string     `gorm:"unique" json:"slot_code"`
	Label          *string    `json:"label"`
	Kickoff        *time.Time `json:"kickoff"`
	HomeLabel      *string    `json:"home_label"`
	AwayLabel      *string    `json:"away_label"`
	ActualHomeTeam *string    `json:"actual_home_team"`
	ActualAwayTeam *string    `json:"actual_away_team"`
	HomeScore      *int       `json:"home_score"`
	AwayScore      *int       `json:"away_score"`
}

type MatchPrediction struct {
	ID        uint      `gorm:"primaryKey" json:"-"`
	UserID    uint      `gorm:"index;uniqueIndex:idx_user_match" json:"user_id"`
	MatchID   uint      `gorm:"index;uniqueIndex:idx_user_match" json:"match_id"`
	HomeScore int       `json:"home_score"`
	AwayScore int       `json:"away_score"`
	UpdatedAt time.Time `json:"updated_at"`

	User  User  `gorm:"foreignKey:UserID" json:"-"`
	Match Match `gorm:"foreignKey:MatchID" json:"-"`
}

type KnockoutPrediction struct {
	ID        uint      `gorm:"primaryKey" json:"-"`
	UserID    uint      `gorm:"index;uniqueIndex:idx_user_slot" json:"user_id"`
	SlotID    uint      `gorm:"index;uniqueIndex:idx_user_slot" json:"slot_id"`
	TeamID    *string   `json:"team_id"` // Legacy, now optional
	HomeScore *int      `json:"home_score"`
	AwayScore *int      `json:"away_score"`
	UpdatedAt time.Time `json:"updated_at"`

	User User         `gorm:"foreignKey:UserID" json:"-"`
	Slot KnockoutSlot `gorm:"foreignKey:SlotID" json:"-"`
}

type AppSetting struct {
	ID     uint      `gorm:"primaryKey;default:1" json:"id"`
	LockAt time.Time `json:"lock_at"`
}
