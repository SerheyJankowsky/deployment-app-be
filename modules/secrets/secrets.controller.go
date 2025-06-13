package secrets

import (
	"strconv"

	"deployer.com/libs"
	"deployer.com/modules/auth/guards"
	"deployer.com/modules/secrets/dto"
	"deployer.com/modules/users"
	"github.com/gofiber/fiber/v2"
)

type SecretsController struct {
	secretsService *SecretsService
	router         *fiber.Router
}

func NewSecretsController(router *fiber.Router, secretsService *SecretsService) *SecretsController {
	return &SecretsController{router: router, secretsService: secretsService}
}

func (c *SecretsController) RegisterRoutes(router *fiber.Router) {
	(*c.router).Get("/", guards.JwtGuard, c.GetSecrets)
	(*c.router).Get("/:id", guards.JwtGuard, c.GetSecret)
	(*c.router).Post("/", guards.JwtGuard, c.CreateSecret)
	(*c.router).Patch("/:id", guards.JwtGuard, c.UpdateSecret)
	(*c.router).Delete("/:id", guards.JwtGuard, c.DeleteSecret)
}

// RegisterApiKeyRoutes creates API key only protected routes for external integrations
func (c *SecretsController) RegisterApiKeyRoutes(router *fiber.Router, apiKeyService guards.ApiKeyService) {
	apiKeyGuard := guards.ApiKeyGuard(apiKeyService)
	combinedGuard := guards.CombinedGuard(apiKeyService)

	// API key only routes (for external integrations)
	(*c.router).Get("/", apiKeyGuard, c.GetSecretsApiKey)
	(*c.router).Get("/:id", apiKeyGuard, c.GetSecretApiKey)

	// Combined auth routes (accept both JWT and API key)
	(*c.router).Post("/", combinedGuard, c.CreateSecret)
}

func (c *SecretsController) GetSecrets(ctx *fiber.Ctx) error {
	userClaims := ctx.Locals("user").(*libs.UserClaims)
	secrets, err := c.secretsService.GetSecrets(uint(userClaims.UserID), userClaims.IV)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(secrets)
}

func (c *SecretsController) GetSecret(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	userClaims := ctx.Locals("user").(*libs.UserClaims)
	secret, err := c.secretsService.GetSecret(uint(id), uint(userClaims.UserID), userClaims.IV)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(secret)
}

func (c *SecretsController) CreateSecret(ctx *fiber.Ctx) error {
	userClaims := ctx.Locals("user").(*libs.UserClaims)
	var body dto.CreateSecretDto
	if err := ctx.BodyParser(&body); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	if err := dto.ValidateCreateSecretDto(body); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	secret, err := c.secretsService.CreateSecret(uint(userClaims.UserID), body, userClaims.IV)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return ctx.Status(fiber.StatusCreated).JSON(secret)
}

func (c *SecretsController) UpdateSecret(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	userClaims := ctx.Locals("user").(*libs.UserClaims)
	var body dto.UpdateSecretDto
	if err := ctx.BodyParser(&body); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	if err := dto.ValidateUpdateSecretDto(body); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	if !body.HasUpdates() {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No updates provided",
		})
	}
	updates, _ := body.GetUpdates()
	secret, err := c.secretsService.UpdateSecret(uint(id), uint(userClaims.UserID), updates, userClaims.IV)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(secret)
}

func (c *SecretsController) DeleteSecret(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	userClaims := ctx.Locals("user").(*libs.UserClaims)
	if err := c.secretsService.DeleteSecret(uint(id), uint(userClaims.UserID)); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Secret deleted successfully",
	})
}

// GetSecretsApiKey handles API key authenticated requests
func (c *SecretsController) GetSecretsApiKey(ctx *fiber.Ctx) error {
	// Get user from API key authentication
	userInterface := ctx.Locals("user")

	// Handle different authentication methods
	authMethod := ctx.Locals("auth_method")
	if authMethod == "api_key" {
		// Import users package for proper type casting
		if user, ok := userInterface.(users.User); ok {
			secrets, err := c.secretsService.GetSecrets(user.ID, user.IV)
			if err != nil {
				return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": err.Error(),
				})
			}
			return ctx.Status(fiber.StatusOK).JSON(secrets)
		}
	}

	return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"error": "Invalid user context",
	})
}

// GetSecretApiKey handles API key authenticated requests for single secret
func (c *SecretsController) GetSecretApiKey(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Get user from API key authentication
	userInterface := ctx.Locals("user")

	// Handle different authentication methods
	authMethod := ctx.Locals("auth_method")
	if authMethod == "api_key" {
		if user, ok := userInterface.(users.User); ok {
			secret, err := c.secretsService.GetSecret(uint(id), user.ID, user.IV)
			if err != nil {
				return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"error": err.Error(),
				})
			}
			return ctx.Status(fiber.StatusOK).JSON(secret)
		}
	}

	return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"error": "Invalid user context",
	})
}
