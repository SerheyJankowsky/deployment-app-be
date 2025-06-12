package domains

import (
	"errors"
	"strings"
	"time"

	"deployer.com/libs"
	"deployer.com/modules/domains/dto"
	"gorm.io/gorm"
)

type SubDomainsService struct {
	db                *gorm.DB
	encryptionService *libs.EncryptionService
}

type SubDomainResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewSubDomainsService(db *gorm.DB) *SubDomainsService {
	return &SubDomainsService{db: db, encryptionService: libs.NewEncryptionService()}
}

func (s *SubDomainsService) GetSubDomains(userId uint, domainId uint, iv string) ([]SubDomainResponse, error) {
	var subDomains []SubDomain
	if err := s.db.Where("user_id = ? AND domain_id = ?", userId, domainId).Select("id, name, created_at, updated_at").Order("created_at DESC").Find(&subDomains).Error; err != nil {
		return nil, err
	}
	result := make([]SubDomainResponse, len(subDomains))
	for i, subDomain := range subDomains {
		result[i] = SubDomainResponse{
			ID:        subDomain.ID,
			Name:      subDomain.Name,
			CreatedAt: subDomain.CreatedAt,
			UpdatedAt: subDomain.UpdatedAt,
		}
	}
	return result, nil
}

func (s *SubDomainsService) GetSubDomain(id, userId uint, iv string) (SubDomainResponse, error) {
	var subDomain SubDomain
	if err := s.db.Where("id = ? AND user_id = ?", id, userId).First(&subDomain).Error; err != nil {
		return SubDomainResponse{}, err
	}

	return SubDomainResponse{
		ID:        subDomain.ID,
		Name:      subDomain.Name,
		CreatedAt: subDomain.CreatedAt,
		UpdatedAt: subDomain.UpdatedAt,
	}, nil
}

func (s *SubDomainsService) CreateSubDomain(userId uint, dto dto.CreateSubDomainDto, iv string) (SubDomainResponse, error) {
	domain := Domain{}
	if err := s.db.Where("id = ? AND user_id = ?", dto.DomainID, userId).First(&domain).Error; err != nil {
		return SubDomainResponse{}, err
	}
	if domain.UserID != userId {
		return SubDomainResponse{}, errors.New("domain not found")
	}
	if err := s.validateSubDomainName(domain.Name, dto.Name); err != nil {
		return SubDomainResponse{}, err
	}
	subDomain := SubDomain{
		Name:     dto.Name,
		DomainID: dto.DomainID,
		UserID:   userId,
	}
	if err := s.db.Create(&subDomain).Error; err != nil {
		return SubDomainResponse{}, err
	}
	return SubDomainResponse{
		ID:        subDomain.ID,
		Name:      subDomain.Name,
		CreatedAt: subDomain.CreatedAt,
		UpdatedAt: subDomain.UpdatedAt,
	}, nil
}

func (s *SubDomainsService) UpdateSubDomain(id, userId uint, updates map[string]interface{}, iv string) (SubDomainResponse, error) {
	var subDomain SubDomain
	if err := s.db.Where("id = ? AND user_id = ?", id, userId).First(&subDomain).Error; err != nil {
		return SubDomainResponse{}, err
	}
	libs.SetStructFieldsFromMap(&subDomain, updates)

	if err := s.db.Save(&subDomain).Error; err != nil {
		return SubDomainResponse{}, err
	}

	return SubDomainResponse{
		ID:        subDomain.ID,
		Name:      subDomain.Name,
		CreatedAt: subDomain.CreatedAt,
		UpdatedAt: subDomain.UpdatedAt,
	}, nil
}

func (s *SubDomainsService) DeleteSubDomain(id, userId uint) error {
	if err := s.db.Where("id = ? AND user_id = ?", id, userId).Delete(&SubDomain{}).Error; err != nil {
		return err
	}
	return nil
}

// validateSubDomainName checks if subDomainName is a valid subdomain of domainName
func (s *SubDomainsService) validateSubDomainName(domainName, subDomainName string) error {
	// Basic validation
	if domainName == "" || subDomainName == "" {
		return errors.New("domain name and subdomain name cannot be empty")
	}

	// Extract root domain (last two parts for most cases)
	rootDomain := extractRootDomain(domainName)
	if rootDomain == "" {
		return errors.New("invalid domain name")
	}

	// Check if subdomain ends with the root domain
	if !strings.HasSuffix(subDomainName, "."+rootDomain) {
		return errors.New("subdomain does not belong to the same domain")
	}

	// Ensure subdomain is actually longer (has additional parts)
	if subDomainName == rootDomain {
		return errors.New("subdomain cannot be the same as root domain")
	}

	return nil
}

// extractRootDomain extracts the root domain from a given domain
// For example: "api.example.com" -> "example.com"
func extractRootDomain(domain string) string {
	parts := strings.Split(domain, ".")
	if len(parts) < 2 {
		return ""
	}

	// For most cases, take the last two parts
	// This handles: example.com, api.example.com, sub.api.example.com
	return strings.Join(parts[len(parts)-2:], ".")
}

// Alternative approach: More flexible validation
func (s *SubDomainsService) validateSubDomainNameFlexible(domainName, subDomainName string) error {
	if domainName == "" || subDomainName == "" {
		return errors.New("domain name and subdomain name cannot be empty")
	}

	// Normalize domains (remove leading/trailing dots, convert to lowercase)
	domainName = strings.ToLower(strings.Trim(domainName, "."))
	subDomainName = strings.ToLower(strings.Trim(subDomainName, "."))

	// The subdomain should end with the domain name
	expectedSuffix := "." + domainName
	if !strings.HasSuffix(subDomainName, expectedSuffix) {
		return errors.New("subdomain does not belong to the specified domain")
	}

	// Ensure it's actually a subdomain (has additional parts before the domain)
	prefix := strings.TrimSuffix(subDomainName, expectedSuffix)
	if prefix == "" {
		return errors.New("provided subdomain is the same as the domain")
	}

	// Validate that the prefix doesn't contain invalid characters
	if strings.Contains(prefix, "..") || strings.HasPrefix(prefix, ".") || strings.HasSuffix(prefix, ".") {
		return errors.New("invalid subdomain format")
	}

	return nil
}
