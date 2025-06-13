package guards

import (
	"strings"

	"deployer.com/libs"
	"github.com/gofiber/fiber/v2"
)

// CombinedGuard allows both JWT and API key authentication
func CombinedGuard(apiKeyService ApiKeyService) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// Try JWT first
		authHeader := ctx.Get("Authorization")
		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			token := strings.TrimPrefix(authHeader, "Bearer ")

			// Check if it looks like a JWT (contains dots)
			if strings.Contains(token, ".") {
				// Try JWT validation
				claims, err := libs.ParseAccessToken(token)
				if err == nil {
					ctx.Locals("user", claims)
					ctx.Locals("auth_method", "jwt")
					return ctx.Next()
				}
			} else {
				// Try as API key
				user, err := apiKeyService.GetUserByApiKey(token)
				if err == nil {
					ctx.Locals("user", user)
					ctx.Locals("auth_method", "api_key")
					return ctx.Next()
				}
			}
		}

		// Try X-API-Key header
		apiKey := ctx.Get("X-API-Key")
		if apiKey != "" {
			user, err := apiKeyService.GetUserByApiKey(apiKey)
			if err == nil {
				ctx.Locals("user", user)
				ctx.Locals("auth_method", "api_key")
				return ctx.Next()
			}
		}

		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
			"error":   "Valid JWT token or API key required",
		})
	}
}
