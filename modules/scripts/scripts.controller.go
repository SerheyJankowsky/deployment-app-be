package scripts

import (
	"strconv"

	"deployer.com/libs"
	"deployer.com/modules/auth/guards"
	"deployer.com/modules/scripts/dto"
	"github.com/gofiber/fiber/v2"
)

type ScriptsController struct {
	scriptsService *ScriptsService
	router         *fiber.Router
}

func NewScriptsController(router *fiber.Router, scriptsService *ScriptsService) *ScriptsController {
	return &ScriptsController{router: router, scriptsService: scriptsService}
}

func (c *ScriptsController) RegisterRoutes(router *fiber.Router) {
	(*c.router).Get("/", guards.JwtGuard, c.GetScripts)
	(*c.router).Get("/:id", guards.JwtGuard, c.GetScript)
	(*c.router).Post("/", guards.JwtGuard, c.CreateScript)
	(*c.router).Patch("/:id", guards.JwtGuard, c.UpdateScript)
	(*c.router).Delete("/:id", guards.JwtGuard, c.DeleteScript)
}

func (c *ScriptsController) GetScripts(ctx *fiber.Ctx) error {
	userClaims := ctx.Locals("user").(*libs.UserClaims)
	scripts, err := c.scriptsService.GetScripts(uint(userClaims.UserID), userClaims.IV)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(scripts)
}

func (c *ScriptsController) GetScript(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	userClaims := ctx.Locals("user").(*libs.UserClaims)
	script, err := c.scriptsService.GetScript(uint(id), uint(userClaims.UserID), userClaims.IV)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(script)
}

func (c *ScriptsController) CreateScript(ctx *fiber.Ctx) error {
	userClaims := ctx.Locals("user").(*libs.UserClaims)
	var body dto.CreateScriptDto
	if err := ctx.BodyParser(&body); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	if err := dto.ValidateCreateScriptDto(body); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	script, err := c.scriptsService.CreateScript(uint(userClaims.UserID), body, userClaims.IV)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return ctx.Status(fiber.StatusCreated).JSON(script)
}

func (c *ScriptsController) UpdateScript(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	userClaims := ctx.Locals("user").(*libs.UserClaims)
	var body dto.UpdateScriptDto
	if err := ctx.BodyParser(&body); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	if err := dto.ValidateUpdateScriptDto(body); err != nil {
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
	script, err := c.scriptsService.UpdateScript(uint(id), uint(userClaims.UserID), updates, userClaims.IV)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(script)
}

func (c *ScriptsController) DeleteScript(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	userClaims := ctx.Locals("user").(*libs.UserClaims)
	if err := c.scriptsService.DeleteScript(uint(id), uint(userClaims.UserID)); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Script deleted successfully",
	})
}
