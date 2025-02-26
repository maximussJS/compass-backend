package models

import (
	"gorm.io/gorm"
	"time"
)

type Category struct {
	Id        uint      `gorm:"primaryKey;autoIncrement" json:"id" `
	Name      string    `gorm:"size:100;not null;unique" json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (c *Category) TableName() string {
	return "categories"
}

func (c *Category) BeforeCreate(_ *gorm.DB) (err error) {
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
	return nil
}

func (c *Category) BeforeUpdate(_ *gorm.DB) (err error) {
	c.UpdatedAt = time.Now()
	return nil
}
