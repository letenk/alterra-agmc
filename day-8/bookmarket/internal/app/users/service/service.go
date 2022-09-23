package service

import (
	"bookmarket/internal/app/users/models"
	"bookmarket/internal/app/users/repository"
	"bookmarket/internal/dto"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	GetAll(ctx context.Context) ([]models.Users, error)
	FindByID(ctx context.Context, userID string) (models.Users, error)
	FindByEmail(ctx context.Context, email string) (models.Users, error)
	Create(ctx context.Context, payload *dto.CreateRequestBody) (models.Users, error)
	Update(ctx context.Context, ID string, payload *dto.UpdateRequestBody) (models.Users, error)
	Delete(ctx context.Context, ID string) error
}

type service struct {
	repository repository.Repository
}

func NewService(repository repository.Repository) *service {
	return &service{repository}
}

func (s *service) GetAll(ctx context.Context) ([]models.Users, error) {
	users, err := s.repository.GetAll(ctx)
	if err != nil {
		return users, err
	}

	return users, nil
}

func (s *service) FindByID(ctx context.Context, userID string) (models.Users, error) {
	user, err := s.repository.FindByID(ctx, userID)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (s *service) FindByEmail(ctx context.Context, email string) (models.Users, error) {
	user, err := s.repository.FindByEmail(ctx, email)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (s *service) Create(ctx context.Context, payload *dto.CreateRequestBody) (models.Users, error) {
	// Passing data to object user
	user := models.Users{}
	user.Fullname = payload.Fullname
	user.Email = payload.Email

	// Generate uuid
	id := uuid.NewString()
	user.ID = id

	// Hash password input with package bcrypy
	passworHash, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.MinCost)
	// Check if bcrypt error
	if err != nil {
		return user, err
	}
	// Insert password hash to struct PasswordHash
	user.Password = string(passworHash)

	// If no, create new user
	newUser, err := s.repository.Create(ctx, user)
	if err != nil {
		return user, err
	}

	// If no error return new user
	return newUser, nil
}

func (s *service) Update(ctx context.Context, ID string, payload *dto.UpdateRequestBody) (models.Users, error) {
	// Find user by id
	user, err := s.repository.FindByID(ctx, ID)
	if err != nil {
		message := fmt.Sprintf("user with id: %s not found", ID)
		return user, errors.New(message)
	}

	// If user is find, passing data payload into object user
	user.Fullname = payload.Fullname

	// Hash password input with package bcrypy
	passworHash, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.MinCost)
	// Check if bcrypt error
	if err != nil {
		return user, err
	}
	// Insert password hash to struct PasswordHash
	user.Password = string(passworHash)
	// Update user
	userUpdated, err := s.repository.Update(ctx, user)
	if err != nil {
		return userUpdated, err
	}

	return userUpdated, nil
}

func (s *service) Delete(ctx context.Context, ID string) error {
	// Find user by id
	user, err := s.repository.FindByID(ctx, ID)
	if err != nil {
		message := fmt.Sprintf("user with id: %s not found", ID)
		return errors.New(message)
	}

	err = s.repository.Delete(context.Background(), user)
	if err != nil {
		return err
	}

	return nil
}
