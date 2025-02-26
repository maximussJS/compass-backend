package models

import (
	"gorm.io/gorm"
	"time"
)

type TeamMember struct {
	TeamId    string    `gorm:"primaryKey" json:"team_id"`
	UserId    string    `gorm:"primaryKey" json:"user_id"`
	Team      Team      `gorm:"foreignKey:TeamId;references:Id" json:"team"`
	User      User      `gorm:"foreignKey:UserId;references:Id" json:"user"`
	CreatedAt time.Time `json:"createdAt"`
}

func (tm *TeamMember) TableName() string {
	return "team_members"
}

func (tm *TeamMember) BeforeCreate(_ *gorm.DB) (err error) {
	tm.CreatedAt = time.Now()
	return nil
}
