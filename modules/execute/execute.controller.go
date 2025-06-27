package execute

import (
	"deployer.com/libs"
	"deployer.com/modules/auth/guards"
	"deployer.com/modules/execute/dto"
	"github.com/gofiber/fiber/v2"
)

type ExecuteController struct {
	executeService *ExecuteService
	router         *fiber.Router
}

func NewExecuteController(router *fiber.Router, executeService *ExecuteService) *ExecuteController {
	return &ExecuteController{router: router, executeService: executeService}
}

func (c *ExecuteController) RegisterExecuteRoutes(router *fiber.Router) {
	(*c.router).Post("/script", guards.JwtGuard, c.RunScript)

}

func (c *ExecuteController) RunScript(ctx *fiber.Ctx) error {
	userClaims := ctx.Locals("user").(*libs.UserClaims)
	var runScriptDto dto.RunScriptDto
	if err := ctx.BodyParser(&runScriptDto); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	if err := c.executeService.RunScript(runScriptDto.ScriptID, uint(userClaims.UserID), runScriptDto.ServerID, runScriptDto.EnvID, userClaims.IV, runScriptDto.LoadEnv); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Script executed successfully",
	})
}
