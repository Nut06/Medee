package useradapter

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	userapp "backend/internal/application/user"
)

type HTTPHandler struct {
	db       *gorm.DB
	usecase  *userapp.Usecase
	repo     *UserRepository
}

func NewHTTPHandler(db *gorm.DB) *HTTPHandler {
	repo := NewUserRepository(db)
	user := NewUserService(repo)
	usecase := userapp.NewUsecase(repo, user)
	return &HTTPHandler{db: db, usecase: usecase, repo: repo}
}

func (h *HTTPHandler) GetUser(c *fiber.Ctx) error {
	req := c.Query("id")
	res, err := h.usecase.GetProfile(c.Context(), req)
	if err != nil {
		return err
	}
	return c.JSON(res)
}
