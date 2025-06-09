package guards

import (
	"deployer.com/libs"
	"github.com/gofiber/fiber/v2"
)

func JwtGuard(ctx *fiber.Ctx) error {
	rfToken, err := libs.ExtractBearerToken(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
			"error":   err.Error(),
		})
	}
	claims, err := libs.ParseAccessToken(rfToken)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
			"error":   err.Error(),
		})
	}
	ctx.Locals("user", claims)
	return ctx.Next()
}
