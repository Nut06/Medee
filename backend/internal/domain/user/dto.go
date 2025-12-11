package user

import "time"

// Profile Commands
type SkillDTO struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name"`
}

type UpdateProfileCommand struct {
	FirstName   string     `json:"first_name,omitempty"`
	LastName    string     `json:"last_name,omitempty"`
	Email       string     `json:"email,omitempty"`
	LinkedInURL string     `json:"linkedin_url,omitempty"`
	GitHubURL   string     `json:"github_url,omitempty"`
	WebsiteURL  string     `json:"website_url,omitempty"`
	UserSkills  []SkillDTO `json:"user_skills,omitempty"`
}

// Experience Commands
type AddExperienceCommand struct {
	Position    string    `json:"position"`
	CompanyName string    `json:"company_name"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	Description string    `json:"description"`
}

type UpdateExperienceCommand struct {
	Position    string    `json:"position"`
	CompanyName string    `json:"company_name"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	Description string    `json:"description"`
}

// Education Commands
type AddEducationCommand struct {
	InstituteName  string `json:"institute_name"`
	Degree         string `json:"degree"`
	FieldOfStudy   string `json:"field_of_study"`
	GraduationYear int    `json:"graduation_year"`
}

type UpdateEducationCommand struct {
	InstituteName  string `json:"institute_name"`
	Degree         string `json:"degree"`
	FieldOfStudy   string `json:"field_of_study"`
	GraduationYear int    `json:"graduation_year"`
}

// Portfolio/Project Commands
type AddProjectCommand struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`
	GithubURL   string `json:"github_url"`
	DemoURL     string `json:"demo_url"`
}

type UpdateProjectCommand struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`
	GithubURL   string `json:"github_url"`
	DemoURL     string `json:"demo_url"`
}

// Skill Commands
type AddUserSkillCommand struct {
	SkillID *string `json:"skill_id"` // Optional: if selecting from existing
	Name    *string `json:"name"`     // Optional: if creating new
}

type UpdateUserSkillCommand struct {
	Level string `json:"level"`
}

// Profile Responses
type UserProfileResponse struct {
	ID          string    `json:"id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Email       string    `json:"email"`
	AvatarURL   *string   `json:"avatar_url,omitempty"`
	Skill       *SkillDTO `json:"skill,omitempty"`
	Bio         *string   `json:"bio,omitempty"`
	LinkedInURL *string   `json:"linkedin_url,omitempty"`
	GitHubURL   *string   `json:"github_url,omitempty"`
	WebsiteURL  *string   `json:"website_url,omitempty"`
}

type FullProfileResponse struct {
	ID          string                   `json:"id"`
	FirstName   string                   `json:"first_name"`
	LastName    string                   `json:"last_name"`
	Email       string                   `json:"email"`
	AvatarURL   *string                  `json:"avatar_url,omitempty"`
	Bio         *string                  `json:"bio,omitempty"`
	LinkedInURL *string                  `json:"linkedin_url,omitempty"`
	GitHubURL   *string                  `json:"github_url,omitempty"`
	WebsiteURL  *string                  `json:"website_url,omitempty"`
	Experiences []WorkExperienceResponse `json:"experiences,omitempty"`
	Educations  []EducationResponse      `json:"educations,omitempty"`
	Skills      []SkillResponse          `json:"skills,omitempty"`
	Projects    []ProjectResponse        `json:"projects,omitempty"`
}

// Work Experience Response
type WorkExperienceResponse struct {
	ID          string  `json:"id"`
	Position    string  `json:"position"`
	CompanyName string  `json:"company_name"`
	StartDate   string  `json:"start_date"`
	EndDate     *string `json:"end_date,omitempty"`
	Description *string `json:"description,omitempty"`
}

// Education Response
type EducationResponse struct {
	ID           string  `json:"id"`
	School       string  `json:"school"`
	Degree       string  `json:"degree"`
	FieldOfStudy *string `json:"field_of_study,omitempty"`
	StartDate    string  `json:"start_date"`
	EndDate      *string `json:"end_date,omitempty"`
	GPA          *string `json:"gpa,omitempty"`
	Description  *string `json:"description,omitempty"`
}

// Skill Response
type SkillResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Project Response
type ProjectResponse struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	Description *string `json:"description,omitempty"`
	ImageURL    *string `json:"image_url,omitempty"`
	URL         *string `json:"url,omitempty"`
	GithubURL   *string `json:"github_url,omitempty"`
}
