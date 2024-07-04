package repository

import (
	"gorm.io/gorm"
	"tek-bank/internal/db/models"
)

type UserRepository interface {
	FindAll() ([]models.User, error)
	FindByID(id int) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	FindByUniqueIdentifier(uniqueIdentifier string) (*models.User, error)
	Create(user models.User) error
	SoftDelete(id int) error
}

type userRepository struct {
	db        *gorm.DB
	tableName string
}

func NewUserRepository(db *gorm.DB) UserRepository {
	var user models.User
	return &userRepository{
		db:        db,
		tableName: user.TableName(),
	}
}

func (r *userRepository) FindAll() ([]models.User, error) {
	var users []models.User
	result := r.db.Table(r.tableName).Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

func (r *userRepository) FindByID(id int) (*models.User, error) {
	var user models.User
	result := r.db.Table(r.tableName).Where("id = ?", id).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	result := r.db.Table(r.tableName).Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (r *userRepository) FindByUniqueIdentifier(uniqueIdentifier string) (*models.User, error) {
	var user models.User
	result := r.db.Table(r.tableName).Where("email = ? OR username = ?", uniqueIdentifier, uniqueIdentifier).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (r *userRepository) Create(user models.User) error {
	result := r.db.Table(r.tableName).Create(&user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *userRepository) SoftDelete(id int) error {
	result := r.db.Table(r.tableName).Where("id = ?", id).Update("is_active", false)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
