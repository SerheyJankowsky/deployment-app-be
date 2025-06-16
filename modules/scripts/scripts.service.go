package scripts

import (
	"time"

	"deployer.com/libs"
	"deployer.com/modules/scripts/dto"
	"gorm.io/gorm"
)

type ScriptsService struct {
	db                *gorm.DB
	encryptionService *libs.EncryptionService
}

type ScriptResponse struct {
	ID          uint       `json:"id"`
	Name        string     `json:"name"`
	Script      string     `json:"script"`
	Description string     `json:"description"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	LastRunAt   *time.Time `json:"last_run_at"`
}

func NewScriptsService(db *gorm.DB) *ScriptsService {
	return &ScriptsService{db: db, encryptionService: libs.NewEncryptionService()}
}

func (s *ScriptsService) GetScripts(userId uint, iv string) ([]ScriptResponse, error) {
	var scripts []Script
	if err := s.db.Where("user_id = ?", userId).Select("id, name, script, created_at, updated_at").Order("created_at DESC").Find(&scripts).Error; err != nil {
		return nil, err
	}
	result := make([]ScriptResponse, len(scripts))
	for i, script := range scripts {
		decoded, err := s.encryptionService.Decrypt(script.Script, iv)
		if err != nil {
			return nil, err
		}
		result[i] = ScriptResponse{
			ID:          script.ID,
			Name:        script.Name,
			Script:      decoded,
			Description: script.Description,
			CreatedAt:   script.CreatedAt,
			UpdatedAt:   script.UpdatedAt,
			LastRunAt:   script.LastRunAt,
		}
	}
	return result, nil
}

func (s *ScriptsService) GetScript(id, userId uint, iv string) (ScriptResponse, error) {
	var script Script
	if err := s.db.Where("id = ? AND user_id = ?", id, userId).First(&script).Error; err != nil {
		return ScriptResponse{}, err
	}
	decoded, err := s.encryptionService.Decrypt(script.Script, iv)
	if err != nil {
		return ScriptResponse{}, err
	}
	return ScriptResponse{
		ID:          script.ID,
		Name:        script.Name,
		Script:      decoded,
		Description: script.Description,
		CreatedAt:   script.CreatedAt,
		UpdatedAt:   script.UpdatedAt,
		LastRunAt:   script.LastRunAt,
	}, nil
}

func (s *ScriptsService) CreateScript(userId uint, dto dto.CreateScriptDto, iv string) (ScriptResponse, error) {
	encrypted, err := s.encryptionService.Encrypt(dto.Script, iv)
	if err != nil {
		return ScriptResponse{}, err
	}
	script := Script{
		Name:        dto.Name,
		Script:      encrypted,
		Description: dto.Description,
		UserID:      userId,
	}
	if err := s.db.Create(&script).Error; err != nil {
		return ScriptResponse{}, err
	}
	return ScriptResponse{
		ID:          script.ID,
		Name:        script.Name,
		Script:      dto.Script,
		Description: dto.Description,
		CreatedAt:   script.CreatedAt,
		UpdatedAt:   script.UpdatedAt,
		LastRunAt:   script.LastRunAt,
	}, nil
}

func (s *ScriptsService) UpdateScript(id, userId uint, updates map[string]interface{}, iv string) (ScriptResponse, error) {
	var script Script
	if err := s.db.Where("id = ? AND user_id = ?", id, userId).First(&script).Error; err != nil {
		return ScriptResponse{}, err
	}
	libs.SetStructFieldsFromMap(&script, updates)
	if updates["script"] != nil {
		encrypted, err := s.encryptionService.Encrypt(script.Script, iv)
		if err != nil {
			return ScriptResponse{}, err
		}
		script.Script = encrypted
	}
	if updates["description"] != nil {
		script.Description = updates["description"].(string)
	}
	if err := s.db.Save(&script).Error; err != nil {
		return ScriptResponse{}, err
	}
	decoded, err := s.encryptionService.Decrypt(script.Script, iv)
	if err != nil {
		return ScriptResponse{}, err
	}
	return ScriptResponse{
		ID:          script.ID,
		Name:        script.Name,
		Script:      decoded,
		Description: script.Description,
		CreatedAt:   script.CreatedAt,
		UpdatedAt:   script.UpdatedAt,
		LastRunAt:   script.LastRunAt,
	}, nil
}

func (s *ScriptsService) DeleteScript(id, userId uint) error {
	if err := s.db.Where("id = ? AND user_id = ?", id, userId).Delete(&Script{}).Error; err != nil {
		return err
	}
	return nil
}
