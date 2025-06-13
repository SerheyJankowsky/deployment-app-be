package auth

import (
	"errors"
	"time"

	"deployer.com/libs"
	"deployer.com/modules/auth/dto"
	"deployer.com/modules/users"
)

type AuthService struct {
	userService *users.UsersService
}

type LoginResponse struct {
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
	User         *users.User `json:"user"`
	Exp          int64       `json:"exp"`
}

type MeResponse struct {
	User *users.User `json:"user"`
}

type RefreshTokenResponse struct {
	AccessToken string `json:"access_token"`
	Exp         int64  `json:"exp"`
}

func NewAuthService(userService *users.UsersService) *AuthService {
	return &AuthService{userService: userService}
}

func (s *AuthService) Login(dto dto.LoginDto) (*LoginResponse, error) {
	user, err := s.userService.GetUserByEmail(dto.Email)
	if err != nil {
		return nil, err
	}
	if !libs.VerifyPassword(dto.Password, user.PasswordHash) {
		return nil, errors.New("invalid password")
	}
	userClaims := libs.UserClaims{
		UserID:   user.ID,
		Email:    user.Email,
		Role:     "",
		Verified: true,
		IV:       user.IV,
	}
	accessToken, err := libs.GenerateAccessToken(userClaims)
	if err != nil {
		return nil, err
	}
	refreshToken, err := libs.GenerateRefreshToken(userClaims)
	if err != nil {
		return nil, err
	}
	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         &user,
		Exp:          time.Now().Add(15 * time.Minute).Unix(),
	}, nil
}

func (s *AuthService) Register(dto dto.RegisterDto) (*LoginResponse, error) {
	isExist, _ := s.userService.IsUserExist(dto.User.Email, dto.User.FirstName+" "+dto.User.LastName)
	if isExist {
		return nil, errors.New("user already exists")
	}
	user, err := s.userService.CreateUser(&dto.User)
	if err != nil {
		return nil, err
	}
	userClaims := libs.UserClaims{
		UserID:   user.ID,
		Email:    user.Email,
		Role:     "",
		Verified: true,
		IV:       user.IV,
	}
	accessToken, err := libs.GenerateAccessToken(userClaims)
	if err != nil {
		return nil, err
	}
	refreshToken, err := libs.GenerateRefreshToken(userClaims)
	if err != nil {
		return nil, err
	}
	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         &user,
		Exp:          time.Now().Add(15 * time.Minute).Unix(),
	}, nil
}

func (s *AuthService) RefreshToken(rfToken string) (*RefreshTokenResponse, error) {
	claims, err := libs.ParseRefreshToken(rfToken)
	if err != nil {
		return nil, err
	}
	user, err := s.userService.GetUser(claims.UserID)
	if err != nil {
		return nil, err
	}
	userClaims := libs.UserClaims{
		UserID:   user.ID,
		Email:    user.Email,
		Role:     "",
		Verified: true,
		IV:       user.IV,
	}
	accessToken, err := libs.GenerateAccessToken(userClaims)
	if err != nil {
		return nil, err
	}
	return &RefreshTokenResponse{
		AccessToken: accessToken,
		Exp:         time.Now().Add(15 * time.Minute).Unix(),
	}, nil
}

func (s *AuthService) Me(userClaims *libs.UserClaims) (*MeResponse, error) {
	user, err := s.userService.GetUser(userClaims.UserID)
	if err != nil {
		return nil, err
	}
	return &MeResponse{
		User: &user,
	}, nil
}

func (s *AuthService) GenerateApiKey(userID uint) (*users.User, error) {
	user, err := s.userService.GenerateApiKey(userID)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *AuthService) RevokeApiKey(userID uint) error {
	return s.userService.RevokeApiKey(userID)
}

func (s *AuthService) GetUserApiKey(userID uint) (*users.User, error) {
	user, err := s.userService.GetUser(userID)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
