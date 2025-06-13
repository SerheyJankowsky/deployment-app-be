package dto

import (
	"deployer.com/modules/deployments"
	"github.com/go-playground/validator/v10"
)

type ProjectDeployments struct {
	Deployment   deployments.Deployment `json:"deployment"`
	DeploymentID uint                   `json:"deployment_id"`
	Order        int                    `json:"order"`
}

type CreateProjectDto struct {
	Name               string               `json:"name" validate:"required,min=1,max=255"`
	ProjectDeployments []ProjectDeployments `json:"project_deployments" validate:"omitempty,dive"`
}

func ValidateCreateProjectDto(dto CreateProjectDto) error {
	validate := validator.New()

	// Register custom validation for complex types if needed
	validate.RegisterValidation("valid_project_deployment", validateProjectDeployment)

	return validate.Struct(dto)
}

// Custom validation function for project deployments
func validateProjectDeployment(fl validator.FieldLevel) bool {
	projectDeployment, ok := fl.Field().Interface().(ProjectDeployments)
	if !ok {
		return false
	}
	// Check if deployment ID is provided and valid
	return projectDeployment.DeploymentID > 0
}

// Helper method to convert IDs to actual objects (to be used in service layer)
func (dto *CreateProjectDto) UseIDsOnly() bool {
	// Return true if any ID fields are populated, indicating client prefers ID-based approach
	return len(dto.ProjectDeployments) > 0
}

// Helper method to check if using full objects
func (dto *CreateProjectDto) UseFullObjects() bool {
	// Return true if any object fields are populated
	return len(dto.ProjectDeployments) > 0
}
