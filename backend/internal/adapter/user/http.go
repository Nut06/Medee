package useradapter

import (
	userapp "backend/internal/application/userapp"
	"backend/internal/domain/domain"

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
	req := c.Query("id")
	res, err := h.usecase.GetProfile(c.Context(), req)
	if err != nil {
		return err
	}
	return c.JSON(res)
}

func (h *HTTPHandler) UploadAvatar(c *fiber.Ctx) error {
	userId := c.Query("id")
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
	userId := c.Query("id")
	var req userapp.UpdateProfileCommand
	if err := c.BodyParser(&req); err != nil {
		return err
	}
	user := &domain.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
	}
	res, err := h.usecase.UpdateProfile(c.Context(), userId, user)
	if err != nil {
		return err
	}
	return c.JSON(res)
}
