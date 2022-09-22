package tests

import (
	"bookmarket/database"
	"bookmarket/internal/app/users/models"
	"bookmarket/internal/app/users/repository"
	"context"
	"log"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/rizkydarmawan-letenk/jabufaker"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func CreateRepositoryCreateRandomUser(t *testing.T) models.Users {
	// Open connection
	db := database.SetupTestDB()

	// Use repository
	userRepository := repository.NewRepository(db)

	// Generate uuid
	id := uuid.NewString()

	// Hash password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	if err != nil {
		log.Panic(err)
	}

	// Create sample data
	newUser := models.Users{
		ID:        id,
		Fullname:  jabufaker.RandomPerson(),
		Email:     jabufaker.RandomEmail(),
		Password:  string(passwordHash),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Create user
	user, err := userRepository.Create(context.Background(), newUser)
	if err != nil {
		log.Panic(err)
	}

	// Test pass
	assert.Equal(t, newUser.ID, user.ID)
	assert.Equal(t, newUser.Fullname, user.Fullname)
	assert.Equal(t, newUser.Email, user.Email)
	assert.Equal(t, newUser.CreatedAt, user.CreatedAt)
	assert.Equal(t, newUser.UpdatedAt, user.UpdatedAt)

	// Compare password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte("password"))
	assert.Nil(t, err)

	return user
}
func TestCreateUserRepositorySuccess(t *testing.T) {
	CreateRepositoryCreateRandomUser(t)
}

func TestGetRepositoryAll(t *testing.T) {
	// Create some user
	for i := 0; i < 5; i++ {
		CreateRepositoryCreateRandomUser(t)
	}

	// Open connection
	db := database.SetupTestDB()

	// Use repository
	userRepository := repository.NewRepository(db)

	// Get all users
	users, err := userRepository.GetAll(context.Background())
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

func TestFindByIDRepositorySuccess(t *testing.T) {
	// Create new user
	newUser := CreateRepositoryCreateRandomUser(t)

	// Open connection
	db := database.SetupTestDB()

	// Use repository
	userRepository := repository.NewRepository(db)

	// Find user by id
	user, err := userRepository.FindByID(context.Background(), newUser.ID)
	if err != nil {
		log.Panic()
	}

	// Test pass
	assert.Equal(t, newUser.ID, user.ID)
	assert.Equal(t, newUser.Fullname, user.Fullname)
	assert.Equal(t, newUser.Email, user.Email)
	assert.NotEmpty(t, user.CreatedAt)
	assert.NotEmpty(t, user.UpdatedAt)

	// Compare password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte("password"))
	assert.Nil(t, err)
}

func TestFindByEmailRepositorySuccess(t *testing.T) {
	// Create new user
	newUser := CreateRepositoryCreateRandomUser(t)

	// Open connection
	db := database.SetupTestDB()

	// Use repository
	userRepository := repository.NewRepository(db)

	// Find user by id
	user, err := userRepository.FindByEmail(context.Background(), newUser.Email)
	if err != nil {
		log.Panic()
	}

	// Test pass
	assert.Equal(t, newUser.ID, user.ID)
	assert.Equal(t, newUser.Fullname, user.Fullname)
	assert.Equal(t, newUser.Email, user.Email)
	assert.NotEmpty(t, user.CreatedAt)
	assert.NotEmpty(t, user.UpdatedAt)

	// Compare password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte("password"))
	assert.Nil(t, err)
}

func TestUpdateRepositorySuccess(t *testing.T) {
	// Create new user
	newUser := CreateRepositoryCreateRandomUser(t)

	// Open connection
	db := database.SetupTestDB()

	// Use repository
	userRepository := repository.NewRepository(db)

	// Hash password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("update"), bcrypt.DefaultCost)
	if err != nil {
		log.Panic(err)
	}

	// Create sample data
	dataUpdate := models.Users{
		ID:        newUser.ID, // use id from newUser
		Fullname:  jabufaker.RandomPerson(),
		Email:     jabufaker.RandomEmail(),
		Password:  string(passwordHash),
		CreatedAt: newUser.CreatedAt, // use from newUser
		UpdatedAt: time.Now(),
	}

	// Update
	user, err := userRepository.Update(context.Background(), dataUpdate)
	if err != nil {
		log.Panic(err)
	}

	// Test pass
	assert.Equal(t, dataUpdate.ID, user.ID)
	assert.Equal(t, dataUpdate.Fullname, user.Fullname)
	assert.Equal(t, dataUpdate.Email, user.Email)
	assert.NotEmpty(t, user.CreatedAt)
	assert.NotEmpty(t, user.UpdatedAt)

	// Compare password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte("update"))
	assert.Nil(t, err)
}

func TestDeleteRepositorySuccess(t *testing.T) {
	// Create new user
	newUser := CreateRepositoryCreateRandomUser(t)

	// Open connection
	db := database.SetupTestDB()

	// Use repository
	userRepository := repository.NewRepository(db)

	// Update
	err := userRepository.Delete(context.Background(), newUser)
	if err != nil {
		log.Panic(err)
	}

	// Test pass
	assert.NoError(t, err)
}
