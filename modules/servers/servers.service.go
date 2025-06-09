package servers

import (
	"time"

	"deployer.com/libs"
	"deployer.com/modules/servers/dto"
	"gorm.io/gorm"
)

type ServersService struct {
	db                *gorm.DB
	encryptionService *libs.EncryptionService
}

type ServerResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Username  string    `json:"username"`
	Host      string    `json:"host"`
	Port      int       `json:"port"`
	SSHKey    string    `json:"ssh_key"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewServersService(db *gorm.DB) *ServersService {
	return &ServersService{db: db, encryptionService: libs.NewEncryptionService()}
}

func (s *ServersService) GetServers(userId uint, iv string) ([]ServerResponse, error) {
	var servers []Server
	if err := s.db.Where("user_id = ?", userId).Select("id, name,username, host, port, ssh_key, password, created_at").Order("created_at DESC").Find(&servers).Error; err != nil {
		return nil, err
	}
	result := make([]ServerResponse, len(servers))
	for i, server := range servers {
		var err error
		var decodedSSHKey string
		if server.SSHKey != nil {
			decodedSSHKey, err = s.encryptionService.Decrypt(*server.SSHKey, iv)
			if err != nil {
				return nil, err
			}
		} else {
			decodedSSHKey = ""
		}
		decodedPassword, err := s.encryptionService.Decrypt(server.Password, iv)
		if err != nil {
			return nil, err
		}
		result[i] = ServerResponse{
			ID:        server.ID,
			Name:      server.Name,
			Username:  server.Username,
			Host:      server.Host,
			Port:      server.Port,
			SSHKey:    decodedSSHKey,
			Password:  decodedPassword,
			CreatedAt: server.CreatedAt,
			UpdatedAt: server.UpdatedAt,
		}
	}
	return result, nil
}

func (s *ServersService) GetServer(id, userId uint, iv string) (ServerResponse, error) {
	var server Server
	if err := s.db.Where("id = ? AND user_id = ?", id, userId).First(&server).Error; err != nil {
		return ServerResponse{}, err
	}
	var err error
	var decodedSSHKey string
	if server.SSHKey != nil {
		decodedSSHKey, err = s.encryptionService.Decrypt(*server.SSHKey, iv)
		if err != nil {
			return ServerResponse{}, err
		}
	} else {
		decodedSSHKey = ""
	}
	decodedPassword, err := s.encryptionService.Decrypt(server.Password, iv)
	if err != nil {
		return ServerResponse{}, err
	}
	return ServerResponse{
		ID:        server.ID,
		Name:      server.Name,
		Username:  server.Username,
		Host:      server.Host,
		Port:      server.Port,
		SSHKey:    decodedSSHKey,
		Password:  decodedPassword,
		CreatedAt: server.CreatedAt,
		UpdatedAt: server.UpdatedAt,
	}, nil
}

func (s *ServersService) CreateServer(userId uint, dto dto.CreateServerDto, iv string) (ServerResponse, error) {
	var err error
	var encryptedSSHKey string
	if dto.SSHKey != nil {
		encryptedSSHKey, err = s.encryptionService.Encrypt(*dto.SSHKey, iv)
		if err != nil {
			return ServerResponse{}, err
		}
	}
	encryptedPassword, err := s.encryptionService.Encrypt(dto.Password, iv)
	if err != nil {
		return ServerResponse{}, err
	}
	server := Server{
		Name:     dto.Name,
		Username: dto.Username,
		Host:     dto.Host,
		Port:     dto.Port,
		UserID:   userId,
		SSHKey:   &encryptedSSHKey,
		Password: encryptedPassword,
	}
	if err := s.db.Create(&server).Error; err != nil {
		return ServerResponse{}, err
	}
	var decodedSSHKey string
	if server.SSHKey != nil {
		decodedSSHKey, err = s.encryptionService.Decrypt(*server.SSHKey, iv)
		if err != nil {
			return ServerResponse{}, err
		}
	} else {
		decodedSSHKey = ""
	}
	decodedPassword, err := s.encryptionService.Decrypt(server.Password, iv)
	if err != nil {
		return ServerResponse{}, err
	}
	return ServerResponse{
		ID:        server.ID,
		Name:      server.Name,
		Username:  server.Username,
		Host:      server.Host,
		Port:      server.Port,
		SSHKey:    decodedSSHKey,
		Password:  decodedPassword,
		CreatedAt: server.CreatedAt,
		UpdatedAt: server.UpdatedAt,
	}, nil
}

func (s *ServersService) UpdateServer(id, userId uint, updates map[string]interface{}, iv string) (ServerResponse, error) {
	var server Server
	if err := s.db.Where("id = ? AND user_id = ?", id, userId).First(&server).Error; err != nil {
		return ServerResponse{}, err
	}
	libs.SetStructFieldsFromMap(&server, updates)

	if updates["ssh_key"] != nil {
		encryptedSSHKey, err := s.encryptionService.Encrypt(updates["ssh_key"].(string), iv)
		if err != nil {
			return ServerResponse{}, err
		}
		server.SSHKey = &encryptedSSHKey
	}
	if updates["password"] != nil {
		encryptedPassword, err := s.encryptionService.Encrypt(updates["password"].(string), iv)
		if err != nil {
			return ServerResponse{}, err
		}
		server.Password = encryptedPassword
	}

	if err := s.db.Save(&server).Error; err != nil {
		return ServerResponse{}, err
	}

	var err error
	var decodedSSHKey string
	if server.SSHKey != nil {
		decodedSSHKey, err = s.encryptionService.Decrypt(*server.SSHKey, iv)
		if err != nil {
			return ServerResponse{}, err
		}
	} else {
		decodedSSHKey = ""
	}
	decodedPassword, err := s.encryptionService.Decrypt(server.Password, iv)
	if err != nil {
		return ServerResponse{}, err
	}
	return ServerResponse{
		ID:        server.ID,
		Name:      server.Name,
		Username:  server.Username,
		Host:      server.Host,
		Port:      server.Port,
		SSHKey:    decodedSSHKey,
		Password:  decodedPassword,
		CreatedAt: server.CreatedAt,
		UpdatedAt: server.UpdatedAt,
	}, nil
}

func (s *ServersService) DeleteServer(id, userId uint) error {
	if err := s.db.Where("id = ? AND user_id = ?", id, userId).Delete(&Server{}).Error; err != nil {
		return err
	}
	return nil
}
