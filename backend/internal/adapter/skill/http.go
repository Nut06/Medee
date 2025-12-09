package skill_adapter

import (
	"backend/internal/application/skillapp"
	"backend/internal/domain/domain"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type HTTPHandler struct {
	usecase *skillapp.Usecase
}

func NewHTTPHandler(db *gorm.DB) *HTTPHandler {
	repo := NewRepository(db)
	usecase := skillapp.NewUsecase(repo)
	return &HTTPHandler{usecase: usecase}
}

func (h *HTTPHandler) GetSkills(c *fiber.Ctx) error {
	query := c.Query("q")
	if query == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "query param 'q' is required"})
	}
	res, err := h.usecase.SearchSkills(c.Context(), query)
	if err != nil {
		return err
	}
	return c.JSON(res)
}

func (h *HTTPHandler) CreateSkill(c *fiber.Ctx) error {
	var req domain.Skill
	if err := c.BodyParser(&req); err != nil {
		return err
	}
	res, err := h.usecase.CreateSkill(c.Context(), &req)
	if err != nil {
		return err
	}
	return c.JSON(res)
}
