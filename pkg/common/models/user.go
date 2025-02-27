package models

import (
	"compass-backend/pkg/common/constants"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type User struct {
	Id         string             `gorm:"primaryKey" json:"id"`
	Name       string             `gorm:"size:100;not null" json:"name"`
	Email      string             `gorm:"size:100;not null;unique" json:"email"`
	Password   string             `gorm:"size:255" json:"-"`
	IsVerified bool               `gorm:"default:false" json:"isVerified"`
	Role       constants.UserRole `gorm:"size:100;not null" json:"role"`
	Teams      []Team             `gorm:"many2many:team_members;joinForeignKey:UserId;joinReferences:TeamId" json:"teams,omitempty"`
	CreatedAt  time.Time          `json:"createdAt"`
	UpdatedAt  time.Time          `json:"updatedAt"`
}

func (u *User) TableName() string {
	return "users"
}

func (u *User) BeforeCreate(_ *gorm.DB) (err error) {
	u.Id = uuid.New().String()
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	return nil
}

func (u *User) BeforeUpdate(_ *gorm.DB) (err error) {
	u.UpdatedAt = time.Now()
	return nil
}

func (u *User) IsClient() bool {
	return u.Role == constants.Client
}

func (u *User) IsTrainer() bool {
	return u.Role == constants.Trainer
}

func (u *User) IsInTeam(teamId string) bool {
	for _, team := range u.Teams {
		if team.Id == teamId {
			return true
		}
	}
	return false
}
