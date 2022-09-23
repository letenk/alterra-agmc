package tests

import (
	"bookmarket/database"
	"bookmarket/internal/app/users/models"
	"bookmarket/internal/app/users/repository"
	"bookmarket/internal/app/users/service"
	"bookmarket/internal/dto"
	"context"
	"log"
	"testing"

	"github.com/rizkydarmawan-letenk/jabufaker"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func CreateRandomUserService(t *testing.T) models.Users {
	// Open connection to db
	db := database.SetupTestDB()

	// Use repository
	userRepository := repository.NewRepository(db)
	// Use service
	userService := service.NewService(userRepository)

	// Create payload sample data
	payload := dto.CreateRequestBody{
		Fullname: jabufaker.RandomPerson(),
		Email:    jabufaker.RandomEmail(),
		Password: "password",
	}

	// Create user
	newUser, err := userService.Create(context.Background(), &payload)
	if err != nil {
		log.Panic(err)
	}

	// Test pass
	assert.Equal(t, payload.Fullname, newUser.Fullname)
	assert.Equal(t, payload.Email, newUser.Email)
	assert.NotEmpty(t, newUser.ID)
	assert.NotEmpty(t, newUser.CreatedAt)
	assert.NotEmpty(t, newUser.UpdatedAt)

	// Compare password
	err = bcrypt.CompareHashAndPassword([]byte(newUser.Password), []byte("password"))
	assert.Nil(t, err)

	return newUser
}

func TestCreateUserServiceSuccess(t *testing.T) {
	CreateRandomUserService(t)
}

func TestGetAllService(t *testing.T) {
	// Create some users
	for i := 0; i < 5; i++ {
		CreateRandomUserService(t)
	}
	// Open connection to db
	db := database.SetupTestDB()

	// Use repository
	userRepository := repository.NewRepository(db)
	// Use service
	userService := service.NewService(userRepository)

	users, err := userService.GetAll(context.Background())
	if err != nil {
		log.Panic(err)
	}

	for _, user := range users {
		assert.NotEmpty(t, user.ID)
		assert.NotEmpty(t, user.Fullname)
		assert.NotEmpty(t, user.Email)
		assert.NotEmpty(t, user.Password)
		assert.NotEmpty(t, user.CreatedAt)
		assert.NotEmpty(t, user.UpdatedAt)
	}
}

func TestGetUserServiceById(t *testing.T) {
	// Create new users
	newUser := CreateRandomUserService(t)

	// Open connection to db
	db := database.SetupTestDB()

	// Use repository
	userRepository := repository.NewRepository(db)
	// Use service
	userService := service.NewService(userRepository)

	// Find by id
	user, err := userService.FindByID(context.Background(), newUser.ID)
	if err != nil {
		log.Panic(err)
	}

	assert.Equal(t, newUser.ID, user.ID)
	assert.Equal(t, newUser.Fullname, user.Fullname)
	assert.Equal(t, newUser.Email, user.Email)
	assert.NotEmpty(t, newUser.CreatedAt)
	assert.NotEmpty(t, newUser.UpdatedAt)

	// Compare password
	err = bcrypt.CompareHashAndPassword([]byte(newUser.Password), []byte("password"))
	assert.Nil(t, err)
}

func TestFindUserServiceByEmail(t *testing.T) {
	// Create new users
	newUser := CreateRandomUserService(t)

	// Open connection to db
	db := database.SetupTestDB()

	// Use repository
	userRepository := repository.NewRepository(db)
	// Use service
	userService := service.NewService(userRepository)

	// Find by id
	user, err := userService.FindByEmail(context.Background(), newUser.Email)
	if err != nil {
		log.Panic(err)
	}

	assert.Equal(t, newUser.ID, user.ID)
	assert.Equal(t, newUser.Fullname, user.Fullname)
	assert.Equal(t, newUser.Email, user.Email)
	assert.NotEmpty(t, newUser.CreatedAt)
	assert.NotEmpty(t, newUser.UpdatedAt)

	// Compare password
	err = bcrypt.CompareHashAndPassword([]byte(newUser.Password), []byte("password"))
	assert.Nil(t, err)
}

func TestUpdateUserService(t *testing.T) {
	// Create new users
	newUser := CreateRandomUserService(t)

	// Open connection to db
	db := database.SetupTestDB()

	// Use repository
	userRepository := repository.NewRepository(db)
	// Use service
	userService := service.NewService(userRepository)

	// Create data payload
	userUpdate := dto.UpdateRequestBody{
		Fullname: jabufaker.RandomPerson(),
		Password: "update",
	}

	// Update
	user, err := userService.Update(context.Background(), newUser.ID, &userUpdate)
	if err != nil {
		log.Panic(err)
	}

	assert.Equal(t, newUser.ID, user.ID)
	assert.Equal(t, userUpdate.Fullname, user.Fullname)
	assert.Equal(t, newUser.Email, user.Email)
	assert.NotEmpty(t, newUser.CreatedAt)
	assert.NotEmpty(t, newUser.UpdatedAt)

	// Compare password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userUpdate.Password))
	assert.Nil(t, err)
}

func TestDeleteUserService(t *testing.T) {
	// Create new users
	newUser := CreateRandomUserService(t)

	// Open connection to db
	db := database.SetupTestDB()

	// Use repository
	userRepository := repository.NewRepository(db)
	// Use service
	userService := service.NewService(userRepository)

	// Delete user
	err := userService.Delete(context.Background(), newUser.ID)
	assert.Nil(t, err)
}
