package authadapter

import (
	authapp "backend/internal/application/auth"
	"backend/internal/domain/auth"
	"context"
	"errors"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type CookieConfig struct {
	AccessCookieName  string
	RefreshCookieName string
	Secure            bool
	SameSite          string
	AccessMaxAge      int
	RefreshMaxAge     int
}

type HTTPHandler struct {
	uc  *authapp.Usecase
	cfg *CookieConfig
}

func NewHTTPHandler(db *gorm.DB) *HTTPHandler {
	repo := NewGormRepository(db)
	cfg := &CookieConfig{
		Secure:   os.Getenv("HTTPS") == "true",
		SameSite: "Strict",
	}
	usecase := authapp.NewUsecase(
		repo,
		NewBcryptHasher(0),
		NewJWTService(os.Getenv("JWT_SECRET"), 15*time.Minute, 7*24*time.Hour),
		repo,
	)
	cfg.AccessCookieName = "access_token"
	cfg.RefreshCookieName = "refresh_token"
	cfg.SameSite = "Strict"
	return &HTTPHandler{uc: usecase, cfg: cfg}
}

func (h *HTTPHandler) Register(c *fiber.Ctx) error {

	var req auth.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	res, tokens, err := h.uc.Register(c.Context(), authapp.RegisterCommand{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  req.Password,
	})
	if err != nil {
		return h.handleError(err)
	}

	h.setCookies(c, tokens)

	return c.Status(fiber.StatusOK).JSON(auth.RegisterResponse{
		ID:        res.ID,
		FirstName: res.FirstName,
		LastName:  res.LastName,
		Email:     res.Email,
	})
}

func (h *HTTPHandler) Login(c *fiber.Ctx) error {
	ctx := h.context(c)

	var req auth.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	res, tokens, err := h.uc.Login(ctx, authapp.LoginCommand{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return h.handleError(err)
	}

	h.setCookies(c, tokens)
	return c.Status(fiber.StatusOK).JSON(auth.LoginResponse{
		ID:        res.ID,
		FirstName: res.FirstName,
		LastName:  res.LastName,
		Email:     res.Email,
		Companies: res.Companies,
	})
}

func (h *HTTPHandler) Logout(c *fiber.Ctx) error {
	ctx := h.context(c)
	rt := c.Cookies(h.cfg.RefreshCookieName)
	if rt == "" {
		return fiber.NewError(fiber.StatusBadRequest, "no refresh token")
	}
	if err := h.uc.Logout(ctx, rt); err != nil {
		return h.handleError(err)
	}
	c.ClearCookie(h.cfg.AccessCookieName)
	c.ClearCookie(h.cfg.RefreshCookieName)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "logout success"})
}

func (h *HTTPHandler) handleError(err error) *fiber.Error {
	switch {
	case err == nil:
		return nil
	case errors.Is(err, auth.ErrEmailAlreadyUsed):
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	case errors.Is(err, auth.ErrInvalidCredential):
		return fiber.NewError(fiber.StatusUnauthorized, err.Error())
	case errors.Is(err, auth.ErrUserNotFound):
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	default:
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
}

func (h *HTTPHandler) setCookies(c *fiber.Ctx, tokens *auth.TokenPair) {
	accessMaxAge := h.cfg.AccessMaxAge
	if accessMaxAge == 0 && !tokens.AccessExpiresAt.IsZero() {
		accessMaxAge = int(time.Until(tokens.AccessExpiresAt).Seconds())
	}
	refreshMaxAge := h.cfg.RefreshMaxAge
	if refreshMaxAge == 0 && !tokens.RefreshExpiresAt.IsZero() {
		refreshMaxAge = int(time.Until(tokens.RefreshExpiresAt).Seconds())
	}

	c.Cookie(&fiber.Cookie{
		Name:     h.cfg.AccessCookieName,
		Value:    tokens.AccessToken,
		HTTPOnly: true,
		Secure:   h.cfg.Secure,
		SameSite: h.cfg.SameSite,
		MaxAge:   accessMaxAge,
	})

	c.Cookie(&fiber.Cookie{
		Name:     h.cfg.RefreshCookieName,
		Value:    tokens.RefreshToken,
		HTTPOnly: true,
		Secure:   h.cfg.Secure,
		SameSite: h.cfg.SameSite,
		MaxAge:   refreshMaxAge,
	})
}

func (h *HTTPHandler) context(c *fiber.Ctx) context.Context {
	if uc := c.UserContext(); uc != nil {
		return uc
	}
	return context.Background()
}
