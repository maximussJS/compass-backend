package models

import (
	common_models "compass-backend/pkg/common/models"
	"gorm.io/gorm"
	"time"
)

type AccessKey struct {
	Id        uint               `gorm:"primaryKey;autoIncrement" json:"id" `
	Name      string             `gorm:"size:100;not null;unique" json:"name"`
	ExpiresAt *time.Time         `json:"expiresAt"`
	UserId    string             `gorm:"not null" json:"userId"`
	User      common_models.User `gorm:"foreignKey:UserId;references:Id" json:"user"`
	TeamId    string             `gorm:"primaryKey" json:"teamId"`
	Team      common_models.Team `gorm:"foreignKey:TeamId;references:Id" json:"team"`
	CreatedAt time.Time          `json:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt"`
}

func (c *AccessKey) TableName() string {
	return "access_keys"
}

func (c *AccessKey) BeforeCreate(_ *gorm.DB) (err error) {
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
	c.ExpiresAt = nil
	return nil
}

func (c *AccessKey) BeforeUpdate(_ *gorm.DB) (err error) {
	c.UpdatedAt = time.Now()
	return nil
}

func (c *AccessKey) IsExpired() bool {
	return c.ExpiresAt != nil && c.ExpiresAt.Before(time.Now())
}
