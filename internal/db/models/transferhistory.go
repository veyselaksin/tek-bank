package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type TransferHistory struct {
	Id     string  `gorm:"primary_key;type:uuid;"`
	From   int64   `gorm:"type:bigint;not null"`
	To     int64   `gorm:"type:bigint;not null"`
	Note   string  `gorm:"default:null"`
	Amount float64 `gorm:"type:numeric;not null"`
	IsFee  bool    `gorm:"default:false"`

	// Audit fields
	CreatedAt time.Time `gorm:"default:current_timestamp"`
	UpdatedAt time.Time `gorm:"default:current_timestamp"`
	CreatedBy string    `gorm:"type:uuid"`
	UpdatedBy string    `gorm:"type:uuid"`
	IsActive  bool      `gorm:"default:true"`

	// Relationship
	FromAccount Account `gorm:"foreignKey:From;references:AccountNumber;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ToAccount   Account `gorm:"foreignKey:To;references:AccountNumber;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (t *TransferHistory) BeforeCreate(tx *gorm.DB) error {
	t.Id = uuid.New().String()
	return nil
}

func (t *TransferHistory) TableName() string {
	return "public.transfer_history"
}
