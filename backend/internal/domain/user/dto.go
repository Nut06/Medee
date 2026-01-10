package user

// Profile Commands
type SkillDTO struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name"`
}

type UpdateProfileCommand struct {
	FirstName   string     `json:"firstName,omitempty"`
	LastName    string     `json:"lastName,omitempty"`
	Email       string     `json:"email,omitempty"`
	AvatarURL   *string    `json:"avatarURL,omitempty"`
	LinkedInURL string     `json:"linkedInURL,omitempty"`
	GitHubURL   string     `json:"githubURL,omitempty"`
	WebsiteURL  string     `json:"websiteURL,omitempty"`
	UserSkills  []SkillDTO `json:"userSkills,omitempty"`
}

// Experience Commands
type AddExperienceCommand struct {
	Position    string  `json:"position"`
	CompanyName string  `json:"companyName"`
	StartDate   *string `json:"startDate,omitempty"`
	EndDate     *string `json:"endDate,omitempty"`
	Description *string `json:"description,omitempty"`
}

type UpdateExperienceCommand struct {
	Position    string  `json:"position"`
	CompanyName string  `json:"companyName"`
	StartDate   *string `json:"startDate,omitempty"`
	EndDate     *string `json:"endDate,omitempty"`
	Description *string `json:"description,omitempty"`
}

// Education Commands
type AddEducationCommand struct {
	InstituteName  string `json:"instituteName"`
	Degree         string `json:"degree"`
	FieldOfStudy   string `json:"fieldOfStudy"`
	GraduationYear int    `json:"graduationYear"`
}

type UpdateEducationCommand struct {
	InstituteName  string `json:"instituteName"`
	Degree         string `json:"degree"`
	FieldOfStudy   string `json:"fieldOfStudy"`
	GraduationYear int    `json:"graduationYear"`
}

// Portfolio/Project Commands
type AddProjectCommand struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	ImageURL    string `json:"imageURL"`
	GithubURL   string `json:"githubURL"`
	DemoURL     string `json:"demoURL"`
}

type UpdateProjectCommand struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	ImageURL    string `json:"imageURL"`
	GithubURL   string `json:"githubURL"`
	DemoURL     string `json:"demoURL"`
}

// Skill Commands
type AddUserSkillCommand struct {
	SkillID *string `json:"skillId"` // Optional: if selecting from existing
	Name    *string `json:"name"`    // Optional: if creating new
}

type UpdateSkillsCommand struct {
	Skills []SkillDTO `json:"skills"`
}

type UpdateUserSkillCommand struct {
	Level string `json:"level"`
}

// Profile Responses
type UserProfileResponse struct {
	ID          string    `json:"id"`
	FirstName   string    `json:"firstName"`
	LastName    string    `json:"lastName"`
	Email       string    `json:"email"`
	AvatarURL   *string   `json:"avatarURL,omitempty"`
	ResumeURL   *string   `json:"resumeURL,omitempty"`
	Skill       *SkillDTO `json:"skill,omitempty"`
	Bio         *string   `json:"bio,omitempty"`
	LinkedInURL *string   `json:"linkedInURL,omitempty"`
	GitHubURL   *string   `json:"githubURL,omitempty"`
	WebsiteURL  *string   `json:"websiteURL,omitempty"`
}

type FullProfileResponse struct {
	ID          string                   `json:"id"`
	FirstName   string                   `json:"firstName"`
	LastName    string                   `json:"lastName"`
	Email       string                   `json:"email"`
	AvatarURL   *string                  `json:"avatarURL,omitempty"`
	ResumeURL   *string                  `json:"resumeURL,omitempty"`
	Bio         *string                  `json:"bio,omitempty"`
	LinkedInURL *string                  `json:"linkedInURL,omitempty"`
	GitHubURL   *string                  `json:"githubURL,omitempty"`
	WebsiteURL  *string                  `json:"websiteURL,omitempty"`
	Experiences []WorkExperienceResponse `json:"experiences,omitempty"`
	Educations  []EducationResponse      `json:"educations,omitempty"`
	Skills      []SkillResponse          `json:"skills,omitempty"`
	Projects    []ProjectResponse        `json:"projects,omitempty"`
}

// Work Experience Response
type WorkExperienceResponse struct {
	ID          string  `json:"id"`
	Position    string  `json:"position"`
	CompanyName string  `json:"companyName"`
	StartDate   string  `json:"startDate"`
	EndDate     *string `json:"endDate,omitempty"`
	Description *string `json:"description,omitempty"`
}

// Education Response
type EducationResponse struct {
	ID           string  `json:"id"`
	School       string  `json:"school"`
	Degree       string  `json:"degree"`
	FieldOfStudy *string `json:"fieldOfStudy,omitempty"`
	StartDate    string  `json:"startDate"`
	EndDate      *string `json:"endDate,omitempty"`
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
	ImageURL    *string `json:"imageURL,omitempty"`
	URL         *string `json:"url,omitempty"`
	GithubURL   *string `json:"githubURL,omitempty"`
}
