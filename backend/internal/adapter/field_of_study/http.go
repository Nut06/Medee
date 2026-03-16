package field_of_study_adapter

import (
	"backend/internal/application/fieldapp"

	"github.com/gofiber/fiber/v3"
)

type HTTPHandler struct {
	usecase *fieldapp.Usecase
}

func NewHTTPHandler(usecase *fieldapp.Usecase) *HTTPHandler {
	return &HTTPHandler{usecase: usecase}
}

func (h *HTTPHandler) SearchFieldOfStudies(c fiber.Ctx) error {
	query := c.Query("q")
	if query == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Query parameter 'q' is required",
		})
	}

	level := c.Query("level", "all") // "all", "broad", "detailed"

	fields, err := h.usecase.SearchFieldOfStudies(c.Context(), query, level)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Limit to 50 results for autocomplete
	if len(fields) > 50 {
		fields = fields[:50]
	}

	return c.JSON(fiber.Map{
		"results": fields,
		"count":   len(fields),
	})
}

func (h *HTTPHandler) GetFieldOfStudy(c fiber.Ctx) error {
	id := c.Params("id")

	field, err := h.usecase.GetFieldOfStudyById(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Field of study not found",
		})
	}

	return c.JSON(field)
}
