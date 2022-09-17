package lib

import (
	"bookmarket/models"
	"time"
)

type CreateOrUpdateBook struct {
	Name   string `json:"name" validate:"required"`
	Author string `json:"author" validate:"required"`
}

type CreateUser struct {
	Fullname string `json:"fullname" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UpdateUser struct {
	Fullname string `json:"fullname" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UserFormatter struct {
	ID        string    `json:"id"`
	Fullname  string    `json:"fullname"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func FormatUser(user models.Users) UserFormatter {
	if user.ID == "" {
		return UserFormatter{}
	}

	formatter := UserFormatter{
		ID:        user.ID,
		Fullname:  user.Fullname,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return formatter
}

func FormatUsers(users []models.Users) []UserFormatter {
	// If data not available, retun empty array
	if len(users) == 0 {
		return []UserFormatter{}
	}

	// If no, iteration users and append into var userFormatter
	var userFormatter []UserFormatter
	for _, user := range users {
		formatter := FormatUser(user)
		userFormatter = append(userFormatter, formatter)
	}

	return userFormatter
}
