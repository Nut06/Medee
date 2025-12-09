package useradapter

import (
	"backend/internal/domain/domain"
	dto "backend/internal/domain/user"
	port "backend/internal/port/user"
	"context"
	"fmt"
	"mime/multipart"

	"github.com/google/uuid"
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
	user.Bio = req.Bio
	// user.Tagline = req.Tagline
	user.LinkedInURL = req.LinkedInURL
	user.GitHubURL = req.GitHubURL
	// user.PortfolioURL = req.PortfolioURL
	// user.ResumeURL = req.ResumeURL
	user.AvatarURL = req.AvatarURL

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

func (r *userRepository) GetUserCompanies(ctx context.Context, userID string) ([]domain.Company, error) {
	var members []domain.CompanyMember
	if err := r.db.WithContext(ctx).Preload("Company").Where("user_id = ?", userID).Find(&members).Error; err != nil {
		return nil, err
	}

	var companies []domain.Company
	for _, m := range members {
		companies = append(companies, m.Company)
	}
	return companies, nil
}

// Candidate Features

func (r *userRepository) GetFullProfile(ctx context.Context, id string) (*domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).
		Preload("ApplicantProfile").
		Preload("WorkExperiences").
		Preload("Educations").
		Preload("UserSkills.Skill").
		Preload("PortfolioItems").
		Where("id = ?", id).First(&user).Error
	if err != nil {
		fmt.Printf("err %v", err)
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) AddExperience(ctx context.Context, experience *domain.WorkExperience) (*domain.WorkExperience, error) {
	if err := r.db.WithContext(ctx).Create(experience).Error; err != nil {
		return nil, err
	}
	return experience, nil
}

func (r *userRepository) UpdateExperience(ctx context.Context, experience *domain.WorkExperience) (*domain.WorkExperience, error) {
	if err := r.db.WithContext(ctx).Save(experience).Error; err != nil {
		return nil, err
	}
	return experience, nil
}

func (r *userRepository) DeleteExperience(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&domain.WorkExperience{}, "id = ?", id).Error
}

func (r *userRepository) AddEducation(ctx context.Context, education *domain.Education) (*domain.Education, error) {
	if err := r.db.WithContext(ctx).Create(education).Error; err != nil {
		return nil, err
	}
	return education, nil
}

func (r *userRepository) UpdateEducation(ctx context.Context, education *domain.Education) (*domain.Education, error) {
	if err := r.db.WithContext(ctx).Save(education).Error; err != nil {
		return nil, err
	}
	return education, nil
}

func (r *userRepository) DeleteEducation(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&domain.Education{}, "id = ?", id).Error
}

func (r *userRepository) UpdateSkills(ctx context.Context, userID string, skills []string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Clear existing skills
		if err := tx.Where("user_id = ?", userID).Delete(&domain.UserSkill{}).Error; err != nil {
			return err
		}

		// 2. Add new skills
		for _, skillName := range skills {
			// Find or Create Skill
			var skill domain.Skill
			if err := tx.Where("name = ?", skillName).FirstOrCreate(&skill, domain.Skill{Name: skillName}).Error; err != nil {
				return err
			}

			// Link to User
			if err := tx.Create(&domain.UserSkill{
				UserID:  uuid.MustParse(userID),
				SkillID: skill.ID,
			}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *userRepository) AddProject(ctx context.Context, project *domain.PortfolioItem) (*domain.PortfolioItem, error) {
	if err := r.db.WithContext(ctx).Create(project).Error; err != nil {
		return nil, err
	}
	return project, nil
}

func (r *userRepository) UpdateProject(ctx context.Context, project *domain.PortfolioItem) (*domain.PortfolioItem, error) {
	if err := r.db.WithContext(ctx).Save(project).Error; err != nil {
		return nil, err
	}
	return project, nil
}

func (r *userRepository) DeleteProject(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&domain.PortfolioItem{}, "id = ?", id).Error
}

func (r *userRepository) UpdateSkill(ctx context.Context, skill *domain.Skill) error {
	return r.db.WithContext(ctx).Save(skill).Error
}

func (r *userRepository) AddSkill(ctx context.Context, userID string, req *dto.AddUserSkillCommand) error {
	var skillID uuid.UUID

	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Determine Skill ID
		if req.SkillID != nil && *req.SkillID != "" {
			// User selected existing skill
			skillID = uuid.MustParse(*req.SkillID)
		} else if req.Name != nil && *req.Name != "" {
			// User entered a new skill name
			var skill domain.Skill
			// Find or Create Master Skill
			if err := tx.Where("name = ?", *req.Name).FirstOrCreate(&skill, domain.Skill{Name: *req.Name}).Error; err != nil {
				return err
			}
			skillID = skill.ID
		} else {
			return dto.ErrSkillRequired
		}

		// 2. Create Relation (UserSkill)
		// Check if already exists to avoid duplicate
		var count int64
		tx.Model(&domain.UserSkill{}).Where("user_id = ? AND skill_id = ?", userID, skillID).Count(&count)
		if count > 0 {
			return dto.ErrSkillAlreadyAdded
		}

		userSkill := domain.UserSkill{
			UserID:  uuid.MustParse(userID),
			SkillID: skillID,
		}
		if err := tx.Create(&userSkill).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *userRepository) DeleteSkill(ctx context.Context, userID string, skillID string) error {
	result := r.db.WithContext(ctx).Delete(&domain.UserSkill{}, "user_id = ? AND skill_id = ?", userID, skillID)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return dto.ErrSkillNotFound
	}
	return nil
}

func (r *userRepository) GetSkills(ctx context.Context, userId string) ([]domain.Skill, error) {
	var skills []domain.Skill
	if err := r.db.WithContext(ctx).Where("user_id = ?", userId).Find(&skills).Error; err != nil {
		return nil, err
	}
	return skills, nil
}
