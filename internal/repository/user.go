package repository

import (
	"context"
	"todo-list/internal/entity"

	"gorm.io/gorm"
)

type UserRepository interface {
	FindAll(ctx context.Context) ([]entity.User, error)
	FindByUsername(ctx context.Context, username string) (*entity.User, error)
	CreateUser(ctx context.Context, user *entity.UserReg) error
}

type userRepository struct {
	db *gorm.DB
}


func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) FindAll(ctx context.Context) ([]entity.User, error) {
	user := make([]entity.User, 0)
	if err := r.db.WithContext(ctx).Find(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) FindByUsername(ctx context.Context, username string) (*entity.User, error) {
	user := new(entity.User)
	if err := r.db.WithContext(ctx).
		Where("username = ?", username).
		First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) CreateUser(ctx context.Context, user *entity.UserReg) error {
    return r.db.WithContext(ctx).Create(user).Error
}