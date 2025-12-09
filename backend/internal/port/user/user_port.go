package userport

import (
	"backend/internal/domain/domain"
	dto "backend/internal/domain/user"
	"context"
	"mime/multipart"
)

type UserService interface {
	GetProfile(ctx context.Context, id string) (*domain.User, error)
	UpdateProfile(ctx context.Context, id string, req *domain.User) (*domain.User, error)
	UploadAvatar(ctx context.Context, id string, file *multipart.FileHeader) (*domain.User, error)
	DeleteAvatar(ctx context.Context, id string) (*domain.User, error)

	// Candidate Features
	GetFullProfile(ctx context.Context, id string) (*domain.User, error)
	AddExperience(ctx context.Context, userID string, req *domain.WorkExperience) (*domain.WorkExperience, error)
	UpdateExperience(ctx context.Context, id string, req *domain.WorkExperience) (*domain.WorkExperience, error)
	DeleteExperience(ctx context.Context, id string) error
	AddEducation(ctx context.Context, userID string, req *domain.Education) (*domain.Education, error)
	UpdateEducation(ctx context.Context, id string, req *domain.Education) (*domain.Education, error)
	DeleteEducation(ctx context.Context, id string) error
	UpdateSkills(ctx context.Context, userID string, skills []string) error
	AddProject(ctx context.Context, userID string, req *domain.PortfolioItem) (*domain.PortfolioItem, error)
	UpdateProject(ctx context.Context, id string, req *domain.PortfolioItem) (*domain.PortfolioItem, error)
	DeleteProject(ctx context.Context, id string) error
	AddSkill(ctx context.Context, userID string, req *dto.AddUserSkillCommand) error
	DeleteSkill(ctx context.Context, userID string, skillID string) error
}

type UserRepository interface {
	FindById(ctx context.Context, id string) (*domain.User, error)
	Update(ctx context.Context, id string, req *domain.User) (*domain.User, error)
	UploadAvatar(ctx context.Context, id string, file *multipart.FileHeader) (*domain.User, error)
	DeleteAvatar(ctx context.Context, id string) (*domain.User, error)
	GetUserCompanies(ctx context.Context, userID string) ([]domain.Company, error)

	// Candidate Features
	GetFullProfile(ctx context.Context, id string) (*domain.User, error)
	AddExperience(ctx context.Context, experience *domain.WorkExperience) (*domain.WorkExperience, error)
	UpdateExperience(ctx context.Context, experience *domain.WorkExperience) (*domain.WorkExperience, error)
	DeleteExperience(ctx context.Context, id string) error
	AddEducation(ctx context.Context, education *domain.Education) (*domain.Education, error)
	UpdateEducation(ctx context.Context, education *domain.Education) (*domain.Education, error)
	DeleteEducation(ctx context.Context, id string) error
	UpdateSkills(ctx context.Context, userID string, skills []string) error
	AddProject(ctx context.Context, project *domain.PortfolioItem) (*domain.PortfolioItem, error)
	UpdateProject(ctx context.Context, project *domain.PortfolioItem) (*domain.PortfolioItem, error)
	DeleteProject(ctx context.Context, id string) error
	AddSkill(ctx context.Context, userID string, req *dto.AddUserSkillCommand) error
	DeleteSkill(ctx context.Context, userID string, skillID string) error
}
