package middleware

import (
	"github.com/gofiber/fiber/v2"
)

const HEADER_AUTHORIZATION_KEY = "authorization"

func (m *HTTPMiddleware) AuthorizationMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authorization := c.Get(HEADER_AUTHORIZATION_KEY)
		if authorization != m.authorizationKey {
			return c.SendStatus(fiber.StatusUnauthorized)
		}
		return c.Next()
	}
}
