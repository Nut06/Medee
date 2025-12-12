package authadapter

import (
	useradapter "backend/internal/adapter/user"
	authapp "backend/internal/application/auth"
	"backend/internal/domain/auth"
	authport "backend/internal/port/auth"
	userport "backend/internal/port/user"
	"context"
	"errors"
	"os"
	"time"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
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
	thirtyMin = 30 * time.Minute
	sevenDays = 30 * 24 * time.Hour // Increased to 30 days
)

var validate = validator.New()

type HTTPHandler struct {
	uc       authport.AuthService
	userRepo userport.UserRepository
	cfg      *CookieConfig
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
		NewJWTService(os.Getenv("JWT_SECRET"), thirtyMin, sevenDays),
		repo,
	)
	cfg.AccessCookieName = "access_token"
	cfg.RefreshCookieName = "refresh_token"
	userRepo := useradapter.NewRepository(db)
	return &HTTPHandler{uc: usecase, userRepo: userRepo, cfg: cfg}
}

func (h *HTTPHandler) Refresh(c *fiber.Ctx) error {
	ctx := h.context(c)

	rt := c.Cookies(h.cfg.RefreshCookieName)
	if rt == "" {
		return h.handleError(auth.ErrInvalidRefreshToken)
	}

	res, tokens, err := h.uc.Refresh(ctx, authport.RefreshCommand{
		RefreshToken: rt,
	})

	if err != nil {
		return h.handleError(err)
	}

	h.setCookies(c, tokens)

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

func (h *HTTPHandler) Register(c *fiber.Ctx) error {
	var req auth.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
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

	h.setCookies(c, tokens)

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

func (h *HTTPHandler) Login(c *fiber.Ctx) error {
	ctx := h.context(c)

	var req auth.LoginRequest
	if err := c.BodyParser(&req); err != nil {
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

	h.setCookies(c, tokens)

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

func (h *HTTPHandler) Logout(c *fiber.Ctx) error {
	ctx := h.context(c)
	rt := c.Cookies(h.cfg.RefreshCookieName)
	if rt == "" {
		return h.handleError(auth.ErrInvalidRefreshToken)
	}
	if err := h.uc.Logout(ctx, rt); err != nil {
		return h.handleError(err)
	}
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

func (h *HTTPHandler) clearCookies(c *fiber.Ctx) {
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
		Path:     "/",
		Domain:   "", // empty means current domain
		HTTPOnly: true,
		Secure:   h.cfg.Secure,
		SameSite: h.cfg.SameSite,
		MaxAge:   -1,
	})
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
		Path:     "/",
		Domain:   "", // empty means current domain
		HTTPOnly: true,
		Secure:   h.cfg.Secure,
		SameSite: "Lax",
		MaxAge:   accessMaxAge,
	})

	c.Cookie(&fiber.Cookie{
		Name:     h.cfg.RefreshCookieName,
		Value:    tokens.RefreshToken,
		Path:     "/",
		Domain:   "", // empty means current domain
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
