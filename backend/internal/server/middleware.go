package server

import (
	authadapter "backend/internal/adapter/auth"
	"backend/internal/domain/auth"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(c *fiber.Ctx) error {
	tokenString := c.Cookies("access_token")
	fmt.Println("Form auth middleware")
	if tokenString == "" {
		return fiber.NewError(fiber.StatusUnauthorized, auth.ErrInvalidToken.Error())
	}

	// Parse token
	claims, err := authadapter.ParseToken(tokenString)
	if err != nil {
		fmt.Println("Error parsing token:", err)
		return fiber.NewError(fiber.StatusUnauthorized, auth.ErrInvalidToken.Error())
	}

	// Set user_id to context
	claimsMap := *claims
	if sub, ok := claimsMap["sub"].(string); ok {
		c.Locals("user_id", sub)
	} else {
		return fiber.NewError(fiber.StatusUnauthorized, auth.ErrInvalidToken.Error())
	}

	return c.Next()
}
