package authadapter

import (
	useradapter "backend/internal/adapter/user"
	authapp "backend/internal/application/auth"
	"backend/internal/domain/auth"
	authport "backend/internal/port/auth"
	userport "backend/internal/port/user"
	"backend/internal/utils"
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/csrf"
	"gorm.io/gorm"
)

type CookieConfig struct {
	AccessCookieName  string
	RefreshCookieName string
	Secure            bool
	SameSite          string // Fiber accepts string: "Lax", "Strict", "None"
	AccessMaxAge      int
	RefreshMaxAge     int
}

var (
	oneHour    = 1 * time.Hour
	thirtyDays = 30 * 24 * time.Hour // Increased to 30 days
)

var validate = validator.New()

type HTTPHandler struct {
	uc       authport.AuthService
	userRepo userport.UserRepository
	cfg      *CookieConfig
}

func NewHTTPHandler(db *gorm.DB, jwtService *JWTService) *HTTPHandler {
	repo := NewGormRepository(db)
	cfg := &CookieConfig{
		Secure:            os.Getenv("HTTPS") == "true",
		SameSite:          "Strict",
		AccessCookieName:  "access_token",
		RefreshCookieName: "refresh_token",
		AccessMaxAge:      int(utils.FiveSec),
		RefreshMaxAge:     int(utils.ThirtyDays),
	}

	usecase := authapp.NewUsecase(
		repo,
		NewBcryptHasher(0),
		jwtService,
	)

	userRepo := useradapter.NewRepository(db)
	return &HTTPHandler{uc: usecase, userRepo: userRepo, cfg: cfg}
}

func (h *HTTPHandler) GetCSRFToken(c fiber.Ctx) error {
	token := csrf.TokenFromContext(c)
	if token == "" {
		return fiber.NewError(fiber.StatusInternalServerError, "csrf token not available")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"csrfToken": token,
	})
}

func (h *HTTPHandler) Refresh(c fiber.Ctx) error {
	ctx := h.context(c)

	rt := c.Cookies(h.cfg.RefreshCookieName)
	if rt == "" {
		fmt.Println("Refresh token not found")
		return h.handleError(auth.ErrInvalidRefreshToken)
	}

	refreshcmd := authport.RefreshCommand{
		RefreshToken: rt,
	}

	res, tokens, err := h.uc.Refresh(ctx, refreshcmd)

	if err != nil {
		return h.handleError(err)
	}

	h.setRefreshCookie(c, tokens)
	h.setAccessCookie(c, tokens)

	// Fetch full user from database
	fullUser, err := h.userRepo.FindById(ctx, res.ID)
	if err != nil {
		return h.handleError(err)
	}

	return c.Status(fiber.StatusOK).JSON(auth.RefreshResponse{
		User:      *auth.ToUserResponse(fullUser),
		Companies: res.Companies,
	})
}

func (h *HTTPHandler) Register(c fiber.Ctx) error {
	var req auth.RegisterRequest
	if err := c.Bind().Body(&req); err != nil {
		return h.handleError(auth.ErrInvalidRequestBody)
	}

	if err := validate.Struct(&req); err != nil {
		return h.handleError(auth.ErrInvalidRequestBody)
	}

	res, tokens, err := h.uc.Register(c.Context(), authport.RegisterCommand{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Password:  req.Password,
	})

	if err != nil {
		return h.handleError(err)
	}

	h.setRefreshCookie(c, tokens)
	h.setAccessCookie(c, tokens)

	// Fetch full user from database
	ctx := h.context(c)
	fullUser, err := h.userRepo.FindById(ctx, res.ID)
	if err != nil {
		return h.handleError(err)
	}

	return c.Status(fiber.StatusOK).JSON(auth.RegisterResponse{
		User:      *auth.ToUserResponse(fullUser),
		Companies: []auth.Company{}, // Empty for new user
	})
}

