package user

import (
	"backend/internal/domain/domain"
	"fmt"
	"time"
)

// ============================================
// Mapper Functions: Domain -> DTO
// ============================================

func ToUserProfileResponse(u *domain.User) *UserProfileResponse {
	if u == nil {
		return nil
	}

	var skillDTO *SkillDTO
	if len(u.UserSkills) > 0 && u.UserSkills[0].Skill != nil {
		skillDTO = &SkillDTO{
			ID:   u.UserSkills[0].Skill.ID.String(),
			Name: u.UserSkills[0].Skill.Name,
		}
	}

	return &UserProfileResponse{
		ID:          u.ID.String(),
		FirstName:   u.FirstName,
		LastName:    u.LastName,
		Email:       u.Email,
		AvatarURL:   u.AvatarURL,
		Skill:       skillDTO,
		Bio:         u.Bio,
		ResumeURL:   u.ResumeURL,
		LinkedInURL: u.LinkedInURL,
		GitHubURL:   u.GitHubURL,
		WebsiteURL:  u.WebsiteURL,
	}
}

func ToFullProfileResponse(u *domain.User) *FullProfileResponse {
	if u == nil {
		return nil
	}

	resp := &FullProfileResponse{
		ID:          u.ID.String(),
		FirstName:   u.FirstName,
		LastName:    u.LastName,
		Email:       u.Email,
		AvatarURL:   u.AvatarURL,
		ResumeURL:   u.ResumeURL,
		Bio:         u.Bio,
		LinkedInURL: u.LinkedInURL,
		GitHubURL:   u.GitHubURL,
		WebsiteURL:  u.WebsiteURL,
	}

	// Map Experiences
	if u.WorkExperiences != nil {
		resp.Experiences = make([]WorkExperienceResponse, len(u.WorkExperiences))
		for i, exp := range u.WorkExperiences {
			resp.Experiences[i] = *ToWorkExperienceResponse(&exp)
		}
	}

	// Map Educations
	if u.Educations != nil {
		resp.Educations = make([]EducationResponse, len(u.Educations))
		for i, edu := range u.Educations {
			resp.Educations[i] = *ToEducationResponse(&edu)
		}
	}

	// Map Skills
	if u.UserSkills != nil {
		resp.Skills = make([]SkillResponse, len(u.UserSkills))
		for i, userSkill := range u.UserSkills {
			resp.Skills[i] = *ToSkillResponse(&userSkill)
		}
	}

	// Map Projects
	if u.PortfolioItems != nil {
		resp.Projects = make([]ProjectResponse, len(u.PortfolioItems))
		for i, proj := range u.PortfolioItems {
			resp.Projects[i] = *ToProjectResponse(&proj)
		}
	}

	return resp
}

func ToWorkExperienceResponse(exp *domain.WorkExperience) *WorkExperienceResponse {
	if exp == nil {
		return nil
	}

	var endDate *string
	if exp.EndDate != nil {
		formatted := exp.EndDate.Format(time.RFC3339)
		endDate = &formatted
	}

	var startDate string
	if exp.StartDate != nil {
		startDate = exp.StartDate.Format(time.RFC3339)
	}

	return &WorkExperienceResponse{
		ID:          exp.ID.String(),
		Position:    exp.Position,
		CompanyName: exp.CompanyName,
		StartDate:   startDate,
		EndDate:     endDate,
		Description: exp.Description,
	}
}

func ToEducationResponse(edu *domain.Education) *EducationResponse {
	if edu == nil {
		return nil
	}

	instituteName := edu.InstituteName
	if instituteName == "" {
		instituteName = edu.InstituteID.String()
	}

	var fieldOfStudyName *string
	if edu.FieldOfStudyName != "" {
		fieldOfStudyName = &edu.FieldOfStudyName
	}

	var startDate string
	if edu.StartDate != nil {
		startDate = edu.StartDate.Format(time.RFC3339)
	}

	var endDate *string
	if edu.EndDate != nil {
		formatted := edu.EndDate.Format(time.RFC3339)
		endDate = &formatted
	}

	var gpa *string
	if edu.GPA != nil {
		gpaStr := fmt.Sprintf("%.2f", *edu.GPA)
		gpa = &gpaStr
	}

	return &EducationResponse{
		ID:           edu.ID.String(),
		School:       instituteName,
		Degree:       edu.Degree,
		FieldOfStudy: fieldOfStudyName,
		StartDate:    startDate,
		EndDate:      endDate,
		GPA:          gpa,
		Description:  edu.Description,
	}
}

func ToSkillResponse(userSkill *domain.UserSkill) *SkillResponse {
	if userSkill == nil || userSkill.Skill == nil {
		return nil
	}
	return &SkillResponse{
		ID:   userSkill.Skill.ID.String(),
		Name: userSkill.Skill.Name,
	}
}

func ToProjectResponse(proj *domain.PortfolioItem) *ProjectResponse {
	if proj == nil {
		return nil
	}
	// Convert Description from string to *string
	desc := proj.Description
	return &ProjectResponse{
		ID:          proj.ID.String(),
		Title:       proj.Title,
		Description: &desc,
		ImageURL:    proj.ImageURL,
		URL:         proj.DemoURL, // Fixed: domain uses DemoURL not URL
		GithubURL:   proj.GithubURL,
	}
}
