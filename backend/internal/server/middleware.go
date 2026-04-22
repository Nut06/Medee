package server

import (
	authadapter "backend/internal/adapter/auth"
	"backend/internal/domain/auth"
	"fmt"

	jwtware "github.com/gofiber/contrib/v3/jwt"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/extractors"
	"github.com/golang-jwt/jwt/v5"
)

// func AuthMiddleware(jwtService *authadapter.JWTService) fiber.Handler {
// 	return func(c fiber.Ctx) error {
// 		fmt.Println("From auth middleware")

// 		tokenString := c.Cookies("access_token")
// 		if tokenString == "" {
// 			return fiber.NewError(fiber.StatusUnauthorized, auth.ErrInvalidToken.Error())
// 		}

// 		// Parse token
// 		claims, err := jwtService.ParseToken(tokenString)
// 		if err != nil {
// 			fmt.Println("Error parsing token:", err)
// 			return fiber.NewError(fiber.StatusUnauthorized, auth.ErrInvalidToken.Error())
// 		}

// 		// Set user_id to context
// 		claimsMap := *claims
// 		if sub, ok := claimsMap["sub"].(string); ok {
// 			c.Locals("user_id", sub)
// 		} else {
// 			return fiber.NewError(fiber.StatusUnauthorized, auth.ErrInvalidToken.Error())
// 		}

// 		return c.Next()
// 	}
// }

func AuthMiddleware(jwtService *authadapter.JWTService) fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: jwtService.GetSecret()},
		Extractor:  extractors.FromCookie("access_token"),
		ErrorHandler: func(c fiber.Ctx, err error) error {
			fmt.Println("JWT Middleware Error", err)
			return fiber.NewError(fiber.StatusUnauthorized, auth.ErrInvalidAccessToken.Error())
		},
		SuccessHandler: func(c fiber.Ctx) error {
			user := jwtware.FromContext(c)
			if user == nil {
				return fiber.NewError(fiber.StatusUnauthorized, auth.ErrInvalidAccessToken.Error())
			}
			claims, ok := user.Claims.(jwt.MapClaims)
			if !ok {
				return fiber.NewError(fiber.StatusUnauthorized, auth.ErrInvalidAccessToken.Error())
			}

			if sub, ok := claims["sub"].(string); ok {
				c.Locals("user_id", sub)
			} else {
				return fiber.NewError(fiber.StatusUnauthorized, auth.ErrInvalidAccessToken.Error())
			}
			return c.Next()
		},
	})
}
