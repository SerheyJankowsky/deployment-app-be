package deployments

import (
	"strconv"

	"deployer.com/libs"
	"deployer.com/modules/auth/guards"
	"deployer.com/modules/deployments/dto"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type DeploymentsController struct {
	deploymentsService *DeploymentsService
	router             *fiber.Router
}

func NewDeploymentsController(router *fiber.Router, deploymentsService *DeploymentsService) *DeploymentsController {
	return &DeploymentsController{router: router, deploymentsService: deploymentsService}
}

func (c *DeploymentsController) RegisterRoutes(router *fiber.Router) {
	(*c.router).Get("/", guards.JwtGuard, c.GetDeployments)
	(*c.router).Get("/:id", guards.JwtGuard, c.GetDeployment)
	(*c.router).Post("/", guards.JwtGuard, c.CreateDeployment)
	(*c.router).Patch("/:id", guards.JwtGuard, c.UpdateDeployment)
	(*c.router).Delete("/:id", guards.JwtGuard, c.DeleteDeployment)
}

func (c *DeploymentsController) GetDeployments(ctx *fiber.Ctx) error {
	userClaims := ctx.Locals("user").(*libs.UserClaims)

	deployments, err := c.deploymentsService.GetDeployments(uint(userClaims.UserID), userClaims.IV)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to retrieve deployments",
			"details": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    deployments,
		"count":   len(deployments),
		"message": "Deployments retrieved successfully",
	})
}

func (c *DeploymentsController) GetDeployment(ctx *fiber.Ctx) error {
	// Parse and validate ID
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid deployment ID format",
		})
	}

	if id == 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Deployment ID must be greater than 0",
		})
	}

	userClaims := ctx.Locals("user").(*libs.UserClaims)

	deployment, err := c.deploymentsService.GetDeployment(uint(id), uint(userClaims.UserID), userClaims.IV)
	if err != nil {
		// Handle specific GORM errors
		if err == gorm.ErrRecordNotFound {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Deployment not found",
			})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to retrieve deployment",
			"details": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    deployment,
		"message": "Deployment retrieved successfully",
	})
}

func (c *DeploymentsController) CreateDeployment(ctx *fiber.Ctx) error {
	userClaims := ctx.Locals("user").(*libs.UserClaims)

	var body dto.CreateDeploymentDto
	if err := ctx.BodyParser(&body); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
	}

	// Validate the DTO
	if err := dto.ValidateCreateDeploymentDto(body); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Validation failed",
			"details": err.Error(),
		})
	}

	deployment, err := c.deploymentsService.CreateDeployment(uint(userClaims.UserID), body)
	if err != nil {
		// Check for specific database errors
		if err.Error() == "UNIQUE constraint failed" ||
			err.Error() == "duplicate key value violates unique constraint" {
			return ctx.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "A deployment with this name already exists",
			})
		}

		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to create deployment",
			"details": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data":    deployment,
		"message": "Deployment created successfully",
	})
}

func (c *DeploymentsController) UpdateDeployment(ctx *fiber.Ctx) error {
	// Parse and validate ID
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid deployment ID format",
		})
	}

	if id == 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Deployment ID must be greater than 0",
		})
	}

	userClaims := ctx.Locals("user").(*libs.UserClaims)

	var body dto.UpdateDeploymentDto
	if err := ctx.BodyParser(&body); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
	}

	// Validate the DTO
	if err := dto.ValidateUpdateDeploymentDto(body); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Validation failed",
			"details": err.Error(),
		})
	}

	// Check if there are any updates
	if !body.HasUpdates() {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No updates provided",
		})
	}

	updates, _ := body.GetUpdates()
	// if err != nil {
	// 	return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
	// 		"error":   "Failed to process updates",
	// 		"details": err.Error(),
	// 	})
	// }

	deployment, err := c.deploymentsService.UpdateDeployment(uint(id), uint(userClaims.UserID), updates)
	if err != nil {
		// Handle specific GORM errors
		if err == gorm.ErrRecordNotFound {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Deployment not found",
			})
		}

		// Check for constraint violations
		if err.Error() == "UNIQUE constraint failed" ||
			err.Error() == "duplicate key value violates unique constraint" {
			return ctx.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": "A deployment with this name already exists",
			})
		}

		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to update deployment",
			"details": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    deployment,
		"message": "Deployment updated successfully",
	})
}

func (c *DeploymentsController) DeleteDeployment(ctx *fiber.Ctx) error {
	// Parse and validate ID
	id, err := strconv.ParseUint(ctx.Params("id"), 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid deployment ID format",
		})
	}

	if id == 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Deployment ID must be greater than 0",
		})
	}

	userClaims := ctx.Locals("user").(*libs.UserClaims)

	if err := c.deploymentsService.DeleteDeployment(uint(id), uint(userClaims.UserID)); err != nil {
		// Handle specific GORM errors
		if err == gorm.ErrRecordNotFound {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Deployment not found",
			})
		}

		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to delete deployment",
			"details": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Deployment deleted successfully",
	})
}

// Additional helper endpoints for managing relationships

// Get deployments with minimal data (just ID and name) for dropdowns/selects
func (c *DeploymentsController) GetDeploymentsList(ctx *fiber.Ctx) error {
	userClaims := ctx.Locals("user").(*libs.UserClaims)

	deployments, err := c.deploymentsService.GetDeployments(uint(userClaims.UserID), userClaims.IV)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Return minimal data for lists/dropdowns
	type DeploymentListItem struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	}

	result := make([]DeploymentListItem, len(deployments))
	for i, deployment := range deployments {
		result[i] = DeploymentListItem{
			ID:   deployment.ID,
			Name: deployment.Name,
		}
	}

	return ctx.Status(fiber.StatusOK).JSON(result)
}

// Register the additional route if needed
func (c *DeploymentsController) RegisterAdditionalRoutes(router *fiber.Router) {
	(*c.router).Get("/list", guards.JwtGuard, c.GetDeploymentsList)
}
