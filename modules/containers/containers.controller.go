package containers

import (
	"strconv"

	"deployer.com/libs"
	"deployer.com/modules/auth/guards"
	"deployer.com/modules/containers/dto"
	"github.com/gofiber/fiber/v2"
)

type ContainersController struct {
	containersService *ContainersService
	router            *fiber.Router
}

func NewContainersController(router *fiber.Router, containersService *ContainersService) *ContainersController {
	return &ContainersController{router: router, containersService: containersService}
}

func (c *ContainersController) RegisterRoutes(router *fiber.Router) {
	(*c.router).Get("/", guards.JwtGuard, c.GetContainers)
	(*c.router).Get("/:id", guards.JwtGuard, c.GetContainer)
	(*c.router).Post("/", guards.JwtGuard, c.CreateContainer)
	(*c.router).Patch("/:id", guards.JwtGuard, c.UpdateContainer)
	(*c.router).Delete("/:id", guards.JwtGuard, c.DeleteContainer)
}

func (c *ContainersController) GetContainers(ctx *fiber.Ctx) error {
	userClaims := ctx.Locals("user").(*libs.UserClaims)
	containers, err := c.containersService.GetContainers(uint(userClaims.UserID), userClaims.IV)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(containers)
}

func (c *ContainersController) GetContainer(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	userClaims := ctx.Locals("user").(*libs.UserClaims)
	container, err := c.containersService.GetContainer(uint(id), uint(userClaims.UserID), userClaims.IV)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(container)
}

func (c *ContainersController) CreateContainer(ctx *fiber.Ctx) error {
	userClaims := ctx.Locals("user").(*libs.UserClaims)
	var body dto.CreateContainerDto
	if err := ctx.BodyParser(&body); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	if err := dto.ValidateCreateContainerDto(body); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	container, err := c.containersService.CreateContainer(uint(userClaims.UserID), body, userClaims.IV)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return ctx.Status(fiber.StatusCreated).JSON(container)
}

func (c *ContainersController) UpdateContainer(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	userClaims := ctx.Locals("user").(*libs.UserClaims)
	var body dto.UpdateContainerDto
	if err := ctx.BodyParser(&body); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	if err := dto.ValidateUpdateContainerDto(body); err != nil {
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
	container, err := c.containersService.UpdateContainer(uint(id), uint(userClaims.UserID), updates, userClaims.IV)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(container)
}

func (c *ContainersController) DeleteContainer(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	userClaims := ctx.Locals("user").(*libs.UserClaims)
	if err := c.containersService.DeleteContainer(uint(id), uint(userClaims.UserID)); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Container deleted successfully",
	})
}
