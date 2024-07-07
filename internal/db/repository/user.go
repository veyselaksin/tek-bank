package repository

import (
	"context"
	"github.com/gofiber/fiber/v2/log"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"tek-bank/internal/db/models"
	"time"
)

//go:generate mockgen -destination=../../mocks/repository/user_repository_mock.go -package=repository tek-bank/internal/db/repository UserRepository
type UserRepository interface {
	FindAll() ([]models.User, error)
	FindByID(id string) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	FindByUniqueIdentifier(uniqueIdentifier string) (*models.User, error)
	Create(user models.User) (*models.User, error)
	SoftDelete(id string) error

	SetTokenBlacklist(ctx *context.Context, key string, value string, exp time.Duration) error
	GetTokenBlacklist(ctx *context.Context, key string) (string, error)

	WithTx(trxHandle *gorm.DB) UserRepository
}

type userRepository struct {
	db          *gorm.DB
	redisClient *redis.Client
	tableName   string
}

func NewUserRepository(db *gorm.DB) UserRepository {
	var user models.User
	return &userRepository{
		db:        db,
		tableName: user.TableName(),
	}
}

func (d *userRepository) WithTx(txHandle *gorm.DB) UserRepository {
	if txHandle == nil {
		log.Error("Transaction not found")
		return d
	}
	d.db = txHandle
	return d
}

func (r *userRepository) FindAll() ([]models.User, error) {
	var users []models.User
	result := r.db.Table(r.tableName).Preload("Accounts").Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

func (r *userRepository) FindByID(id string) (*models.User, error) {
	var user models.User
	result := r.db.Table(r.tableName).Preload("Accounts").Where("id = ?", id).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	result := r.db.Table(r.tableName).Where("email = ?", email).First(&user)
	if result.Error != nil {
		log.Error(result.Error)
		return nil, result.Error
	}

	return &user, nil
}

func (r *userRepository) FindByUniqueIdentifier(uniqueIdentifier string) (*models.User, error) {
	// Find by identity number or customer number
	var user models.User
	result := r.db.Table(r.tableName).Where("\"identity_number\" = ? OR \"customer_number\" = ?", uniqueIdentifier, uniqueIdentifier).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

func (r *userRepository) Create(user models.User) (*models.User, error) {
	result := r.db.Table(r.tableName).Create(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (r *userRepository) SoftDelete(id string) error {
	result := r.db.Table(r.tableName).Where("id = ?", id).Update("is_active", false)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *userRepository) SetTokenBlacklist(ctx *context.Context, key string, value string, exp time.Duration) error {
	err := r.redisClient.Set(*ctx, key, value, exp).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *userRepository) GetTokenBlacklist(ctx *context.Context, key string) (string, error) {
	value, err := r.redisClient.Get(*ctx, key).Result()
	if err != nil {
		return "", err
	}
	return value, nil
}
