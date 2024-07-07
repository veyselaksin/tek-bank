package repository

import (
	"gorm.io/gorm"
	"tek-bank/internal/db/models"
)

//go:generate mockgen -destination=../../mocks/repository/transfer_history_repository_mock.go -package=repository tek-bank/internal/db/repository TransferHistoryRepository
type TransferHistoryRepository interface {
	Create(transferHistory []models.TransferHistory) error
	FetchByAccountNumber(accountNumber int64) ([]models.TransferHistory, error)
}

type transferHistoryRepository struct {
	db        *gorm.DB
	tableName string
}

func NewTransferHistoryRepository(db *gorm.DB) TransferHistoryRepository {
	var transferHistory models.TransferHistory
	return &transferHistoryRepository{
		db:        db,
		tableName: transferHistory.TableName(),
	}
}

func (d *transferHistoryRepository) Create(transferHistory []models.TransferHistory) error {
	result := d.db.Table(d.tableName).Create(&transferHistory)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (d *transferHistoryRepository) FetchByAccountNumber(accountNumber int64) ([]models.TransferHistory, error) {
	var transferHistory []models.TransferHistory
	result := d.db.Table(d.tableName).Where("\"from\" = ? OR \"to\" = ?", accountNumber, accountNumber).Find(&transferHistory)
	if result.Error != nil {
		return nil, result.Error
	}
	return transferHistory, nil
}
