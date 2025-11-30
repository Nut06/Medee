package userapp

import (
	"backend/internal/domain/domain"
	port "backend/internal/port/user"
	"context"
	"mime/multipart"

	"github.com/google/uuid"
)

type UpdateProfileCommand struct {
	FirstName string
	LastName  string
	Email     string
}

type Usecase struct {
	repo port.UserRepository
}

func NewUsecase(repo port.UserRepository) *Usecase {
	return &Usecase{repo: repo}
}

func (s *Usecase) GetProfile(ctx context.Context, id string) (*domain.User, error) {
	return s.repo.FindById(ctx, id)
}

func (s *Usecase) UpdateProfile(ctx context.Context, id string, req *domain.User) (*domain.User, error) {
	return s.repo.Update(ctx, id, req)
}

func (s *Usecase) UploadAvatar(ctx context.Context, id string, file *multipart.FileHeader) (*domain.User, error) {
	return s.repo.UploadAvatar(ctx, id, file)
}

func (s *Usecase) DeleteAvatar(ctx context.Context, id string) (*domain.User, error) {
	return s.repo.DeleteAvatar(ctx, id)
}

// Candidate Features

func (s *Usecase) GetFullProfile(ctx context.Context, id string) (*domain.User, error) {
	return s.repo.GetFullProfile(ctx, id)
}

func (s *Usecase) AddExperience(ctx context.Context, userID string, req *domain.WorkExperience) (*domain.WorkExperience, error) {
	// Validate request if needed
	req.UserID = uuid.MustParse(userID)
	return s.repo.AddExperience(ctx, req)
}

func (s *Usecase) UpdateExperience(ctx context.Context, id string, req *domain.WorkExperience) (*domain.WorkExperience, error) {
	// Check ownership if needed (omitted for brevity, usually done in repo or by fetching first)
	req.ID = uuid.MustParse(id)
	return s.repo.UpdateExperience(ctx, req)
}

func (s *Usecase) DeleteExperience(ctx context.Context, id string) error {
	return s.repo.DeleteExperience(ctx, id)
}

func (s *Usecase) AddEducation(ctx context.Context, userID string, req *domain.Education) (*domain.Education, error) {
	req.UserID = uuid.MustParse(userID)
	return s.repo.AddEducation(ctx, req)
}

func (s *Usecase) UpdateEducation(ctx context.Context, id string, req *domain.Education) (*domain.Education, error) {
	req.ID = uuid.MustParse(id)
	return s.repo.UpdateEducation(ctx, req)
}

func (s *Usecase) DeleteEducation(ctx context.Context, id string) error {
	return s.repo.DeleteEducation(ctx, id)
}

func (s *Usecase) UpdateSkills(ctx context.Context, userID string, skills []string) error {
	return s.repo.UpdateSkills(ctx, userID, skills)
}

func (s *Usecase) AddProject(ctx context.Context, userID string, req *domain.PortfolioItem) (*domain.PortfolioItem, error) {
	req.UserID = uuid.MustParse(userID)
	return s.repo.AddProject(ctx, req)
}

func (s *Usecase) UpdateProject(ctx context.Context, id string, req *domain.PortfolioItem) (*domain.PortfolioItem, error) {
	req.ID = uuid.MustParse(id)
	return s.repo.UpdateProject(ctx, req)
}

func (s *Usecase) DeleteProject(ctx context.Context, id string) error {
	return s.repo.DeleteProject(ctx, id)
}
