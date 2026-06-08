package models

import (
	"time"
)

type User struct {
	ID        string    `gorm:"primaryKey" json:"id" bson:"_id,omitempty"`
	Email     string    `gorm:"uniqueIndex;not null" json:"email" bson:"email"`
	Name      string    `json:"name" bson:"name"`
	Role      string    `gorm:"default:'player'" json:"role" bson:"role"` // 'player' or 'admin'
	Points    int       `gorm:"default:0" json:"points" bson:"points"`    // cached leaderboard score
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

type OTP struct {
	ID        string    `gorm:"primaryKey" bson:"_id,omitempty"`
	Email     string    `gorm:"index" bson:"email"`
	Code      string    `bson:"code"`
	ExpiresAt time.Time `bson:"expires_at"`
	CreatedAt time.Time `bson:"created_at"`
}

type Team struct {
	ID        string `gorm:"primaryKey;type:varchar(5)" json:"id" bson:"_id,omitempty"` // 3-letter code, e.g. 'BRA'
	Name      string `json:"name" bson:"name"`
	GroupCode string `json:"group_code" bson:"group_code"`
	FlagEmoji string `json:"flag_emoji" bson:"flag_emoji"`
}

type Match struct {
	ID          string    `gorm:"primaryKey" json:"id" bson:"_id,omitempty"`
	GroupCode   string    `json:"group_code" bson:"group_code"`
	MatchNumber *int      `json:"match_number" bson:"match_number"`
	Kickoff     time.Time `json:"kickoff" bson:"kickoff"`
	HomeTeam    string    `json:"home_team" bson:"home_team"`
	AwayTeam    string    `json:"away_team" bson:"away_team"`
	HomeScore   int       `json:"home_score" bson:"home_score"`
	AwayScore   int       `json:"away_score" bson:"away_score"`
	Finished    bool      `gorm:"default:false" json:"finished" bson:"finished"`
}

type KnockoutSlot struct {
	ID             string     `gorm:"primaryKey" json:"id" bson:"_id,omitempty"`
	Round          string     `json:"round" bson:"round"`
	SlotCode       string     `gorm:"unique" json:"slot_code" bson:"slot_code"`
	Label          *string    `json:"label" bson:"label"`
	Kickoff        *time.Time `json:"kickoff" bson:"kickoff"`
	HomeLabel      *string    `json:"home_label" bson:"home_label"`
	AwayLabel      *string    `json:"away_label" bson:"away_label"`
	ActualHomeTeam *string    `json:"actual_home_team" bson:"actual_home_team"`
	ActualAwayTeam *string    `json:"actual_away_team" bson:"actual_away_team"`
	HomeScore      int        `json:"home_score" bson:"home_score"`
	AwayScore      int        `json:"away_score" bson:"away_score"`
	Finished       bool       `gorm:"default:false" json:"finished" bson:"finished"`
}

type MatchPrediction struct {
	ID        string    `gorm:"primaryKey" json:"-" bson:"_id,omitempty"`
	UserID    string    `gorm:"index;uniqueIndex:idx_user_match" json:"user_id" bson:"user_id"`
	MatchID   string    `gorm:"index;uniqueIndex:idx_user_match" json:"match_id" bson:"match_id"`
	HomeScore int       `json:"home_score" bson:"home_score"`
	AwayScore int       `json:"away_score" bson:"away_score"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`

	User  User  `gorm:"foreignKey:UserID" json:"-" bson:"-"`
	Match Match `gorm:"foreignKey:MatchID" json:"-" bson:"-"`
}

type KnockoutPrediction struct {
	ID        string    `gorm:"primaryKey" json:"-" bson:"_id,omitempty"`
	UserID    string    `gorm:"index;uniqueIndex:idx_user_slot" json:"user_id" bson:"user_id"`
	SlotID    string    `gorm:"index;uniqueIndex:idx_user_slot" json:"slot_id" bson:"slot_id"`
	TeamID    *string   `json:"team_id" bson:"team_id"` // Legacy, now optional
	HomeScore int       `json:"home_score" bson:"home_score"`
	AwayScore int       `json:"away_score" bson:"away_score"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`

	User User         `gorm:"foreignKey:UserID" json:"-" bson:"-"`
	Slot KnockoutSlot `gorm:"foreignKey:SlotID" json:"-" bson:"-"`
}

type AppSetting struct {
	ID     string    `gorm:"primaryKey;default:1" json:"id" bson:"_id,omitempty"`
	LockAt time.Time `json:"lock_at" bson:"lock_at"`
}
