package authadapter

import (
	"backend/internal/domain/auth"
	"backend/internal/domain/domain"
	authport "backend/internal/port/auth"
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

// FindByEmail returns a domain user by email or ErrUserNotFound.
func (r *GormRepository) FindByEmail(ctx context.Context, email string) (auth.User, error) {
	var u domain.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return auth.User{}, auth.ErrUserNotFound
		}
		return auth.User{}, err
	}
	return mapToDomainUser(&u), nil
}

func (r *GormRepository) FindByID(ctx context.Context, id string) (auth.User, error) {
	var u domain.User
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return auth.User{}, auth.ErrUserNotFound
		}
		return auth.User{}, err
	}
	return mapToDomainUser(&u), nil
}

// Create inserts a user and returns the persisted domain model.
func (r *GormRepository) Create(ctx context.Context, user auth.User) (auth.User, error) {
	entity := mapToEntityUser(user)
	if err := r.db.WithContext(ctx).Create(&entity).Error; err != nil {
		if isDuplicateKey(err) {
			return auth.User{}, auth.ErrEmailAlreadyUsed
		}
		return auth.User{}, err
	}
	return mapToDomainUser(entity), nil
}

// Save persists a refresh token (implements RefreshTokenStore).
func (r *GormRepository) Save(ctx context.Context, token string, userID uuid.UUID, expiresAt time.Time) error {
	rt := domain.RefreshToken{
		Token:     token,
		ExpiresAt: expiresAt,
		UserID:    userID.String(),
	}
	return r.db.WithContext(ctx).Create(&rt).Error
}

// Delete removes a refresh token by its value.a
func (r *GormRepository) Delete(ctx context.Context, token string) error {
	return r.db.WithContext(ctx).Where("token = ?", token).Delete(&domain.RefreshToken{}).Error
}

func (r *GormRepository) GetUserCompanies(ctx context.Context, userID string) ([]auth.Company, error) {
	var members []domain.CompanyMember
	if err := r.db.WithContext(ctx).Preload("Company").Where("user_id = ?", userID).Find(&members).Error; err != nil {
		return nil, err
	}

	var companies []auth.Company
	for _, m := range members {
		companies = append(companies, auth.Company{
			ID:   m.Company.ID.String(),
			Name: m.Company.Name,
			Role: string(m.Role),
		})
	}
	return companies, nil
}

func mapToDomainUser(u *domain.User) auth.User {
	pw := ""
	if u.Password != nil {
		pw = *u.Password
	}
	return auth.User{
		ID:           u.ID,
		FirstName:    u.FirstName,
		LastName:     u.LastName,
		Email:        u.Email,
		PasswordHash: pw,
	}
}

func mapToEntityUser(u auth.User) *domain.User {
	pw := u.PasswordHash

	return &domain.User{
		ID:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
		Password:  &pw,
	}
}

func isDuplicateKey(err error) bool {
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return true
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "duplicate") || strings.Contains(msg, "unique constraint")
}

// Interfaces satisfaction
var (
	_ authport.UserRepository    = (*GormRepository)(nil)
	_ authport.RefreshTokenStore = (*GormRepository)(nil)
)
