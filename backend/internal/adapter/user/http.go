package useradapter

import (
	userapp "backend/internal/application/userapp"
	"backend/internal/domain/domain"
	user "backend/internal/domain/user"
	skilladapter "backend/internal/adapter/skill"
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type HTTPHandler struct {
	usecase *userapp.Usecase
}

func NewHTTPHandler(db *gorm.DB) *HTTPHandler {
	userRepo := NewRepository(db)
	skillRepo := skilladapter.NewRepository(db)
	usecase := userapp.NewUsecase(userRepo, skillRepo)
	return &HTTPHandler{usecase: usecase}
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
	var req domain.WorkExperience
	if err := c.BodyParser(&req); err != nil {
		return h.handleError(err)
	}
	res, err := h.usecase.AddExperience(c.Context(), userId, &req)
	if err != nil {
		return h.handleError(err)
	}
	return c.JSON(res)
}

func (h *HTTPHandler) UpdateExperience(c *fiber.Ctx) error {
	id := c.Params("experienceId")
	var req domain.WorkExperience
	if err := c.BodyParser(&req); err != nil {
		return h.handleError(err)
	}
	res, err := h.usecase.UpdateExperience(c.Context(), id, &req)
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
	var req domain.Education
	if err := c.BodyParser(&req); err != nil {
		return h.handleError(err)
	}
	res, err := h.usecase.AddEducation(c.Context(), userId, &req)
	if err != nil {
		return h.handleError(err)
	}
	return c.JSON(res)
}

func (h *HTTPHandler) UpdateEducation(c *fiber.Ctx) error {
	id := c.Params("educationId")
	var req domain.Education
	if err := c.BodyParser(&req); err != nil {
		return h.handleError(err)
	}
	res, err := h.usecase.UpdateEducation(c.Context(), id, &req)
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
