package useradapter

import (
	field_of_study_adapter "backend/internal/adapter/field_of_study"
	institute_adapter "backend/internal/adapter/institute"
	skilladapter "backend/internal/adapter/skill"
	"backend/internal/application/fieldapp"
	"backend/internal/application/instituteapp"
	userapp "backend/internal/application/userapp"
	"backend/internal/domain/domain"
	user "backend/internal/domain/user"
	"errors"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type HTTPHandler struct {
	usecase          *userapp.Usecase
	instituteUsecase *instituteapp.Usecase
	fieldUsecase     *fieldapp.Usecase
}

func NewHTTPHandler(db *gorm.DB) *HTTPHandler {
	userRepo := NewRepository(db)
	skillRepo := skilladapter.NewRepository(db)
	usecase := userapp.NewUsecase(userRepo, skillRepo)

	// Institute and FieldOfStudy usecases
	instituteRepo := institute_adapter.NewInstituteRepository(db)
	instituteUsecase := instituteapp.NewUsecase(instituteRepo)

	fieldRepo := field_of_study_adapter.NewFieldOfStudyRepository(db)
	fieldUsecase := fieldapp.NewUsecase(fieldRepo)

	return &HTTPHandler{
		usecase:          usecase,
		instituteUsecase: instituteUsecase,
		fieldUsecase:     fieldUsecase,
	}
}

func (h *HTTPHandler) GetUser(c *fiber.Ctx) error {
	userId := c.Locals("user_id").(string)
	res, err := h.usecase.GetProfile(c.Context(), userId)
	if err != nil {
		return h.handleError(err)
	}
	return c.JSON(res)
}

func (h *HTTPHandler) UploadAvatar(c *fiber.Ctx) error {
	userId := c.Locals("user_id").(string)
	file, err := c.FormFile("avatar")
	if err != nil {
		return h.handleError(user.ErrInvalidRequestBody)
	}

	user, err := h.usecase.UploadAvatar(c.Context(), userId, file)
	if err != nil {
		return h.handleError(err)
	}
	return c.JSON(user)
}

func (h *HTTPHandler) UpdateProfile(c *fiber.Ctx) error {
	userId := c.Locals("user_id").(string)
	var req user.UpdateProfileCommand
	if err := c.BodyParser(&req); err != nil {
		return h.handleError(user.ErrInvalidRequestBody)
	}

	res, err := h.usecase.UpdateProfile(c.Context(), userId, &req)
	if err != nil {
		return h.handleError(err)
	}
	return c.JSON(res)
}

func (h *HTTPHandler) DeleteAvatar(c *fiber.Ctx) error {
	userId := c.Locals("user_id").(string)
	user, err := h.usecase.DeleteAvatar(c.Context(), userId)
	if err != nil {
		return h.handleError(err)
	}
	return c.JSON(user)
}

func (h *HTTPHandler) UploadResume(c *fiber.Ctx) error {
	userId := c.Locals("user_id").(string)
	file, err := c.FormFile("resume")
	if err != nil {
		return h.handleError(user.ErrInvalidRequestBody)
	}

	res, err := h.usecase.UploadResume(c.Context(), userId, file)
	if err != nil {
		return h.handleError(err)
	}
	return c.JSON(res)
}

func (h *HTTPHandler) DeleteResume(c *fiber.Ctx) error {
	userId := c.Locals("user_id").(string)
	res, err := h.usecase.DeleteResume(c.Context(), userId)
	if err != nil {
		return h.handleError(err)
	}
	return c.JSON(res)
}

func (h *HTTPHandler) GetFullProfile(c *fiber.Ctx) error {
	userId := c.Locals("user_id").(string)
	fmt.Printf("id: %v  from get full profile http handler", userId)
	res, err := h.usecase.GetFullProfile(c.Context(), userId)
	if err != nil {
		return h.handleError(err)
	}
	return c.JSON(res)
}

func (h *HTTPHandler) AddExperience(c *fiber.Ctx) error {
	userId := c.Locals("user_id").(string)
	var req user.AddExperienceCommand
	if err := c.BodyParser(&req); err != nil {
		return h.handleError(user.ErrInvalidRequestBody)
	}

	// Convert string dates to *time.Time
	var startDate, endDate *time.Time
	if req.StartDate != nil && *req.StartDate != "" {
		t, err := time.Parse("2006-01-02", *req.StartDate)
		if err != nil {
			return h.handleError(user.ErrInvalidRequestBody)
		}
		startDate = &t
	}
	if req.EndDate != nil && *req.EndDate != "" {
		t, err := time.Parse("2006-01-02", *req.EndDate)
		if err != nil {
			return h.handleError(user.ErrInvalidRequestBody)
		}
		endDate = &t
	}

	exp := &domain.WorkExperience{
		Position:    req.Position,
		CompanyName: req.CompanyName,
		StartDate:   startDate,
		EndDate:     endDate,
		Description: req.Description,
	}

	res, err := h.usecase.AddExperience(c.Context(), userId, exp)
	if err != nil {
		return h.handleError(err)
	}
	return c.JSON(res)
}

func (h *HTTPHandler) UpdateExperience(c *fiber.Ctx) error {
	id := c.Params("experienceId")
	var req user.UpdateExperienceCommand
	if err := c.BodyParser(&req); err != nil {
		return h.handleError(user.ErrInvalidRequestBody)
	}

	// Convert string dates to *time.Time
	var startDate, endDate *time.Time
	if req.StartDate != nil && *req.StartDate != "" {
		t, err := time.Parse("2006-01-02", *req.StartDate)
		if err != nil {
			return h.handleError(user.ErrInvalidRequestBody)
		}
		startDate = &t
	}
	if req.EndDate != nil && *req.EndDate != "" {
		t, err := time.Parse("2006-01-02", *req.EndDate)
		if err != nil {
			return h.handleError(user.ErrInvalidRequestBody)
		}
		endDate = &t
	}

	exp := &domain.WorkExperience{
		Position:    req.Position,
		CompanyName: req.CompanyName,
		StartDate:   startDate,
		EndDate:     endDate,
		Description: req.Description,
	}

	res, err := h.usecase.UpdateExperience(c.Context(), id, exp)
	if err != nil {
		return h.handleError(err)
	}
	return c.JSON(res)
}

func (h *HTTPHandler) DeleteExperience(c *fiber.Ctx) error {
	id := c.Params("experienceId")
	err := h.usecase.DeleteExperience(c.Context(), id)
	if err != nil {
		return h.handleError(err)
	}
	return c.SendStatus(fiber.StatusOK)
}

func (h *HTTPHandler) AddEducation(c *fiber.Ctx) error {
	userId := c.Locals("user_id").(string)
	var req user.AddEducationCommand
	if err := c.BodyParser(&req); err != nil {
		return h.handleError(user.ErrInvalidRequestBody)
	}

	// Find or create Institute
	institute, err := h.instituteUsecase.FindOrCreateInstitute(c.Context(), req.InstituteName)
	if err != nil {
		return h.handleError(err)
	}

	// Find or create FieldOfStudy
	field, err := h.fieldUsecase.FindOrCreateFieldOfStudy(c.Context(), req.FieldOfStudy)
	if err != nil {
		return h.handleError(err)
	}

	// Create Education entity
	edu := &domain.Education{
		UserID:         uuid.MustParse(userId),
		InstituteID:    institute.ID,
		FieldOfStudyID: field.ID,
		Degree:         req.Degree,
		GraduationYear: &req.GraduationYear,
	}

	res, err := h.usecase.AddEducation(c.Context(), userId, edu)
	if err != nil {
		return h.handleError(err)
	}
	return c.JSON(res)
}

func (h *HTTPHandler) UpdateEducation(c *fiber.Ctx) error {
	id := c.Params("educationId")
	var req user.UpdateEducationCommand
	if err := c.BodyParser(&req); err != nil {
		return h.handleError(user.ErrInvalidRequestBody)
	}

	// Find or create Institute
	institute, err := h.instituteUsecase.FindOrCreateInstitute(c.Context(), req.InstituteName)
	if err != nil {
		return h.handleError(err)
	}

	// Find or create FieldOfStudy
	field, err := h.fieldUsecase.FindOrCreateFieldOfStudy(c.Context(), req.FieldOfStudy)
	if err != nil {
		return h.handleError(err)
	}

	// Create Education entity
	edu := &domain.Education{
		InstituteID:    institute.ID,
		FieldOfStudyID: field.ID,
		Degree:         req.Degree,
		GraduationYear: &req.GraduationYear,
	}

	res, err := h.usecase.UpdateEducation(c.Context(), id, edu)
	if err != nil {
		return h.handleError(err)
	}
	return c.JSON(res)
}

func (h *HTTPHandler) DeleteEducation(c *fiber.Ctx) error {
	id := c.Params("educationId")
	err := h.usecase.DeleteEducation(c.Context(), id)
	if err != nil {
		return h.handleError(err)
	}
	return c.SendStatus(fiber.StatusOK)
}

func (h *HTTPHandler) UpdateSkills(c *fiber.Ctx) error {
	userId := c.Locals("user_id").(string)
	var req struct {
		Skills []string `json:"skills"`
	}
	if err := c.BodyParser(&req); err != nil {
		return h.handleError(err)
	}
	err := h.usecase.UpdateSkills(c.Context(), userId, req.Skills)
	if err != nil {
		return h.handleError(err)
	}
	return c.SendStatus(fiber.StatusOK)
}

func (h *HTTPHandler) AddProject(c *fiber.Ctx) error {
	userId := c.Locals("user_id").(string)
	var req domain.PortfolioItem
	if err := c.BodyParser(&req); err != nil {
		return h.handleError(err)
	}
	res, err := h.usecase.AddProject(c.Context(), userId, &req)
	if err != nil {
		return h.handleError(err)
	}
	return c.JSON(res)
}

func (h *HTTPHandler) UpdateProject(c *fiber.Ctx) error {
	id := c.Params("projectId")
	var req domain.PortfolioItem
	if err := c.BodyParser(&req); err != nil {
		return h.handleError(err)
	}
	res, err := h.usecase.UpdateProject(c.Context(), id, &req)
	if err != nil {
		return h.handleError(err)
	}
	return c.JSON(res)
}

func (h *HTTPHandler) DeleteProject(c *fiber.Ctx) error {
	id := c.Params("projectId")
	err := h.usecase.DeleteProject(c.Context(), id)
	if err != nil {
		return h.handleError(err)
	}
	return c.SendStatus(fiber.StatusOK)
}

func (h *HTTPHandler) AddSkill(c *fiber.Ctx) error {
	userId := c.Locals("user_id").(string)
	var req user.AddUserSkillCommand
	if err := c.BodyParser(&req); err != nil {
		return h.handleError(err)
	}
	err := h.usecase.AddSkill(c.Context(), userId, &req)
	if err != nil {
		return h.handleError(err)
	}
	return c.SendStatus(fiber.StatusCreated)
}

func (h *HTTPHandler) DeleteSkill(c *fiber.Ctx) error {
	userId := c.Locals("user_id").(string)
	id := c.Params("skillId")
	err := h.usecase.DeleteSkill(c.Context(), userId, id)
	if err != nil {
		return h.handleError(err)
	}
	return c.SendStatus(fiber.StatusOK)
}

func (h *HTTPHandler) handleError(err error) *fiber.Error {
	switch {
	case err == nil:
		return nil
	case errors.Is(err, user.ErrInvalidRequestBody):
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	case errors.Is(err, user.ErrUserNotFound):
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	case errors.Is(err, user.ErrSkillNotFound):
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	case errors.Is(err, user.ErrSkillAlreadyAdded):
		return fiber.NewError(fiber.StatusConflict, err.Error())
	case errors.Is(err, user.ErrSkillRequired):
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	case errors.Is(err, user.ErrExperienceNotFound):
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	case errors.Is(err, user.ErrEducationNotFound):
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	case errors.Is(err, user.ErrProjectNotFound):
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	case errors.Is(err, user.ErrFileUploadFailed):
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	case errors.Is(err, user.ErrInvalidFileType):
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	case errors.Is(err, user.ErrFileSizeTooLarge):
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	default:
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
}
