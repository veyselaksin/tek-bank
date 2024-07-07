package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type User struct {
	Id             string `gorm:"primary_key;type:uuid;"`
	IdentityNumber int64  `gorm:"unique;not null"`
	CustomerNumber int64  `gorm:"unique;not null"`
	FirstName      string `gorm:"not null"`
	LastName       string `gorm:"not null"`
	Email          string `gorm:"unique;not null"`
	PhoneNumber    uint64 `gorm:"unique;not null"`
	Password       string `gorm:"not null"`

	// Audit fields
	CreatedAt time.Time `gorm:"default:current_timestamp"`
	UpdatedAt time.Time `gorm:"default:current_timestamp"`
	IsActive  bool      `gorm:"default:true"`

	// Relationship
	Accounts []Account `gorm:"foreignKey:OwnerId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	u.Id = uuid.New().String()
	return nil
}

func (u *User) TableName() string {
	return "public.users"
}
