package useradapter

import (
	"backend/internal/domain/domain"
	port "backend/internal/port/user"
	"context"
	"mime/multipart"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) port.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindById(ctx context.Context, id string) (*domain.User, error) {
	var user domain.User
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(ctx context.Context, id string, req *domain.User) (*domain.User, error) {
	var user domain.User
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}

	// Update fields
	user.FirstName = req.FirstName
	user.LastName = req.LastName
	user.PhoneNumber = req.PhoneNumber
	// user.Bio = req.Bio // Assuming Bio field exists or will be added

	if err := r.db.WithContext(ctx).Save(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) UploadAvatar(ctx context.Context, id string, file *multipart.FileHeader) (*domain.User, error) {
	// In a real application, you would upload the file to a storage service (S3, GCS, etc.) here.
	// For now, we'll just simulate it by returning the filename as the URL.
	// You might want to implement a separate Storage Port/Adapter for this.

	avatarURL := "https://example.com/uploads/" + file.Filename // Placeholder

	var user domain.User
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}

	user.AvatarURL = &avatarURL
	if err := r.db.WithContext(ctx).Save(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) DeleteAvatar(ctx context.Context, id string) (*domain.User, error) {
	var user domain.User
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}

	user.AvatarURL = nil
	if err := r.db.WithContext(ctx).Save(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
