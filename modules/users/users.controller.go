package users

import (
	"strconv"

	"deployer.com/libs"
	"deployer.com/modules/auth/guards"
	"deployer.com/modules/users/dto"
	"github.com/gofiber/fiber/v2"
)

type UsersController struct {
	usersService *UsersService
	router       *fiber.Router
}

func NewUsersController(router *fiber.Router, usersService *UsersService) *UsersController {
	return &UsersController{usersService: usersService, router: router}
}

func (c *UsersController) RegisterRoutes(router *fiber.Router) {
	// (*c.router).Get("/:id", guards.JwtGuard, c.GetUser)
	(*c.router).Patch("/", guards.JwtGuard, c.UpdateUser)
	(*c.router).Delete("/", guards.JwtGuard, c.DeleteUser)

}

func (c *UsersController) GetUser(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	user, err := c.usersService.GetUser(uint(id))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return ctx.JSON(user)
}

func (c *UsersController) UpdateUser(ctx *fiber.Ctx) error {
	userClaims := ctx.Locals("user").(*libs.UserClaims)
	var updateUser dto.UpdateUserDto
	updateUser.ID = uint(userClaims.UserID)
	if err := ctx.BodyParser(&updateUser); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	if err := dto.ValidateUpdateUser(updateUser); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	user, err := c.usersService.UpdateUser(&updateUser)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return ctx.JSON(user)
}

func (c *UsersController) DeleteUser(ctx *fiber.Ctx) error {
	userClaims := ctx.Locals("user").(*libs.UserClaims)
	if err := c.usersService.DeleteUser(uint(userClaims.UserID)); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User deleted successfully",
	})
}
