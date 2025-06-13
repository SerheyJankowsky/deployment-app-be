package users

import (
	"fmt"

	"deployer.com/libs"
	"deployer.com/modules/users/dto"
	"gorm.io/gorm"
)

type UsersService struct {
	db                *gorm.DB
	encryptionService *libs.EncryptionService
}

func NewUsersService(db *gorm.DB) *UsersService {
	return &UsersService{db: db, encryptionService: libs.NewEncryptionService()}
}

func (s *UsersService) GetUser(id uint) (User, error) {
	var user User
	if err := s.db.First(&user, id).Error; err != nil {
		return User{}, err
	}
	if user.ApiKey != "" {
		decryptedApiKey, err := s.encryptionService.Decrypt(user.ApiKey, user.IV)
		if err != nil {
			return User{}, err
		}
		user.ApiKey = decryptedApiKey
	}
	return user, nil
}

func (s *UsersService) GetUserByEmail(email string) (User, error) {
	var user User
	if err := s.db.Where("email = ?", email).First(&user).Error; err != nil {
		return User{}, err
	}
	return user, nil
}

func (s *UsersService) IsUserExist(email string, username string) (bool, error) {
	var user User
	if err := s.db.Where("email = ? OR username = ?", email, username).First(&user).Error; err != nil {
		return false, err
	}
	return true, nil
}

func (s *UsersService) CreateUser(user *dto.CreateUserDto) (User, error) {
	hashedPassword, err := libs.HashPassword(user.Password)
	if err != nil {
		return User{}, err
	}
	iv := s.encryptionService.GenIv()

	userEntity := User{
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Username:     fmt.Sprintf("%s%s", user.FirstName, user.LastName),
		Email:        user.Email,
		Phone:        user.Phone,
		Country:      user.Country,
		PasswordHash: hashedPassword,
		IV:           iv,
	}

	if err := s.db.Create(&userEntity).Error; err != nil {
		return User{}, err
	}

	return userEntity, nil
}

func (s *UsersService) UpdateUser(user *dto.UpdateUserDto) (User, error) {
	var userEntity User
	if err := s.db.Model(&User{}).Where("id = ?", user.ID).Updates(user).Error; err != nil {
		return User{}, err
	}
	if err := s.db.First(&userEntity, user.ID).Error; err != nil {
		return User{}, err
	}
	return userEntity, nil
}

func (s *UsersService) UpdateUserApiKey(userId uint, iv string) (User, error) {
	apiKey := libs.GenerateApiKey()
	encryptedApiKey, err := s.encryptionService.Encrypt(apiKey, iv)
	if err != nil {
		return User{}, err
	}
	if err := s.db.Model(&User{}).Where("id = ?", userId).Update("api_key", encryptedApiKey).Error; err != nil {
		return User{}, err
	}
	return s.GetUser(userId)
}

func (s *UsersService) DeleteUser(id uint) error {
	if err := s.db.Delete(&User{}, id).Error; err != nil {
		return err
	}
	return nil
}
