package useradapter

import (
	"backend/internal/domain/domain"
	"backend/internal/domain/user"
	"context"
	"errors"
	"mime/multipart"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindById(ctx context.Context, id string) (*domain.User, error) {
	var u domain.User
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, user.ErrUserNotFound
		}
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) Update(ctx context.Context, id string, req *domain.User) (*domain.User, error) {
	var u domain.User
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, user.ErrUserNotFound
		}
		return nil, err
	}
	u.Name = req.Name
	u.Email = req.Email
	u.Role = req.Role
	u.Password = req.Password
	if err := r.db.WithContext(ctx).Save(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) UploadAvatar(ctx context.Context, id string, file *multipart.FileHeader) (*domain.User, error) {
	var u domain.User
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, user.ErrUserNotFound
		}
		return nil, err
	}
	u.AvatarURL = &file.Filename
	if err := r.db.WithContext(ctx).Save(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}