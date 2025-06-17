package projects

import (
	"errors"
	"time"

	"deployer.com/modules/projects/dto"
	"gorm.io/gorm"
)

type ProjectsService struct {
	db *gorm.DB
}

type DeploymentResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ProjectDeploymentResponse struct {
	ID         uint               `json:"id"`
	Deployment DeploymentResponse `json:"deployment"`
	Order      int                `json:"order"`
	Status     string             `json:"status"`
	Logs       string             `json:"logs"`
}

type ProjectResponse struct {
	ID                 uint                        `json:"id"`
	Name               string                      `json:"name"`
	ProjectDeployments []ProjectDeploymentResponse `json:"project_deployments"`
	CreatedAt          time.Time                   `json:"created_at"`
	UpdatedAt          time.Time                   `json:"updated_at"`
}

func NewProjectsService(db *gorm.DB) *ProjectsService {
	return &ProjectsService{db: db}
}

func (s *ProjectsService) GetProjects(userId uint) ([]ProjectResponse, error) {
	var projects []Project
	if err := s.db.Where("user_id = ?", userId).Preload("ProjectDeployments.Deployment").Order("created_at DESC").Find(&projects).Error; err != nil {
		return nil, err
	}

	result := make([]ProjectResponse, len(projects))
	for i, project := range projects {
		projectDeployments := make([]ProjectDeploymentResponse, len(project.ProjectDeployments))
		for j, pd := range project.ProjectDeployments {
			projectDeployments[j] = ProjectDeploymentResponse{
				ID: pd.ID,
				Deployment: DeploymentResponse{
					ID:        pd.Deployment.ID,
					Name:      pd.Deployment.Name,
					CreatedAt: pd.Deployment.CreatedAt,
					UpdatedAt: pd.Deployment.UpdatedAt,
				},
				Order:  pd.Order,
				Status: pd.Status,
				Logs:   pd.Logs,
			}
		}

		result[i] = ProjectResponse{
			ID:                 project.ID,
			Name:               project.Name,
			ProjectDeployments: projectDeployments,
			CreatedAt:          project.CreatedAt,
			UpdatedAt:          project.UpdatedAt,
		}
	}
	return result, nil
}

func (s *ProjectsService) GetProject(id, userId uint) (ProjectResponse, error) {
	var project Project
	if err := s.db.Where("id = ? AND user_id = ?", id, userId).Preload("ProjectDeployments.Deployment").First(&project).Error; err != nil {
		return ProjectResponse{}, err
	}

	projectDeployments := make([]ProjectDeploymentResponse, len(project.ProjectDeployments))
	for j, pd := range project.ProjectDeployments {
		projectDeployments[j] = ProjectDeploymentResponse{
			ID: pd.ID,
			Deployment: DeploymentResponse{
				ID:        pd.Deployment.ID,
				Name:      pd.Deployment.Name,
				CreatedAt: pd.Deployment.CreatedAt,
				UpdatedAt: pd.Deployment.UpdatedAt,
			},
			Order:  pd.Order,
			Status: pd.Status,
			Logs:   pd.Logs,
		}
	}

	return ProjectResponse{
		ID:                 project.ID,
		Name:               project.Name,
		ProjectDeployments: projectDeployments,
		CreatedAt:          project.CreatedAt,
		UpdatedAt:          project.UpdatedAt,
	}, nil
}

// validateUserDeployments checks that all deployment IDs exist and belong to the user
func (s *ProjectsService) validateUserDeployments(deploymentIDs []uint, userId uint) error {
	if len(deploymentIDs) == 0 {
		return nil
	}

	// First check if all deployments exist
	var existingCount int64
	if err := s.db.Table("deployments").Where("id IN ?", deploymentIDs).Count(&existingCount).Error; err != nil {
		return err
	}

	if int(existingCount) != len(deploymentIDs) {
		return errors.New("some deployment IDs do not exist")
	}

	// Then check if all existing deployments belong to the user
	var userCount int64
	if err := s.db.Table("deployments").Where("id IN ? AND user_id = ?", deploymentIDs, userId).Count(&userCount).Error; err != nil {
		return err
	}

	if int(userCount) != len(deploymentIDs) {
		return errors.New("some deployments do not belong to this user")
	}

	return nil
}

