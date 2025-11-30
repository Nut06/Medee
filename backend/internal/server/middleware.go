package server

import (
	authadapter "backend/internal/adapter/auth"
	"backend/internal/domain/auth"

	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(c *fiber.Ctx) error {
	tokenString := c.Cookies("access_token")
	if tokenString == "" {
		return fiber.NewError(fiber.StatusUnauthorized, auth.ErrInvalidToken.Error())
	}
	
	// Parse token
	claims, err := authadapter.ParseToken(tokenString)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, auth.ErrInvalidToken.Error())
	}

	// Set user_id to context
	c.Locals("user_id", claims.Subject)

	return c.Next()
}
