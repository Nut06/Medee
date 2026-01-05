package institute_adapter

import (
	"backend/internal/application/instituteapp"

	"github.com/gofiber/fiber/v2"
)

type HTTPHandler struct {
	usecase *instituteapp.Usecase
}

func NewHTTPHandler(usecase *instituteapp.Usecase) *HTTPHandler {
	return &HTTPHandler{usecase: usecase}
}

// SearchInstitutes handles GET /institutes/search?q={query}
func (h *HTTPHandler) SearchInstitutes(c *fiber.Ctx) error {
	query := c.Query("q")
	if len(query) < 2 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Query must be at least 2 characters",
		})
	}

	country := c.Query("country", "")

	institutes, err := h.usecase.SearchInstitutes(c.Context(), query, country)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"results": institutes,
		"count":   len(institutes),
	})
}

// GetInstitute handles GET /institutes/:id
func (h *HTTPHandler) GetInstitute(c *fiber.Ctx) error {
	id := c.Params("id")

	institute, err := h.usecase.GetInstituteById(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Institute not found",
		})
	}

	return c.JSON(institute)
}