func (h *HTTPHandler) Login(c fiber.Ctx) error {
	ctx := h.context(c)

	var req auth.LoginRequest
	if err := c.Bind().Body(&req); err != nil {
		return h.handleError(auth.ErrInvalidRequestBody)
	}

	if err := validate.Struct(&req); err != nil {
		return h.handleError(auth.ErrInvalidRequestBody)
	}

	res, tokens, err := h.uc.Login(ctx, authport.LoginCommand{
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		return h.handleError(err)
	}

	h.setRefreshCookie(c, tokens)
	h.setAccessCookie(c, tokens)

	// Fetch full user from database
	fullUser, err := h.userRepo.FindById(ctx, res.ID)
	if err != nil {
		return h.handleError(err)
	}

	return c.Status(fiber.StatusOK).JSON(auth.LoginResponse{
		User:      *auth.ToUserResponse(fullUser),
		Companies: res.Companies,
	})
}

func (h *HTTPHandler) Logout(c fiber.Ctx) error {
	h.clearCookies(c)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "logout success"})
}

func (h *HTTPHandler) handleError(err error) *fiber.Error {
	switch {
	case err == nil:
		return nil
	case errors.Is(err, auth.ErrInvalidRequestBody):
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	case errors.Is(err, auth.ErrEmailAlreadyUsed):
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	case errors.Is(err, auth.ErrInvalidCredential):
		return fiber.NewError(fiber.StatusUnauthorized, err.Error())
	case errors.Is(err, auth.ErrInvalidToken):
		return fiber.NewError(fiber.StatusUnauthorized, err.Error())
	case errors.Is(err, auth.ErrInvalidRefreshToken):
		return fiber.NewError(fiber.StatusUnauthorized, err.Error())
	case errors.Is(err, auth.ErrInvalidAccessToken):
		return fiber.NewError(fiber.StatusUnauthorized, err.Error())
	case errors.Is(err, auth.ErrUserNotFound):
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	case errors.Is(err, auth.ErrInvalidUser):
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	default:
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
}

func (h *HTTPHandler) clearCookies(c fiber.Ctx) {
	c.Cookie(&fiber.Cookie{
		Name:     h.cfg.AccessCookieName,
		Value:    "",
		Path:     "/",
		Domain:   "", // empty means current domain
		HTTPOnly: true,
		Secure:   h.cfg.Secure,
		SameSite: h.cfg.SameSite,
		MaxAge:   -1,
	})

	c.Cookie(&fiber.Cookie{
		Name:     h.cfg.RefreshCookieName,
		Value:    "",
		Path:     "/auth/refresh",
		Domain:   "", // empty means current domain
		HTTPOnly: true,
		Secure:   h.cfg.Secure,
		SameSite: h.cfg.SameSite,
		MaxAge:   -1,
	})
}

func (h *HTTPHandler) setRefreshCookie(c fiber.Ctx, tokens *auth.TokenPair) {
	refreshMaxAge := h.cfg.RefreshMaxAge

	c.Cookie(&fiber.Cookie{
		Name:     h.cfg.RefreshCookieName,
		Value:    tokens.RefreshToken,
		Path:     "/auth/refresh",
		Domain:   "", // empty means current domain
		HTTPOnly: true,
		Secure:   h.cfg.Secure,
		SameSite: h.cfg.SameSite,
		MaxAge:   refreshMaxAge,
	})
}

func (h *HTTPHandler) setAccessCookie(c fiber.Ctx, tokens *auth.TokenPair) {
	c.Cookie(&fiber.Cookie{
		Name:     h.cfg.AccessCookieName,
		Value:    tokens.AccessToken,
		Path:     "/",
		Domain:   "", // empty means current domain
		HTTPOnly: true,
		Secure:   h.cfg.Secure,
		SameSite: h.cfg.SameSite,
		MaxAge:   h.cfg.AccessMaxAge,
	})
}

func (h *HTTPHandler) context(c fiber.Ctx) context.Context {
	if uc := c.Context(); uc != nil {
		return uc
	}
	return context.Background()
}
