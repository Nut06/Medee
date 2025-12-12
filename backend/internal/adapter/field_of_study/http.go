package field_of_study_adapter

import (
	"backend/internal/application/fieldapp"

	"github.com/gofiber/fiber/v2"
)

type HTTPHandler struct {
	usecase *fieldapp.Usecase
}

func NewHTTPHandler(usecase *fieldapp.Usecase) *HTTPHandler {
	return &HTTPHandler{usecase: usecase}
}

func (h *HTTPHandler) SearchFieldOfStudies(c *fiber.Ctx) error {
	query := c.Query("q")
	if query == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Query parameter 'q' is required",
		})
	}

	fields, err := h.usecase.SearchFieldOfStudies(c.Context(), query)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"results": fields,
	})
}

func (h *HTTPHandler) GetFieldOfStudy(c *fiber.Ctx) error {
	id := c.Params("id")

	field, err := h.usecase.GetFieldOfStudyById(c.Context(), id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Field of study not found",
		})
	}

	return c.JSON(field)
}
