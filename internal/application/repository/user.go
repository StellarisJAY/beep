package repository

import (
	"beep/internal/types"
	"beep/internal/types/interfaces"
	"context"

	"gorm.io/gorm"
)

type UserRepoImpl struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) interfaces.UserRepo {
	return &UserRepoImpl{db: db}
}

func (u *UserRepoImpl) Create(ctx context.Context, user *types.User) error {
	return u.db.WithContext(ctx).Create(user).Error
}

func (u *UserRepoImpl) Update(ctx context.Context, user *types.User) error {
	return u.db.Model(user).WithContext(ctx).Updates(user).Error
}

func (u *UserRepoImpl) Delete(ctx context.Context, user *types.User) error {
	return u.db.WithContext(ctx).Delete(user, "id = ?", user.ID).Error
}

func (u *UserRepoImpl) FindById(ctx context.Context, userId string) (*types.User, error) {
	var user *types.User
	if err := u.db.WithContext(ctx).
		Model(&types.User{}).
		Where("id = ?", userId).
		First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserRepoImpl) FindByEmail(ctx context.Context, email string) (*types.User, error) {
	var user *types.User
	if err := u.db.WithContext(ctx).
		Model(&types.User{}).
		Where("email = ?", email).
		First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserRepoImpl) CheckPassword(ctx context.Context, email string, password string) (*types.User, error) {
	var user *types.User
	if err := u.db.WithContext(ctx).
		Model(&types.User{}).
		Where("email = ?", email).
		Where("password = ?", password).
		First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}
