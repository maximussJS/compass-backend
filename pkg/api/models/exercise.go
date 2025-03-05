package models

import (
	common_models "compass-backend/pkg/common/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Exercise struct {
	Id          string             `gorm:"primaryKey" json:"id"`
	Name        string             `gorm:"size:100;not null;unique" json:"name"`
	Description string             `gorm:"size:255;not null" json:"description"`
	CreatorId   string             `gorm:"primaryKey" json:"creatorId"`
	Creator     common_models.User `gorm:"foreignKey:CreatorId;references:Id" json:"creator"`
	Files       []ExerciseMedia    `gorm:"foreignKey:ExerciseId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"mediaFiles"`
	CreatedAt   time.Time          `json:"createdAt"`
	UpdatedAt   time.Time          `json:"updatedAt"`
}

func (c *Exercise) TableName() string {
	return "exercises"
}

func (c *Exercise) BeforeCreate(_ *gorm.DB) (err error) {
	c.Id = uuid.New().String()
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
	return nil
}

func (c *Exercise) BeforeUpdate(_ *gorm.DB) (err error) {
	c.UpdatedAt = time.Now()
	return nil
}
