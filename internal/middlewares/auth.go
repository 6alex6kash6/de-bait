package middlewares

import (
	"fmt"
	"strings"

	"github.com/de-bait/internal/db/repository"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func Auth(ur repository.UserRepository) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		authorization := c.Get("Authorization")
		token := strings.TrimPrefix(authorization, "Bearer ")

		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "message": "You are not logged in"})
		}

		tokenData, err := jwt.Parse(token, func(jwtToken *jwt.Token) (interface{}, error) {
			if _, ok := jwtToken.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %s", jwtToken.Header["alg"])
			}
			return []byte("secret"), nil
		})

		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "message": fmt.Sprintf("invalidate token: %v", err)})
		}
		claims, ok := tokenData.Claims.(jwt.MapClaims)
		if !ok || !tokenData.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "message": "invalid token claim"})

		}
		id := int(claims["sub"].(float64))
		u, err := ur.FindOne(id)

		if err != nil {
			return fmt.Errorf("failed finding user: %w", err)
		}

		c.Locals("user", u)

		return c.Next()
	}
}
