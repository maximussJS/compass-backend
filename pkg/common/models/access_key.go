package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type AccessKey struct {
	Id        string     `gorm:"primaryKey" json:"id"`
	Name      string     `gorm:"size:100;not null;unique" json:"name"`
	ExpiresAt *time.Time `json:"expiresAt"`
	UserId    string     `gorm:"not null" json:"userId"`
	User      User       `gorm:"foreignKey:UserId;references:Id" json:"user"`
	TeamId    string     `gorm:"primaryKey" json:"teamId"`
	Team      Team       `gorm:"foreignKey:TeamId;references:Id" json:"team"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
}

func (c *AccessKey) TableName() string {
	return "access_keys"
}

func (c *AccessKey) BeforeCreate(_ *gorm.DB) (err error) {
	c.Id = uuid.New().String()
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
