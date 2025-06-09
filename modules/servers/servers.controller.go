package servers

import (
	"strconv"

	"deployer.com/libs"
	"deployer.com/modules/auth/guards"
	"deployer.com/modules/servers/dto"
	"github.com/gofiber/fiber/v2"
)

type ServersController struct {
	serversService *ServersService
	router         *fiber.Router
}

func NewServersController(router *fiber.Router, serversService *ServersService) *ServersController {
	return &ServersController{router: router, serversService: serversService}
}

func (c *ServersController) RegisterRoutes(router *fiber.Router) {
	(*c.router).Get("/", guards.JwtGuard, c.GetServers)
	(*c.router).Get("/:id", guards.JwtGuard, c.GetServer)
	(*c.router).Post("/", guards.JwtGuard, c.CreateServer)
	(*c.router).Patch("/:id", guards.JwtGuard, c.UpdateServer)
	(*c.router).Delete("/:id", guards.JwtGuard, c.DeleteServer)
}

func (c *ServersController) GetServers(ctx *fiber.Ctx) error {
	userClaims := ctx.Locals("user").(*libs.UserClaims)
	servers, err := c.serversService.GetServers(uint(userClaims.UserID), userClaims.IV)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(servers)
}

func (c *ServersController) GetServer(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	userClaims := ctx.Locals("user").(*libs.UserClaims)
	server, err := c.serversService.GetServer(uint(id), uint(userClaims.UserID), userClaims.IV)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(server)
}

func (c *ServersController) CreateServer(ctx *fiber.Ctx) error {
	userClaims := ctx.Locals("user").(*libs.UserClaims)
	var body dto.CreateServerDto
	if err := ctx.BodyParser(&body); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	if err := dto.ValidateCreateServerDto(body); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	server, err := c.serversService.CreateServer(uint(userClaims.UserID), body, userClaims.IV)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return ctx.Status(fiber.StatusCreated).JSON(server)
}

func (c *ServersController) UpdateServer(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	userClaims := ctx.Locals("user").(*libs.UserClaims)
	var body dto.UpdateServerDto
	if err := ctx.BodyParser(&body); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	if err := dto.ValidateUpdateServerDto(body); err != nil {
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
	server, err := c.serversService.UpdateServer(uint(id), uint(userClaims.UserID), updates, userClaims.IV)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(server)
}

func (c *ServersController) DeleteServer(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	userClaims := ctx.Locals("user").(*libs.UserClaims)
	if err := c.serversService.DeleteServer(uint(id), uint(userClaims.UserID)); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Secret deleted successfully",
	})
}
