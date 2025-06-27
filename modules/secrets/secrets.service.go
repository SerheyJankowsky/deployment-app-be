package secrets

import (
	"strings"
	"time"

	"deployer.com/libs"
	"deployer.com/modules/secrets/dto"
	"gorm.io/gorm"
)

type SecretsService struct {
	db                *gorm.DB
	encryptionService *libs.EncryptionService
}

type SecretResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewSecretsService(db *gorm.DB) *SecretsService {
	return &SecretsService{db: db, encryptionService: libs.NewEncryptionService()}
}

func (s *SecretsService) GetSecrets(userId uint, iv string) ([]SecretResponse, error) {
	var secrets []Secret
	if err := s.db.Where("user_id = ?", userId).Select("id, name, content, created_at").Order("created_at DESC").Find(&secrets).Error; err != nil {
		return nil, err
	}
	result := make([]SecretResponse, len(secrets))
	for i, secret := range secrets {
		decoded, err := s.encryptionService.Decrypt(secret.Content, iv)
		if err != nil {
			return nil, err
		}
		result[i] = SecretResponse{
			ID:        secret.ID,
			Name:      secret.Name,
			Content:   decoded,
			CreatedAt: secret.CreatedAt,
			UpdatedAt: secret.UpdatedAt,
		}
	}
	return result, nil
}

func (s *SecretsService) GetSecret(id, userId uint, iv string) (SecretResponse, error) {
	var secret Secret
	if err := s.db.Where("id = ? AND user_id = ?", id, userId).First(&secret).Error; err != nil {
		return SecretResponse{}, err
	}
	decoded, err := s.encryptionService.Decrypt(secret.Content, iv)
	if err != nil {
		return SecretResponse{}, err
	}
	return SecretResponse{
		ID:        secret.ID,
		Name:      secret.Name,
		Content:   decoded,
		CreatedAt: secret.CreatedAt,
		UpdatedAt: secret.UpdatedAt,
	}, nil
}

func (s *SecretsService) CreateSecret(userId uint, dto dto.CreateSecretDto, iv string) (SecretResponse, error) {
	encrypted, err := s.encryptionService.Encrypt(dto.Content, iv)
	if err != nil {
		return SecretResponse{}, err
	}
	secret := Secret{
		Name:    dto.Name,
		Content: encrypted,
		UserID:  userId,
	}
	if err := s.db.Create(&secret).Error; err != nil {
		return SecretResponse{}, err
	}
	return SecretResponse{
		ID:        secret.ID,
		Name:      secret.Name,
		Content:   dto.Content,
		CreatedAt: secret.CreatedAt,
		UpdatedAt: secret.UpdatedAt,
	}, nil
}

func (s *SecretsService) UpdateSecret(id, userId uint, updates map[string]interface{}, iv string) (SecretResponse, error) {
	var secret Secret
	if err := s.db.Where("id = ? AND user_id = ?", id, userId).First(&secret).Error; err != nil {
		return SecretResponse{}, err
	}
	libs.SetStructFieldsFromMap(&secret, updates)
	if updates["content"] != nil {
		encrypted, err := s.encryptionService.Encrypt(secret.Content, iv)
		if err != nil {
			return SecretResponse{}, err
		}
		secret.Content = encrypted
	}
	if err := s.db.Save(&secret).Error; err != nil {
		return SecretResponse{}, err
	}
	decoded, err := s.encryptionService.Decrypt(secret.Content, iv)
	if err != nil {
		return SecretResponse{}, err
	}
	return SecretResponse{
		ID:        secret.ID,
		Name:      secret.Name,
		Content:   decoded,
		CreatedAt: secret.CreatedAt,
		UpdatedAt: secret.UpdatedAt,
	}, nil
}

func (s *SecretsService) DeleteSecret(id, userId uint) error {
	if err := s.db.Where("id = ? AND user_id = ?", id, userId).Delete(&Secret{}).Error; err != nil {
		return err
	}
	return nil
}

func (s *SecretsService) GetEnvMap(secret SecretResponse) map[string]string {
	envMap := make(map[string]string)
	lines := strings.Split(secret.Content, "\n")
	for _, line := range lines {
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			envMap[parts[0]] = parts[1]
		}
	}
	return envMap
}
