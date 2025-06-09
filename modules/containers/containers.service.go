package containers

import (
	"time"

	"deployer.com/libs"
	"deployer.com/modules/containers/dto"
	"gorm.io/gorm"
)

type ContainersService struct {
	db                *gorm.DB
	encryptionService *libs.EncryptionService
}

type ContainerResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Registry  string    `json:"registry"`
	Image     string    `json:"image"`
	Tag       string    `json:"tag"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	SecretKey string    `json:"secret_key"`
	Params    string    `json:"params"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewContainersService(db *gorm.DB) *ContainersService {
	return &ContainersService{db: db, encryptionService: libs.NewEncryptionService()}
}

func (s *ContainersService) GetContainers(userId uint, iv string) ([]ContainerResponse, error) {
	var containers []Container
	if err := s.db.Where("user_id = ?", userId).Select("id, name, registry, image, tag, username, password, secret_key, params, created_at").Order("created_at DESC").Find(&containers).Error; err != nil {
		return nil, err
	}
	result := make([]ContainerResponse, len(containers))
	for i, container := range containers {
		decodedPassword, err := s.encryptionService.Decrypt(container.Password, iv)
		if err != nil {
			return nil, err
		}
		decodedSecretKey, err := s.encryptionService.Decrypt(container.SecretKey, iv)
		if err != nil {
			return nil, err
		}
		result[i] = ContainerResponse{
			ID:        container.ID,
			Name:      container.Name,
			Registry:  container.Registry,
			Image:     container.Image,
			Tag:       container.Tag,
			Username:  container.Username,
			Password:  decodedPassword,
			SecretKey: decodedSecretKey,
			Params:    container.Params,
			CreatedAt: container.CreatedAt,
			UpdatedAt: container.UpdatedAt,
		}
	}
	return result, nil
}

func (s *ContainersService) GetContainer(id, userId uint, iv string) (ContainerResponse, error) {
	var container Container
	if err := s.db.Where("id = ? AND user_id = ?", id, userId).First(&container).Error; err != nil {
		return ContainerResponse{}, err
	}
	decodedPassword, err := s.encryptionService.Decrypt(container.Password, iv)
	if err != nil {
		return ContainerResponse{}, err
	}
	decodedSecretKey, err := s.encryptionService.Decrypt(container.SecretKey, iv)
	if err != nil {
		return ContainerResponse{}, err
	}
	return ContainerResponse{
		ID:        container.ID,
		Name:      container.Name,
		Registry:  container.Registry,
		Image:     container.Image,
		Tag:       container.Tag,
		Username:  container.Username,
		Password:  decodedPassword,
		SecretKey: decodedSecretKey,
		Params:    container.Params,
		CreatedAt: container.CreatedAt,
		UpdatedAt: container.UpdatedAt,
	}, nil
}

func (s *ContainersService) CreateContainer(userId uint, dto dto.CreateContainerDto, iv string) (ContainerResponse, error) {
	encryptedPassword, err := s.encryptionService.Encrypt(dto.Password, iv)
	if err != nil {
		return ContainerResponse{}, err
	}
	encryptedSecretKey, err := s.encryptionService.Encrypt(dto.SecretKey, iv)
	if err != nil {
		return ContainerResponse{}, err
	}
	container := Container{
		Name:      dto.Name,
		Registry:  dto.Registry,
		Image:     dto.Image,
		Tag:       dto.Tag,
		Username:  dto.Username,
		Password:  encryptedPassword,
		SecretKey: encryptedSecretKey,
		Params:    dto.Params,
		UserID:    userId,
	}
	if err := s.db.Create(&container).Error; err != nil {
		return ContainerResponse{}, err
	}
	return ContainerResponse{
		ID:        container.ID,
		Name:      container.Name,
		Registry:  container.Registry,
		Image:     container.Image,
		Tag:       container.Tag,
		Username:  container.Username,
		Password:  dto.Password,
		SecretKey: encryptedSecretKey,
		Params:    dto.Params,
		CreatedAt: container.CreatedAt,
		UpdatedAt: container.UpdatedAt,
	}, nil
}

func (s *ContainersService) UpdateContainer(id, userId uint, updates map[string]interface{}, iv string) (ContainerResponse, error) {
	var container Container
	if err := s.db.Where("id = ? AND user_id = ?", id, userId).First(&container).Error; err != nil {
		return ContainerResponse{}, err
	}
	libs.SetStructFieldsFromMap(&container, updates)
	if updates["password"] != nil {
		encrypted, err := s.encryptionService.Encrypt(container.Password, iv)
		if err != nil {
			return ContainerResponse{}, err
		}
		container.Password = encrypted
	}
	if updates["secret_key"] != nil {
		encrypted, err := s.encryptionService.Encrypt(container.SecretKey, iv)
		if err != nil {
			return ContainerResponse{}, err
		}
		container.SecretKey = encrypted
	}
	if err := s.db.Save(&container).Error; err != nil {
		return ContainerResponse{}, err
	}
	decodedPassword, err := s.encryptionService.Decrypt(container.Password, iv)
	if err != nil {
		return ContainerResponse{}, err
	}
	decodedSecretKey, err := s.encryptionService.Decrypt(container.SecretKey, iv)
	if err != nil {
		return ContainerResponse{}, err
	}
	return ContainerResponse{
		ID:        container.ID,
		Name:      container.Name,
		Registry:  container.Registry,
		Image:     container.Image,
		Tag:       container.Tag,
		Username:  container.Username,
		Password:  decodedPassword,
		SecretKey: decodedSecretKey,
		Params:    container.Params,
		CreatedAt: container.CreatedAt,
		UpdatedAt: container.UpdatedAt,
	}, nil
}

func (s *ContainersService) DeleteContainer(id, userId uint) error {
	if err := s.db.Where("id = ? AND user_id = ?", id, userId).Delete(&Container{}).Error; err != nil {
		return err
	}
	return nil
}
