package models

import (
	"compass-backend/pkg/common/constants"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type TeamInvite struct {
	Id        string                     `gorm:"primaryKey" json:"id"`
	TeamId    string                     `gorm:"size:100;not null" json:"teamId"`
	Team      *Team                      `gorm:"foreignKey:TeamId" json:"team"`
	Status    constants.TeamInviteStatus `gorm:"size:100;not null" json:"status"`
	Token     string                     `gorm:"size:400;not null" json:"token"`
	Email     string                     `gorm:"size:100;not null" json:"email"`
	IsSent    bool                       `gorm:"default:false" json:"isSent"`
	ExpiresAt time.Time                  `json:"expiresAt"`
	CreatedAt time.Time                  `json:"createdAt"`
}

func (ti *TeamInvite) TableName() string {
	return "team_invites"
}

func (ti *TeamInvite) BeforeCreate(_ *gorm.DB) (err error) {
	ti.Id = uuid.New().String()
	ti.Status = constants.TeamInviteStatusPending
	ti.CreatedAt = time.Now()
	return nil
}

func (ti *TeamInvite) IsExpired() bool {
	return time.Now().After(ti.ExpiresAt)
}

func (ti *TeamInvite) IsAccepted() bool {
	return ti.Status == constants.TeamInviteStatusAccepted
}

func (ti *TeamInvite) IsPending() bool {
	return ti.Status == constants.TeamInviteStatusPending
}

func (ti *TeamInvite) IsCancelled() bool {
	return ti.Status == constants.TeamInviteStatusCancelled
}

func (ti *TeamInvite) TeamName() string {
	if ti.Team.Id == "" {
		panic("Team is not loaded for team invite")
	}
	return ti.Team.Name
}

func (ti *TeamInvite) TeamOwnerName() string {
	if ti.Team.Id == "" {
		panic("Team is not loaded for team invite")
	}

	if ti.Team.Owner.Id == "" {
		panic("Team owner is not loaded for team invite")
	}

	return ti.Team.Owner.Name
}

func (ti *TeamInvite) TeamOwnerEmail() string {
	if ti.Team.Id == "" {
		panic("Team is not loaded for team invite")
	}

	if ti.Team.Owner.Id == "" {
		panic("Team owner is not loaded for team invite")
	}

	return ti.Team.Owner.Email
}
