package middleware

import (
	"strings"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

// JWT için gizli anahtar. Bu service-user'daki ile AYNI OLMALI!
var jwtSecret = []byte("super-secret-key")

// AuthRequired, bir JWT'nin varlığını ve geçerliliğini kontrol eder.
func AuthRequired() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing authorization header"})
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid authorization header format"})
		}
		
		tokenString := parts[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Doğrulama metodu (algoritma) kontrolü
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.NewError(fiber.StatusUnauthorized, "unexpected signing method")
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token"})
		}
		
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token claims"})
		}

		// Token'dan gelen bilgileri bir sonraki handler'ın kullanabilmesi için context'e ekliyoruz.
		c.Locals("user_id", claims["user_id"])
		c.Locals("role", claims["role"])

		return c.Next()
	}
}

// RoleRequired, belirli bir role sahip olmayı gerektiren bir middleware'dir.
func RoleRequired(requiredRole string) fiber.Handler {
    return func(c *fiber.Ctx) error {
        role, ok := c.Locals("role").(string)
        if !ok || role != requiredRole {
            return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "permission denied"})
        }
        return c.Next()
    }
}