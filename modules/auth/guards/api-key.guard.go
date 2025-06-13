package guards

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

// ApiKeyService interface to avoid import cycles
type ApiKeyService interface {
	GetUserByApiKey(apiKey string) (interface{}, error)
}

func ApiKeyGuard(apiKeyService ApiKeyService) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// Extract API key from X-API-Key header
		apiKey := ctx.Get("API-Key")
		if apiKey == "" {
			// Also check Authorization header with "Bearer" prefix for API keys
			authHeader := ctx.Get("Authorization")
			if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
				apiKey = strings.TrimPrefix(authHeader, "Bearer ")
			}
		}

		if apiKey == "" {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "API key required",
				"error":   "Missing X-API-Key header",
			})
		}

		// Validate API key
		user, err := apiKeyService.GetUserByApiKey(apiKey)
		if err != nil {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Invalid API key",
				"error":   "API key not found or invalid",
			})
		}

		// Store user in context for use in handlers
		ctx.Locals("user", user)
		ctx.Locals("auth_method", "api_key")

		return ctx.Next()
	}
}
