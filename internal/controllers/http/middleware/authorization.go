package middleware

import (
	"encoding/base64"
	"faceScanner/pkg/password"
	"faceScanner/pkg/validator"
	"github.com/gofiber/fiber/v2"
	"log/slog"
	"net/http"
	"strings"
)

const HEADER_AUTHORIZATION_KEY = "authorization"

func (m *HTTPMiddleware) AuthorizationMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authorization := c.Get(HEADER_AUTHORIZATION_KEY)

		if !strings.HasPrefix(authorization, "Basic ") {
			return c.Status(http.StatusUnauthorized).
				SendString("Invalid authorization scheme")
		}

		encodedCredentials := strings.TrimPrefix(authorization, "Basic ")
		decodedCredentials, err := base64.StdEncoding.DecodeString(encodedCredentials)
		if err != nil {
			slog.Error(err.Error())
			return c.Status(http.StatusUnauthorized).
				SendString("Invalid authorization header")
		}

		credentials := string(decodedCredentials)
		parts := strings.SplitN(credentials, ":", 2)
		if len(parts) != 2 {
			slog.Error("Invalid authorization credentials")
			return c.Status(http.StatusUnauthorized).
				SendString("Invalid authorization credentials")
		}

		if valid := validator.ValidateEmail(parts[0]); !valid {
			return c.Status(fiber.StatusBadRequest).
				SendString("Login must be an email")
		}

		hashedPassword, err := m.authUsecase.GetUserCredentials(c.Context(), parts[0])
		if err != nil {
			slog.Error(err.Error())
			return c.Status(http.StatusUnauthorized).
				SendString("Credentials not found")
		}

		if err = password.ComparePassword(hashedPassword, parts[1]); err != nil {
			return c.Status(http.StatusUnauthorized).
				SendString("Wrong password")
		}

		return c.Next()
	}
}
