package useradapter

import (
	userapp "backend/internal/application/userapp"
	"backend/internal/domain/domain"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type HTTPHandler struct {
	usecase *userapp.Usecase
}

func NewHTTPHandler(db *gorm.DB) *HTTPHandler {
	repo := NewRepository(db)
	usecase := userapp.NewUsecase(repo)
	return &HTTPHandler{usecase: usecase}
}

func (h *HTTPHandler) GetUser(c *fiber.Ctx) error {
	userId := c.Locals("user_id").(string)
	res, err := h.usecase.GetProfile(c.Context(), userId)
	if err != nil {
		return err
	}
	return c.JSON(res)
}

func (h *HTTPHandler) UploadAvatar(c *fiber.Ctx) error {
	userId := c.Locals("user_id").(string)
	file, err := c.FormFile("avatar")
	if err != nil {
		return err
	}

	user, err := h.usecase.UploadAvatar(c.Context(), userId, file)
	if err != nil {
		return err
	}
	return c.JSON(user)
}

func (h *HTTPHandler) UpdateProfile(c *fiber.Ctx) error {
	val := c.Locals("user_id")
	if val == nil {
		return c.Status(500).JSON(fiber.Map{"error": "user_id missing from context"})
	}
	userId := val.(string)
	var req userapp.UpdateProfileCommand
	if err := c.BodyParser(&req); err != nil {
		return err
	}
	user := &domain.User{
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Email:       req.Email,
		LinkedInURL: &req.LinkedInURL,
		GitHubURL:   &req.GitHubURL,
		WebsiteURL:  &req.WebsiteURL,
	}
	fmt.Printf("user: %v from update profile http handler", user)
	res, err := h.usecase.UpdateProfile(c.Context(), userId, user)
	if err != nil {
		return err
	}
	return c.JSON(res)
}

func (h *HTTPHandler) DeleteAvatar(c *fiber.Ctx) error {
	userId := c.Locals("user_id").(string)
	user, err := h.usecase.DeleteAvatar(c.Context(), userId)
	if err != nil {
		return err
	}
	return c.JSON(user)
}

// Candidate Features

func (h *HTTPHandler) GetFullProfile(c *fiber.Ctx) error {
	userId := c.Locals("user_id").(string)
	fmt.Printf("id: %v  from get full profile http handler", userId)
	res, err := h.usecase.GetFullProfile(c.Context(), userId)
	if err != nil {
		return err
	}
	return c.JSON(res)
}

func (h *HTTPHandler) AddExperience(c *fiber.Ctx) error {
	userId := c.Locals("user_id").(string)
	var req domain.WorkExperience
	if err := c.BodyParser(&req); err != nil {
		return err
	}
	res, err := h.usecase.AddExperience(c.Context(), userId, &req)
	if err != nil {
		return err
	}
	return c.JSON(res)
}

func (h *HTTPHandler) UpdateExperience(c *fiber.Ctx) error {
	id := c.Params("experienceId")
	var req domain.WorkExperience
	if err := c.BodyParser(&req); err != nil {
		return err
	}
	res, err := h.usecase.UpdateExperience(c.Context(), id, &req)
	if err != nil {
		return err
	}
	return c.JSON(res)
}

func (h *HTTPHandler) DeleteExperience(c *fiber.Ctx) error {
	id := c.Params("experienceId")
	err := h.usecase.DeleteExperience(c.Context(), id)
	if err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusOK)
}

func (h *HTTPHandler) AddEducation(c *fiber.Ctx) error {
	userId := c.Locals("user_id").(string)
	var req domain.Education
	if err := c.BodyParser(&req); err != nil {
		return err
	}
	res, err := h.usecase.AddEducation(c.Context(), userId, &req)
	if err != nil {
		return err
	}
	return c.JSON(res)
}

func (h *HTTPHandler) UpdateEducation(c *fiber.Ctx) error {
	id := c.Params("educationId")
	var req domain.Education
	if err := c.BodyParser(&req); err != nil {
		return err
	}
	res, err := h.usecase.UpdateEducation(c.Context(), id, &req)
	if err != nil {
		return err
	}
	return c.JSON(res)
}

func (h *HTTPHandler) DeleteEducation(c *fiber.Ctx) error {
	id := c.Params("educationId")
	err := h.usecase.DeleteEducation(c.Context(), id)
	if err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusOK)
}

func (h *HTTPHandler) UpdateSkills(c *fiber.Ctx) error {
	userId := c.Locals("user_id").(string)
	var req struct {
		Skills []string `json:"skills"`
	}
	if err := c.BodyParser(&req); err != nil {
		return err
	}
	err := h.usecase.UpdateSkills(c.Context(), userId, req.Skills)
	if err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusOK)
}

func (h *HTTPHandler) AddProject(c *fiber.Ctx) error {
	userId := c.Locals("user_id").(string)
	var req domain.PortfolioItem
	if err := c.BodyParser(&req); err != nil {
		return err
	}
	res, err := h.usecase.AddProject(c.Context(), userId, &req)
	if err != nil {
		return err
	}
	return c.JSON(res)
}

func (h *HTTPHandler) UpdateProject(c *fiber.Ctx) error {
	id := c.Params("projectId")
	var req domain.PortfolioItem
	if err := c.BodyParser(&req); err != nil {
		return err
	}
	res, err := h.usecase.UpdateProject(c.Context(), id, &req)
	if err != nil {
		return err
	}
	return c.JSON(res)
}

func (h *HTTPHandler) DeleteProject(c *fiber.Ctx) error {
	id := c.Params("projectId")
	err := h.usecase.DeleteProject(c.Context(), id)
	if err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusOK)
}
