package domains

import (
	"fmt"
	"time"

	"deployer.com/libs"
	"deployer.com/modules/domains/dto"
	"gorm.io/gorm"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

type DomainsService struct {
	db                *gorm.DB
	encryptionService *libs.EncryptionService
}

type DomainResponse struct {
	ID         uint        `json:"id"`
	Name       string      `json:"name"`
	SSLCert    string      `json:"ssl_cert"`
	SSLKey     string      `json:"ssl_key"`
	SubDomains []SubDomain `json:"sub_domains"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
}

func NewDomainsService(db *gorm.DB) *DomainsService {
	return &DomainsService{db: db, encryptionService: libs.NewEncryptionService()}
}

func (s *DomainsService) GetDomains(userId uint, iv string) ([]DomainResponse, error) {
	var domains []Domain
	if err := s.db.Where("user_id = ?", userId).Preload("SubDomains").Select("id, name, ssl_cert, ssl_key, created_at, updated_at").Order("created_at DESC").Find(&domains).Error; err != nil {
		return nil, err
	}
	result := make([]DomainResponse, len(domains))
	for i, domain := range domains {
		var decodedCert, decodedKey string

		// Decrypt SSL certificate only if not empty
		if domain.SSLCert != "" {
			fmt.Printf("Original SSL Cert (first 50 chars): %s\n", domain.SSLCert[:min(50, len(domain.SSLCert))])
			fmt.Printf("IV: %s\n", iv)
			decoded, err := s.encryptionService.Decrypt(domain.SSLCert, iv)
			if err != nil {
				// If decryption fails, return as-is (might be already decrypted)
				fmt.Println("Decryption failed for SSL cert", err)
				decodedCert = domain.SSLCert
			} else {
				fmt.Printf("Decrypted SSL Cert (first 50 chars): %s\n", decoded[:min(50, len(decoded))])
				decodedCert = decoded
			}
		}

		// Decrypt SSL key only if not empty
		if domain.SSLKey != "" {
			fmt.Printf("Original SSL Key (first 50 chars): %s\n", domain.SSLKey[:min(50, len(domain.SSLKey))])
			decoded, err := s.encryptionService.Decrypt(domain.SSLKey, iv)
			if err != nil {
				// If decryption fails, return as-is (might be already decrypted)
				fmt.Println("Decryption failed for SSL key", err)
				decodedKey = domain.SSLKey
			} else {
				fmt.Printf("Decrypted SSL Key (first 50 chars): %s\n", decoded[:min(50, len(decoded))])
				decodedKey = decoded
			}
		}

		result[i] = DomainResponse{
			ID:         domain.ID,
			Name:       domain.Name,
			SSLCert:    decodedCert,
			SSLKey:     decodedKey,
			SubDomains: domain.SubDomains,
			CreatedAt:  domain.CreatedAt,
			UpdatedAt:  domain.UpdatedAt,
		}
	}
	return result, nil
}

func (s *DomainsService) GetDomain(id, userId uint, iv string) (DomainResponse, error) {
	var domain Domain
	if err := s.db.Where("id = ? AND user_id = ?", id, userId).Preload("SubDomains").First(&domain).Error; err != nil {
		return DomainResponse{}, err
	}

	var decodedCert, decodedKey string

	// Decrypt SSL certificate only if not empty
	if domain.SSLCert != "" {
		fmt.Printf("GetDomain - Original SSL Cert (first 50 chars): %s\n", domain.SSLCert[:min(50, len(domain.SSLCert))])
		fmt.Printf("GetDomain - IV: %s\n", iv)
		decoded, err := s.encryptionService.Decrypt(domain.SSLCert, iv)
		if err != nil {
			// If decryption fails, return as-is (might be already decrypted)
			fmt.Println("GetDomain - Decryption failed for SSL cert", err)
			decodedCert = domain.SSLCert
		} else {
			fmt.Printf("GetDomain - Decrypted SSL Cert (first 50 chars): %s\n", decoded[:min(50, len(decoded))])
			decodedCert = decoded
		}
	}

	// Decrypt SSL key only if not empty
	if domain.SSLKey != "" {
		fmt.Printf("GetDomain - Original SSL Key (first 50 chars): %s\n", domain.SSLKey[:min(50, len(domain.SSLKey))])
		decoded, err := s.encryptionService.Decrypt(domain.SSLKey, iv)
		if err != nil {
			// If decryption fails, return as-is (might be already decrypted)
			fmt.Println("GetDomain - Decryption failed for SSL key", err)
			decodedKey = domain.SSLKey
		} else {
			fmt.Printf("GetDomain - Decrypted SSL Key (first 50 chars): %s\n", decoded[:min(50, len(decoded))])
			decodedKey = decoded
		}
	}

	return DomainResponse{
		ID:         domain.ID,
		Name:       domain.Name,
		SSLCert:    decodedCert,
		SSLKey:     decodedKey,
		SubDomains: domain.SubDomains,
		CreatedAt:  domain.CreatedAt,
		UpdatedAt:  domain.UpdatedAt,
	}, nil
}

func (s *DomainsService) CreateDomain(userId uint, dto dto.CreateDomainDto, iv string) (DomainResponse, error) {
	encryptedCert, err := s.encryptionService.Encrypt(dto.SSLCert, iv)
	if err != nil {
		return DomainResponse{}, err
	}
	encryptedKey, err := s.encryptionService.Encrypt(dto.SSLKey, iv)
	if err != nil {
		return DomainResponse{}, err
	}
	domain := Domain{
		Name:    dto.Name,
		SSLCert: encryptedCert,
		SSLKey:  encryptedKey,
		UserID:  userId,
	}
	if err := s.db.Create(&domain).Error; err != nil {
		return DomainResponse{}, err
	}
	return DomainResponse{
		ID:        domain.ID,
		Name:      domain.Name,
		SSLCert:   dto.SSLCert,
		SSLKey:    dto.SSLKey,
		CreatedAt: domain.CreatedAt,
		UpdatedAt: domain.UpdatedAt,
	}, nil
}

func (s *DomainsService) UpdateDomain(id, userId uint, updates map[string]interface{}, iv string) (DomainResponse, error) {
	var domain Domain
	if err := s.db.Where("id = ? AND user_id = ?", id, userId).Preload("SubDomains").First(&domain).Error; err != nil {
		return DomainResponse{}, err
	}
	libs.SetStructFieldsFromMap(&domain, updates)

	// Encrypt SSL certificate if updated
	if updates["ssl_cert"] != nil {
		encrypted, err := s.encryptionService.Encrypt(domain.SSLCert, iv)
		if err != nil {
			return DomainResponse{}, err
		}
		domain.SSLCert = encrypted
	}

	// Encrypt SSL key if updated
	if updates["ssl_key"] != nil {
		encrypted, err := s.encryptionService.Encrypt(domain.SSLKey, iv)
		if err != nil {
			return DomainResponse{}, err
		}
		domain.SSLKey = encrypted
	}

	if err := s.db.Save(&domain).Error; err != nil {
		return DomainResponse{}, err
	}
	decodedCert, err := s.encryptionService.Decrypt(domain.SSLCert, iv)
	if err != nil {
		return DomainResponse{}, err
	}
	decodedKey, err := s.encryptionService.Decrypt(domain.SSLKey, iv)
	if err != nil {
		return DomainResponse{}, err
	}
	return DomainResponse{
		ID:         domain.ID,
		Name:       domain.Name,
		SSLCert:    decodedCert,
		SSLKey:     decodedKey,
		SubDomains: domain.SubDomains,
		CreatedAt:  domain.CreatedAt,
		UpdatedAt:  domain.UpdatedAt,
	}, nil
}

func (s *DomainsService) DeleteDomain(id, userId uint) error {
	if err := s.db.Where("id = ? AND user_id = ?", id, userId).Delete(&Domain{}).Error; err != nil {
		return err
	}
	return nil
}
