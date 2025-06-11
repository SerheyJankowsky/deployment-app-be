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

func (s *SubDomainsService) validateSubDomainName(domainName, subDomainName string) error {
	splittedDomainName := strings.Split(domainName, ".")
	if len(splittedDomainName) < 2 {
		return errors.New("invalid domain name")
	}
	subDomainName = strings.TrimPrefix(subDomainName, splittedDomainName[0]+".")
	if subDomainName == "" {
		return errors.New("invalid sub domain name")
	}
	return nil
}
