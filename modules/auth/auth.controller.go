package auth

import (
	"fmt"

	"deployer.com/libs"
	"deployer.com/modules/auth/dto"
	"deployer.com/modules/auth/guards"
	"github.com/gofiber/fiber/v2"
)

type AuthController struct {
	authService *AuthService
}

func NewAuthController(authService *AuthService) *AuthController {
	return &AuthController{authService: authService}
}

func (c *AuthController) RegisterRoutes(router *fiber.Router) {
	(*router).Post("/login", c.Login)
	(*router).Post("/register", c.Register)
	(*router).Post("/refresh", c.RefreshToken)
	(*router).Get("/me", guards.JwtGuard, c.Me)
}

func (c *AuthController) Login(ctx *fiber.Ctx) error {
	var body dto.LoginDto
	if err := ctx.BodyParser(&body); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}
	if err := dto.ValidateLogin(body); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}
	response, err := c.authService.Login(body)
	if err != nil {
		fmt.Println(err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(response)
}

func (c *AuthController) Register(ctx *fiber.Ctx) error {
	var body dto.RegisterDto
	if err := ctx.BodyParser(&body); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}
	if err := dto.ValidateRegister(body); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}
	response, err := c.authService.Register(body)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
			"error":   err.Error(),
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(response)
}

func (c *AuthController) RefreshToken(ctx *fiber.Ctx) error {
	var body dto.RefreshTokenDto
	if err := ctx.BodyParser(&body); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}
	if err := dto.ValidateRefreshToken(body); err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
			"error":   err.Error(),
		})
	}
	response, err := c.authService.RefreshToken(body.RefreshToken)
	if err != nil {
		fmt.Println(err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
			"error":   err.Error(),
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(response)
}

func (c *AuthController) Me(ctx *fiber.Ctx) error {
	userClaims := ctx.Locals("user").(*libs.UserClaims)
	user, err := c.authService.Me(userClaims)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal server error",
			"error":   err.Error(),
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(user)
}
