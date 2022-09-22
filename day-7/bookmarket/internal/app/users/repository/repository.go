package repository

import (
	"bookmarket/internal/app/users/models"
	"context"

	"gorm.io/gorm"
)

type Repository interface {
	GetAll(ctx context.Context) ([]models.Users, error)
	FindByID(ctx context.Context, userID string) (models.Users, error)
	FindByEmail(ctx context.Context, email string) (models.Users, error)
	Create(ctx context.Context, user models.Users) (models.Users, error)
	Update(ctx context.Context, user models.Users) (models.Users, error)
	Delete(ctx context.Context, user models.Users) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) GetAll(ctx context.Context) ([]models.Users, error) {
	var users []models.Users

	err := r.db.WithContext(ctx).Model(&models.Users{}).Find(&users).Error
	if err != nil {
		return users, err
	}

	return users, nil
}

func (r *repository) FindByID(ctx context.Context, userID string) (models.Users, error) {
	var user models.Users

	if err := r.db.WithContext(ctx).Model(&models.Users{}).Where("id = ?", userID).First(&user).Error; err != nil {
		return user, err
	}

	return user, nil
}

func (r *repository) FindByEmail(ctx context.Context, email string) (models.Users, error) {
	var user models.Users

	if err := r.db.WithContext(ctx).Model(&models.Users{}).Where("email = ?", email).First(&user).Error; err != nil {
		return user, err
	}

	return user, nil
}

func (r *repository) Create(ctx context.Context, user models.Users) (models.Users, error) {
	err := r.db.WithContext(ctx).Model(&models.Users{}).Create(&user).Error
	if err != nil {
		return user, err
	}

	// IF success return user and error nil
	return user, nil
}

func (r *repository) Update(ctx context.Context, user models.Users) (models.Users, error) {
	if err := r.db.WithContext(ctx).Save(&user).Find(&user).Error; err != nil {
		return user, err
	}

	// IF success return user and error nil
	return user, nil
}

func (r *repository) Delete(ctx context.Context, user models.Users) error {
	if err := r.db.WithContext(ctx).Delete(&user).Error; err != nil {
		return err
	}

	return nil
}
