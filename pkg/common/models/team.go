package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Team struct {
	Id        string    `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:100;not null" json:"name"`
	OwnerId   string    `gorm:"size:100;not null" json:"-"`
	Owner     User      `gorm:"foreignKey:OwnerId;references:Id" json:"owner"`
	Members   []User    `gorm:"many2many:team_members;joinForeignKey:TeamId;joinReferences:UserId" json:"members"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (t *Team) TableName() string {
	return "teams"
}

func (t *Team) BeforeCreate(_ *gorm.DB) (err error) {
	t.Id = uuid.New().String()
	t.CreatedAt = time.Now()
	t.UpdatedAt = time.Now()
	return nil
}

func (t *Team) BeforeUpdate(_ *gorm.DB) (err error) {
	t.UpdatedAt = time.Now()
	return nil
}

func (t *Team) IsOwner(user *User) bool {
	return t.OwnerId == user.Id
}
