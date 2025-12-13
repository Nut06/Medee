package userapp

import (
	"backend/internal/domain/domain"
	"backend/internal/domain/user"
	skillport "backend/internal/port/skill"
	userport "backend/internal/port/user"
	"context"
	"mime/multipart"

	"github.com/google/uuid"
)

type Usecase struct {
	userRepo  userport.UserRepository
	skillRepo skillport.SkillRepository
}

func NewUsecase(repo userport.UserRepository, skillRepo skillport.SkillRepository) *Usecase {
	return &Usecase{userRepo: repo, skillRepo: skillRepo}
}

func (s *Usecase) GetProfile(ctx context.Context, id string) (*user.UserProfileResponse, error) {
	domainUser, err := s.userRepo.FindById(ctx, id)
	if err != nil {
		return nil, err
	}
	return user.ToUserProfileResponse(domainUser), nil
}

func (s *Usecase) UpdateProfile(ctx context.Context, id string, req *user.UpdateProfileCommand) (*user.UserProfileResponse, error) {
	var skillIds []string
	if req.UserSkills != nil {
		for _, skill := range req.UserSkills {
			if skill.ID != "" {
				skillObj, err := s.skillRepo.FindSkillById(ctx, skill.ID)
				if err != nil {
					return nil, err
				}
				skillIds = append(skillIds, skillObj.ID.String())

			} else if skill.Name != "" {
				skillObj, err := s.skillRepo.FindSkillByName(ctx, skill.Name)
				if err != nil {
					return nil, err
				}
				skillIds = append(skillIds, skillObj.ID.String())
			}
		}
	}
	userInput := domain.User{
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Email:       req.Email,
		LinkedInURL: &req.LinkedInURL,
		GitHubURL:   &req.GitHubURL,
		WebsiteURL:  &req.WebsiteURL,
	}
	domainUser, err := s.userRepo.Update(ctx, id, &userInput)
	if err != nil {
		return nil, err
	}

	if len(skillIds) > 0 {
		err = s.userRepo.UpdateSkills(ctx, id, skillIds)
		if err != nil {
			return nil, err
		}
	}

	return user.ToUserProfileResponse(domainUser), nil
}

// UploadAvatar uploads avatar and returns updated profile
func (s *Usecase) UploadAvatar(ctx context.Context, id string, file *multipart.FileHeader) (*user.UserProfileResponse, error) {
	domainUser, err := s.userRepo.UploadAvatar(ctx, id, file)
	if err != nil {
		return nil, err
	}
	return user.ToUserProfileResponse(domainUser), nil
}

// DeleteAvatar deletes avatar and returns updated profile
func (s *Usecase) DeleteAvatar(ctx context.Context, id string) (*user.UserProfileResponse, error) {
	domainUser, err := s.userRepo.DeleteAvatar(ctx, id)
	if err != nil {
		return nil, err
	}
	return user.ToUserProfileResponse(domainUser), nil
}

func (s *Usecase) UploadResume(ctx context.Context, id string, file *multipart.FileHeader) (*user.UserProfileResponse, error) {
	domainUser, err := s.userRepo.UploadResume(ctx, id, file)
	if err != nil {
		return nil, err
	}
	return user.ToUserProfileResponse(domainUser), nil
}

// DeleteResume deletes resume and returns updated profile
func (s *Usecase) DeleteResume(ctx context.Context, id string) (*user.UserProfileResponse, error) {
	domainUser, err := s.userRepo.DeleteResume(ctx, id)
	if err != nil {
		return nil, err
	}
	return user.ToUserProfileResponse(domainUser), nil
}

// GetFullProfile returns complete user profile with all related data
func (s *Usecase) GetFullProfile(ctx context.Context, id string) (*user.FullProfileResponse, error) {
	domainUser, err := s.userRepo.GetFullProfile(ctx, id)
	if err != nil {
		return nil, err
	}
	return user.ToFullProfileResponse(domainUser), nil
}

// AddExperience adds work experience and returns the created experience
func (s *Usecase) AddExperience(ctx context.Context, userID string, req *domain.WorkExperience) (*user.WorkExperienceResponse, error) {
	req.UserID = uuid.MustParse(userID)
	domainExp, err := s.userRepo.AddExperience(ctx, req)
	if err != nil {
		return nil, err
	}
	return user.ToWorkExperienceResponse(domainExp), nil
}

// UpdateExperience updates work experience and returns the updated experience
func (s *Usecase) UpdateExperience(ctx context.Context, id string, req *domain.WorkExperience) (*user.WorkExperienceResponse, error) {
	req.ID = uuid.MustParse(id)
	domainExp, err := s.userRepo.UpdateExperience(ctx, req)
	if err != nil {
		return nil, err
	}
	return user.ToWorkExperienceResponse(domainExp), nil
}

// DeleteExperience deletes work experience
func (s *Usecase) DeleteExperience(ctx context.Context, id string) error {
	return s.userRepo.DeleteExperience(ctx, id)
}

// AddEducation adds education and returns the created education
func (s *Usecase) AddEducation(ctx context.Context, userID string, req *domain.Education) (*user.EducationResponse, error) {
	req.UserID = uuid.MustParse(userID)
	domainEdu, err := s.userRepo.AddEducation(ctx, req)
	if err != nil {
		return nil, err
	}
	return user.ToEducationResponse(domainEdu), nil
}

// UpdateEducation updates education and returns the updated education
func (s *Usecase) UpdateEducation(ctx context.Context, id string, req *domain.Education) (*user.EducationResponse, error) {
	req.ID = uuid.MustParse(id)
	domainEdu, err := s.userRepo.UpdateEducation(ctx, req)
	if err != nil {
		return nil, err
	}
	return user.ToEducationResponse(domainEdu), nil
}

// DeleteEducation deletes education
func (s *Usecase) DeleteEducation(ctx context.Context, id string) error {
	return s.userRepo.DeleteEducation(ctx, id)
}

// UpdateSkills updates user skills
func (s *Usecase) UpdateSkills(ctx context.Context, userID string, skills []string) error {
	return s.userRepo.UpdateSkills(ctx, userID, skills)
}

// AddProject adds project and returns the created project
func (s *Usecase) AddProject(ctx context.Context, userID string, req *domain.PortfolioItem) (*user.ProjectResponse, error) {
	req.UserID = uuid.MustParse(userID)
	domainProj, err := s.userRepo.AddProject(ctx, req)
	if err != nil {
		return nil, err
	}
	return user.ToProjectResponse(domainProj), nil
}

// UpdateProject updates project and returns the updated project
func (s *Usecase) UpdateProject(ctx context.Context, id string, req *domain.PortfolioItem) (*user.ProjectResponse, error) {
	req.ID = uuid.MustParse(id)
	domainProj, err := s.userRepo.UpdateProject(ctx, req)
	if err != nil {
		return nil, err
	}
	return user.ToProjectResponse(domainProj), nil
}

// DeleteProject deletes project
func (s *Usecase) DeleteProject(ctx context.Context, id string) error {
	return s.userRepo.DeleteProject(ctx, id)
}

// AddSkill adds a skill to user profile
func (s *Usecase) AddSkill(ctx context.Context, userID string, req *user.AddUserSkillCommand) error {
	return s.userRepo.AddSkill(ctx, userID, req)
}

// DeleteSkill removes a skill from user profile
func (s *Usecase) DeleteSkill(ctx context.Context, userID string, skillId string) error {
	return s.userRepo.DeleteSkill(ctx, userID, skillId)
}
