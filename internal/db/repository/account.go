package repository

import (
	"context"
	"github.com/gofiber/fiber/v2/log"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"sync"
	"tek-bank/internal/db/models"
)

type AccountRepository interface {
	Create(account models.Account) (*models.Account, error)
	UpdateBalance(amount float64, id string) error
	FindByAccountNumber(accountNumber int64) (*models.Account, error)
	FindByIBAN(iban string) (*models.Account, error)
	FindByOwnerId(ownerId string) ([]models.Account, error)

	// Redis operations
	SetToken(ctx context.Context, key string, value string) error
	GetToken(ctx context.Context, key string) (*string, error)
	DeleteToken(ctx context.Context, key string) error

	WithTx(trxHandle *gorm.DB) AccountRepository
}

//go:generate mockgen -destination=../../mocks/repository/account_repository_mock.go -package=repository tek-bank/internal/db/repository AccountRepository
type accountRepository struct {
	db          *gorm.DB
	redisClient *redis.Client
	tableName   string
	dbMutex     sync.Mutex
}

func NewAccountRepository(db *gorm.DB, client *redis.Client) AccountRepository {
	var account models.Account
	return &accountRepository{
		db:          db,
		redisClient: client,
		tableName:   account.TableName(),
	}
}

func (d *accountRepository) WithTx(txHandle *gorm.DB) AccountRepository {
	if txHandle == nil {
		log.Error("Transaction not found")
		return d
	}
	d.db = txHandle
	return d
}

func (r *accountRepository) Create(account models.Account) (*models.Account, error) {
	result := r.db.Table(r.tableName).Preload("Owner").Create(&account).First(&account)
	if result.Error != nil {
		return nil, result.Error
	}
	return &account, nil
}

func (r *accountRepository) UpdateBalance(amount float64, id string) error {
	r.dbMutex.Lock()
	defer r.dbMutex.Unlock()

	// Add money to the account
	result := r.db.Table(r.tableName).Preload("Owner").Where("id = ?", id).Update("balance", amount)
	if result.Error != nil {
		return result.Error
	}
	return nil

}

func (r *accountRepository) FindByAccountNumber(accountNumber int64) (*models.Account, error) {
	var account models.Account
	// CustomerNumber is in the Owner table, so we need to preload the Owner relationship
	// after the preload, we can use the Owner.CustomerNumber to filter the account
	result := r.db.Table(r.tableName).Preload("Owner").Where("\"account_number\" = ?", accountNumber).First(&account)
	if result.Error != nil {
		return nil, result.Error
	}
	return &account, nil
}

func (r *accountRepository) FindByIBAN(iban string) (*models.Account, error) {
	var account models.Account
	result := r.db.Table(r.tableName).Preload("Owner").Where("iban = ?", iban).First(&account)
	if result.Error != nil {
		return nil, result.Error
	}
	return &account, nil
}

func (r *accountRepository) FindByOwnerId(ownerId string) ([]models.Account, error) {
	var accounts []models.Account
	result := r.db.Table(r.tableName).Preload("Owner").Where("owner_id = ?", ownerId).Find(&accounts)
	if result.Error != nil {
		return nil, result.Error
	}
	return accounts, nil
}

func (r *accountRepository) SetToken(ctx context.Context, key string, value string) error {

	result := r.redisClient.Set(ctx, key, value, 0)
	if result.Err() != nil {
		return result.Err()
	}
	return nil
}

func (r *accountRepository) GetToken(ctx context.Context, key string) (*string, error) {
	result := r.redisClient.Get(ctx, key)
	if result.Err() != nil {
		return nil, result.Err()
	}
	val := result.Val()
	return &val, nil
}

func (r *accountRepository) DeleteToken(ctx context.Context, key string) error {
	result := r.redisClient.Del(ctx, key)
	if result.Err() != nil {
		return result.Err()
	}
	return nil
}
