package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"secureops/backend-go/api/model"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.User{}).Where("username = ?", username).Count(&count).Error
	return count > 0, err
}

func (r *UserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.User{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}

func (r *UserRepository) Save(ctx context.Context, user model.User) error {
	return r.db.WithContext(ctx).Create(&user).Error
}

func (r *UserRepository) FindByUsernameOrEmail(ctx context.Context, userOrEmail string) (model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).
		Where("username = ? OR email = ?", userOrEmail, userOrEmail).
		First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return model.User{}, gorm.ErrRecordNotFound
	}
	return user, err
}
