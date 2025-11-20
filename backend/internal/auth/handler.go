package auth

import (
	"os"

	"github.com/gofiber/fiber/v2"
)

var INVALID_REQUEST_BODY string = "Invalid request Body"
var NO_REFRESH_TOKEN string = "No refresh token"

func RegisterHandler(c *fiber.Ctx) error {
	var req *RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, INVALID_REQUEST_BODY)
	}

	res, err := register(req)

	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	return c.Status(fiber.StatusOK).JSON(res)
}

var fifteenMinutes = 60 * 15

func LoginHandler(c *fiber.Ctx) error {
	var Https bool
	if os.Getenv("HTTPS") == "true" {
		Https = true
	}
	Https = false
	var req *LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, INVALID_REQUEST_BODY)
	}

	res, token, err := login(req)

	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}
	
	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    token.AccessToken,
		HTTPOnly: true,
		Secure:   Https,
		SameSite: "Strict",
		MaxAge:   fifteenMinutes,
	})

	c.Cookie(
		&fiber.Cookie{
			Name:     "refresh_token",
			Value:    token.RefreshToken,
			HTTPOnly: true,
			Secure:   Https,
			SameSite: "Strict",
			MaxAge:   fifteenMinutes,
		})
	return c.Status(fiber.StatusOK).JSON(res)
}

func LogoutHandler(c *fiber.Ctx) error {
	refreshToken := c.Cookies("refresh_token")
	if refreshToken == "" {
		return fiber.NewError(fiber.StatusBadRequest, NO_REFRESH_TOKEN)
	}
	err := logout(refreshToken)

	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to del rf_token")
	}
	
	c.ClearCookie("access_token")
	c.ClearCookie("refresh_token")
	return c.JSON(fiber.Map{
		"message":"Logout success",
	})
}
