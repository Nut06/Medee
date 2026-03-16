package server

import (
	authadapter "backend/internal/adapter/auth"
	"backend/internal/domain/auth"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v3"
)

func AuthMiddleware(jwtService *authadapter.JWTService) fiber.Handler {
	return func(c fiber.Ctx) error {
		fmt.Println("From auth middleware")
		
		tokenString := c.Get("Authorization")
		if strings.HasPrefix(tokenString, "Bearer ") {
			tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		}else if c.Cookies("refresh_token") != "" {
			tokenString = c.Cookies("refresh_token")
		}
		if tokenString == "" {
			return fiber.NewError(fiber.StatusUnauthorized, auth.ErrInvalidToken.Error())
		}


		// Parse token
		claims, err := jwtService.ParseToken(tokenString)
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
}