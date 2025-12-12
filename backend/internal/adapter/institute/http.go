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
	if query == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Query parameter 'q' is required",
		})
	}

	institutes, err := h.usecase.SearchInstitutes(c.Context(), query)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"results": institutes,
	})
}

// GetInstitute handles GET /institutes/:id
func (h *HTTPHandler) GetInstitute(c *fiber.Ctx) error {
	id := c.Params("id")

	institute, err := h.usecase.GetInstituteById(c.Context(), id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Institute not found",
		})
	}

	return c.JSON(institute)
}