func (s *ProjectsService) CreateProject(userId uint, createDto dto.CreateProjectDto) (ProjectResponse, error) {
	// Validate that all deployments belong to the user
	if len(createDto.ProjectDeployments) > 0 {
		deploymentIDs := make([]uint, len(createDto.ProjectDeployments))
		for i, pd := range createDto.ProjectDeployments {
			deploymentIDs[i] = pd.DeploymentID
		}

		if err := s.validateUserDeployments(deploymentIDs, userId); err != nil {
			return ProjectResponse{}, err
		}
	}

	// Start a transaction
	tx := s.db.Begin()
	if tx.Error != nil {
		return ProjectResponse{}, tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create the project
	project := Project{
		Name:   createDto.Name,
		UserID: userId,
	}

	if err := tx.Create(&project).Error; err != nil {
		tx.Rollback()
		return ProjectResponse{}, err
	}

	// Create project deployments if provided
	projectDeployments := make([]ProjectDeployments, 0)
	for _, pd := range createDto.ProjectDeployments {
		projectDeployment := ProjectDeployments{
			ProjectID:    project.ID,
			DeploymentID: pd.DeploymentID,
			Order:        pd.Order,
			Status:       "pending", // Default status
		}
		if err := tx.Create(&projectDeployment).Error; err != nil {
			tx.Rollback()
			return ProjectResponse{}, err
		}
		projectDeployments = append(projectDeployments, projectDeployment)
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return ProjectResponse{}, err
	}

	// Return the created project
	return s.GetProject(project.ID, userId)
}

func (s *ProjectsService) UpdateProject(id, userId uint, updates map[string]interface{}) (ProjectResponse, error) {
	var project Project
	if err := s.db.Where("id = ? AND user_id = ?", id, userId).First(&project).Error; err != nil {
		return ProjectResponse{}, err
	}

	// Update only allowed fields
	allowedFields := map[string]bool{
		"name": true,
	}

	filteredUpdates := make(map[string]interface{})
	for key, value := range updates {
		if allowedFields[key] {
			filteredUpdates[key] = value
		}
	}

	if len(filteredUpdates) > 0 {
		if err := s.db.Model(&project).Updates(filteredUpdates).Error; err != nil {
			return ProjectResponse{}, err
		}
	}

	// Handle project deployments update if provided
	if projectDeployments, exists := updates["project_deployments"]; exists {
		if deployments, ok := projectDeployments.([]dto.ProjectDeployments); ok {
			// Validate that all deployments belong to the user
			if len(deployments) > 0 {
				deploymentIDs := make([]uint, len(deployments))
				for i, pd := range deployments {
					deploymentIDs[i] = pd.DeploymentID
				}

				if err := s.validateUserDeployments(deploymentIDs, userId); err != nil {
					return ProjectResponse{}, err
				}
			}

			// Start transaction for updating deployments
			tx := s.db.Begin()
			if tx.Error != nil {
				return ProjectResponse{}, tx.Error
			}

			// Get existing project deployments
			var existingDeployments []ProjectDeployments
			if err := tx.Where("project_id = ?", id).Find(&existingDeployments).Error; err != nil {
				tx.Rollback()
				return ProjectResponse{}, err
			}

			// Create maps for easier lookup
			existingMap := make(map[uint]*ProjectDeployments)
			for i := range existingDeployments {
				existingMap[existingDeployments[i].DeploymentID] = &existingDeployments[i]
			}

			newDeploymentIDs := make(map[uint]bool)

			// Process incoming deployments
			for _, pd := range deployments {
				newDeploymentIDs[pd.DeploymentID] = true

				if existing, found := existingMap[pd.DeploymentID]; found {
					// Update existing ProjectDeployment
					existing.Order = pd.Order
					// Update status if provided, otherwise keep existing
					if pd.Status != "" {
						existing.Status = pd.Status
					}
					if err := tx.Save(existing).Error; err != nil {
						tx.Rollback()
						return ProjectResponse{}, err
					}
				} else {
					// Create new ProjectDeployment
					status := "pending" // Default status for new ones
					if pd.Status != "" {
						status = pd.Status
					}
					projectDeployment := ProjectDeployments{
						ProjectID:    id,
						DeploymentID: pd.DeploymentID,
						Order:        pd.Order,
						Status:       status,
					}
					if err := tx.Create(&projectDeployment).Error; err != nil {
						tx.Rollback()
						return ProjectResponse{}, err
					}
				}
			}

			// Remove ProjectDeployments that are not in the new list
			for deploymentID, existing := range existingMap {
				if !newDeploymentIDs[deploymentID] {
					if err := tx.Delete(existing).Error; err != nil {
						tx.Rollback()
						return ProjectResponse{}, err
					}
				}
			}

			if err := tx.Commit().Error; err != nil {
				return ProjectResponse{}, err
			}
		}
	}

	return s.GetProject(id, userId)
}

func (s *ProjectsService) DeleteProject(id, userId uint) error {
	// Start transaction
	tx := s.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Delete project deployments first (foreign key constraint)
	if err := tx.Where("project_id = ?", id).Delete(&ProjectDeployments{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Delete the project
	if err := tx.Where("id = ? AND user_id = ?", id, userId).Delete(&Project{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
