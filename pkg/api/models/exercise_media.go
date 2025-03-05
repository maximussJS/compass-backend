package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type ExerciseMedia struct {
	Id                 string    `gorm:"primaryKey" json:"id"`
	ExerciseId         string    `gorm:"not null;index" json:"exerciseId"`
	Exercise           Exercise  `gorm:"foreignKey:ExerciseId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
	CloudinaryPublicId string    `gorm:"size:255;not null" json:"-"`
	CloudinaryUrl      string    `gorm:"size:1024;not null" json:"url"`
	CreatedAt          time.Time `json:"createdAt"`
	UpdatedAt          time.Time `json:"updatedAt"`
}

func (c *ExerciseMedia) TableName() string {
	return "exercise_media"
}

func (c *ExerciseMedia) BeforeCreate(_ *gorm.DB) (err error) {
	c.Id = uuid.New().String()
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
	return nil
}

func (c *ExerciseMedia) BeforeUpdate(_ *gorm.DB) (err error) {
	c.UpdatedAt = time.Now()
	return nil
}
