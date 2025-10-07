package middleware

import (
	"strings"

	"github.com/google/uuid"
	"github.com/kenziehh/cashflow-be/pkg/errx"
	"github.com/kenziehh/cashflow-be/pkg/jwt"

	"github.com/gofiber/fiber/v2"
)

func JWTAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return errx.ErrUnauthorized
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return errx.ErrUnauthorized
		}

		token := parts[1]
		claims, err := jwt.ValidateToken(token)
		if err != nil {
			return errx.ErrUnauthorized
		}

		userUUID, err := uuid.Parse(claims.UserID)
		if err != nil {
			return errx.NewUnauthorizedError("Invalid user ID format in token")
		}

		c.Locals("userID", userUUID)
		return c.Next()
	}
}
