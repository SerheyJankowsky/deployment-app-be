package domains

import (
	"strconv"

	"deployer.com/libs"
	"deployer.com/modules/auth/guards"
	"deployer.com/modules/domains/dto"
	"github.com/gofiber/fiber/v2"
)

type SubDomainsController struct {
	subDomainsService *SubDomainsService
	router            *fiber.Router
}

func NewSubDomainsController(router *fiber.Router, subDomainsService *SubDomainsService) *SubDomainsController {
	return &SubDomainsController{router: router, subDomainsService: subDomainsService}
}

func (c *SubDomainsController) RegisterSubDomainsRoutes(router *fiber.Router) {
	(*c.router).Get("/:domainId", guards.JwtGuard, c.GetSubDomains)
	(*c.router).Get("/one/:id", guards.JwtGuard, c.GetSubDomain)
	(*c.router).Post("/", guards.JwtGuard, c.CreateSubDomain)
	(*c.router).Patch("/:id", guards.JwtGuard, c.UpdateSubDomain)
	(*c.router).Delete("/:id", guards.JwtGuard, c.DeleteSubDomain)
}

func (c *SubDomainsController) GetSubDomains(ctx *fiber.Ctx) error {
	domainId, err := strconv.ParseUint(ctx.Params("domainId"), 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	userClaims := ctx.Locals("user").(*libs.UserClaims)
	subDomains, err := c.subDomainsService.GetSubDomains(uint(userClaims.UserID), uint(domainId), userClaims.IV)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(subDomains)
}

func (c *SubDomainsController) GetSubDomain(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	userClaims := ctx.Locals("user").(*libs.UserClaims)
	subDomain, err := c.subDomainsService.GetSubDomain(uint(id), uint(userClaims.UserID), userClaims.IV)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(subDomain)
}

func (c *SubDomainsController) CreateSubDomain(ctx *fiber.Ctx) error {
	userClaims := ctx.Locals("user").(*libs.UserClaims)
	var body dto.CreateSubDomainDto
	if err := ctx.BodyParser(&body); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	if err := dto.ValidateCreateSubDomainDto(body); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	subDomain, err := c.subDomainsService.CreateSubDomain(uint(userClaims.UserID), body, userClaims.IV)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return ctx.Status(fiber.StatusCreated).JSON(subDomain)
}

func (c *SubDomainsController) UpdateSubDomain(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	userClaims := ctx.Locals("user").(*libs.UserClaims)
	var body dto.UpdateSubDomainDto
	if err := ctx.BodyParser(&body); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	if err := dto.ValidateUpdateSubDomainDto(body); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	if !body.HasUpdatesSubDomain() {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No updates provided",
		})
	}
	updates, _ := body.GetUpdatesSubDomain()
	subDomain, err := c.subDomainsService.UpdateSubDomain(uint(id), uint(userClaims.UserID), updates, userClaims.IV)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(subDomain)
}

func (c *SubDomainsController) DeleteSubDomain(ctx *fiber.Ctx) error {
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	userClaims := ctx.Locals("user").(*libs.UserClaims)
	if err := c.subDomainsService.DeleteSubDomain(uint(id), uint(userClaims.UserID)); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "SubDomain deleted successfully",
	})
}
