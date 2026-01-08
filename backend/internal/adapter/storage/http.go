package storage

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type StorageHandler struct {
	storage *SupabaseStorage
}

func NewStorageHandler() *StorageHandler {
	return &StorageHandler{storage: NewSupabaseStorage()}
}

type SignedURLRequest struct {
	Filename    string `json:"filename"`
	ContentType string `json:"contentType"`
}

func (h *StorageHandler) GetSignedUploadURL(c *fiber.Ctx) error {
	var req SignedURLRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if req.Filename == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "filename is required"})
	}

	// Generate unique filename to prevent collisions
	uniqueName := fmt.Sprintf("%s_%s", uuid.New().String(), req.Filename)

	result, err := h.storage.CreateSignedUploadURL(uniqueName)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"uploadUrl": result.SignedURL,
		"path":      uniqueName,
		"publicUrl": h.storage.GetPublicURL(uniqueName),
	})
}
