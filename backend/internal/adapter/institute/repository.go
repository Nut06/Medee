package institute_adapter

import (
	"backend/internal/domain/domain"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gofiber/storage/redis/v3"
	"gorm.io/gorm"
)

type instituteRepository struct {
	db    *gorm.DB
	redis *redis.Storage
}

func NewInstituteRepository(db *gorm.DB, redisClient *redis.Storage) *instituteRepository {
	return &instituteRepository{
		db:    db,
		redis: redisClient,
	}
}

// CreateInstitute creates a new institute
func (r *instituteRepository) CreateInstitute(ctx context.Context, institute *domain.Institute) (*domain.Institute, error) {
	if err := r.db.WithContext(ctx).Create(institute).Error; err != nil {
		return nil, err
	}
	return institute, nil
}

// SearchInstitutes searches institutes by name (for autocomplete)
// Uses flexible matching to handle apostrophes and special characters
func (r *instituteRepository) SearchInstitutes(ctx context.Context, query string) ([]domain.Institute, error) {
	var institutes []domain.Institute

	// Normalize the search by removing apostrophes and special chars
	// This allows "Mongkuts" to match "Mongkut's"
	err := r.db.WithContext(ctx).
		Where("REGEXP_REPLACE(name, '[^a-zA-Z0-9 ]', '', 'g') ILIKE ?", "%"+query+"%").
		Limit(10).
		Order("name ASC").
		Find(&institutes).Error

	return institutes, err
}

// FindInstituteById finds institute by ID
func (r *instituteRepository) FindInstituteById(ctx context.Context, id string) (*domain.Institute, error) {
	var institute domain.Institute
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&institute).Error
	return &institute, err
}

// FindInstituteByName finds institute by exact name
func (r *instituteRepository) FindInstituteByName(ctx context.Context, name string) (*domain.Institute, error) {
	var institute domain.Institute
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&institute).Error
	return &institute, err
}

// FindOrCreateInstitute finds existing or creates new institute
func (r *instituteRepository) FindOrCreateInstitute(ctx context.Context, name string) (*domain.Institute, error) {
	// Try to find first
	institute, err := r.FindInstituteByName(ctx, name)
	if err == nil {
		return institute, nil
	}

	// Not found, create new
	institute = &domain.Institute{Name: name}
	return r.CreateInstitute(ctx, institute)
}

// NEW: Cache methods for external API integration
func (r *instituteRepository) SearchInstitutesFromCache(ctx context.Context, query string, country string) ([]domain.Institute, bool, error) {
	if r.redis == nil {
		return nil, false, nil
	}

	cacheKey := fmt.Sprintf("universities:search:%s:%s", query, country)

	// val, err := r.redis.Get(ctx, cacheKey).Result()
	val, err := r.redis.Get(cacheKey)
	
	if err != nil {
		return nil, false, err // Redis error
	}

	var institutes []domain.Institute
	if err := json.Unmarshal([]byte(val), &institutes); err != nil {
		return nil, false, err
	}

	return institutes, true, nil
}

func (r *instituteRepository) CacheInstitutes(ctx context.Context, key string, institutes []domain.Institute, ttl time.Duration) error {
	if r.redis == nil {
		return nil
	}

	data, err := json.Marshal(institutes)
	if err != nil {
		return err
	}

	return r.redis.Set(key, data, ttl)
}

func (r *instituteRepository) SearchInstitutesWithCountry(ctx context.Context, query string, country string) ([]domain.Institute, error) {
	var institutes []domain.Institute

	db := r.db.WithContext(ctx).Where("name ILIKE ?", "%"+query+"%")
	if country != "" {
		db = db.Where("country = ?", country)
	}

	if err := db.Limit(50).Find(&institutes).Error; err != nil {
		return nil, err
	}

	return institutes, nil
}
