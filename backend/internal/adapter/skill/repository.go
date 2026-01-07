package skill_adapter

import (
	"backend/internal/domain/domain"
	port "backend/internal/port/skill"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type skillRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) port.SkillRepository {
	return &skillRepository{db: db}
}

func (r *skillRepository) FindSkillById(ctx context.Context, id string) (*domain.Skill, error) {
	var skill domain.Skill
	if err := r.db.WithContext(ctx).First(&skill, uuid.MustParse(id)).Error; err != nil {
		return nil, err
	}
	return &skill, nil
}

func (r *skillRepository) FindSkillByName(ctx context.Context, name string) (*domain.Skill, error) {
	var skill domain.Skill
	if err := r.db.WithContext(ctx).First(&skill, name).Error; err != nil {
		return nil, err
	}
	return &skill, nil
}

func (r *skillRepository) CreateSkill(ctx context.Context, skill *domain.Skill) (*domain.Skill, error) {
	if err := r.db.WithContext(ctx).Create(skill).Error; err != nil {
		return nil, err
	}
	return skill, nil
}

func (r *skillRepository) SearchSkills(ctx context.Context, query string) ([]domain.Skill, error) {
	var skills []domain.Skill
	if err := r.db.WithContext(ctx).Where("name ILIKE ?", "%"+query+"%").Find(&skills).Error; err != nil {
		return nil, err
	}
	return skills, nil
}
