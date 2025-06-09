package libs

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func ExtractRefreshToken(ctx *fiber.Ctx) (string, error) {
	bearerToken := ctx.Request().Header.Peek("Authorization")
	if string(bearerToken) == "" {
		return "", errors.New("unauthorized")
	}
	return strings.TrimPrefix(string(bearerToken), "Refresh "), nil
}

func ExtractBearerToken(ctx *fiber.Ctx) (string, error) {
	bearerToken := ctx.Request().Header.Peek("Authorization")
	if string(bearerToken) == "" {
		return "", errors.New("unauthorized")
	}
	return strings.TrimPrefix(string(bearerToken), "Bearer "), nil
}
