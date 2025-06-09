package tests

import (
	"errors"
	"os"
	"testing"

	"deployer.com/modules/users"
	"deployer.com/modules/users/dto"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestDB() (*gorm.DB, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return nil, errors.New("DATABASE_URL is not set")
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func TestCreateGetDeleteUser(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}
	svc := users.NewUsersService(db)
	createDto := &dto.CreateUserDto{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
		Password:  "password123",
		Phone:     "+1234567890",
		Country:   "US",
	}
	user, err := svc.CreateUser(createDto)
	assert.NoError(t, err)
	assert.Equal(t, createDto.Email, user.Email)
	assert.NotEmpty(t, user.PasswordHash)

	fetched, err := svc.GetUser(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, user.Email, fetched.Email)

	err = svc.DeleteUser(user.ID)
	assert.NoError(t, err)
	_, err = svc.GetUser(user.ID)
	assert.Error(t, err)
}

func TestUpdateUser(t *testing.T) {
	db, err := setupTestDB()
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}
	svc := users.NewUsersService(db)
	createDto := &dto.CreateUserDto{
		FirstName: "Jane",
		LastName:  "Smith",
		Email:     "jane@example.com",
		Password:  "password123",
		Phone:     "+1234567890",
		Country:   "US",
	}
	user, _ := svc.CreateUser(createDto)
	updateDto := &dto.UpdateUserDto{
		ID:        user.ID,
		FirstName: "Janet",
		LastName:  "Smith",
		// Email:     "janet@example.com",
		Phone:   user.Phone,
		Country: user.Country,
	}
	updated, err := svc.UpdateUser(updateDto)
	assert.NoError(t, err)
	assert.Equal(t, updateDto.FirstName, updated.FirstName)
	assert.Equal(t, updateDto.LastName, updated.LastName)
	assert.Equal(t, updateDto.Phone, updated.Phone)
	assert.Equal(t, updateDto.Country, updated.Country)
	err = svc.DeleteUser(user.ID)
	assert.NoError(t, err)
	_, err = svc.GetUser(user.ID)
	assert.Error(t, err)
}
