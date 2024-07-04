package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Account struct {
	Id            string `gorm:"primary_key;type:uuid;"`
	AccountNumber string `gorm:"unique;not null"`
	OwnerID       string `gorm:"type:uuid;not null"`
	Balance       float64

	// Audit fields
	CreatedAt time.Time `gorm:"default:current_timestamp"`
	UpdatedAt time.Time `gorm:"default:current_timestamp"`
	CreatedBy string    `gorm:"type:uuid"`
	UpdatedBy string    `gorm:"type:uuid"`
	IsActive  bool      `gorm:"default:true"`

	// Relationship
	Owner User `gorm:"foreignKey:OwnerID"`
}

func (a *Account) BeforeCreate(tx *gorm.DB) error {
	a.Id = uuid.New().String()
	return nil
}

func (a *Account) TableName() string {
	return "public.accounts"
}
