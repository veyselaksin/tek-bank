package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type User struct {
	Id        string `gorm:"primary_key;type:uuid;"`
	FirstName string `gorm:"not null"`
	LastName  string `gorm:"not null"`
	Username  string `gorm:"unique;not null"`
	Email     string `gorm:"unique;not null"`
	Password  string `gorm:"not null"`

	// Audit fields
	CreatedAt time.Time `gorm:"default:current_timestamp"`
	UpdatedAt time.Time `gorm:"default:current_timestamp"`
	IsActive  bool      `gorm:"default:true"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	u.Id = uuid.New().String()
	return nil
}

func (u *User) TableName() string {
	return "public.users"
}
