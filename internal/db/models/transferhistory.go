package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type TransferHistory struct {
	Id     string `gorm:"primary_key;type:uuid;"`
	From   string `gorm:"type:uuid;not null"`
	To     string `gorm:"type:uuid;not null"`
	Status string `gorm:"default:'pending'"`
	Note   string `gorm:"default:null"`
	Amount float64

	// Audit fields
	CreatedAt time.Time `gorm:"default:current_timestamp"`
	UpdatedAt time.Time `gorm:"default:current_timestamp"`
	CreatedBy string    `gorm:"type:uuid"`
	UpdatedBy string    `gorm:"type:uuid"`
	IsActive  bool      `gorm:"default:true"`

	// Relationship
	FromAccount Account `gorm:"foreignKey:From"`
	ToAccount   Account `gorm:"foreignKey:To"`
}

func (t *TransferHistory) BeforeCreate(tx *gorm.DB) error {
	t.Id = uuid.New().String()
	return nil
}

func (t *TransferHistory) TableName() string {
	return "public.transfer_history"
}
