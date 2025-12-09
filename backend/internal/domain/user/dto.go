package user

import "time"

// Profile Commands
type UpdateProfileCommand struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	LinkedInURL string `json:"linkedin_url"`
	GitHubURL   string `json:"github_url"`
	WebsiteURL  string `json:"website_url"`
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
