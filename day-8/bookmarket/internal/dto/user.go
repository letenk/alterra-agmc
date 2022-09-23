package dto

import (
	"bookmarket/internal/app/users/models"
	"time"
)

type (
	LoginPayload struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}

	CreateRequestBody struct {
		Fullname string `json:"fullname" validate:"required"`
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}

	UpdateRequestBody struct {
		Fullname string `json:"fullname" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	UserResponse struct {
		ID        string    `json:"id"`
		Fullname  string    `json:"fullname"`
		Email     string    `json:"email"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}
)

func FormatUserResponse(user models.Users) UserResponse {
	if user.ID == "" {
		return UserResponse{}
	}

	formatter := UserResponse{
		ID:        user.ID,
		Fullname:  user.Fullname,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return formatter
}

func FormatUsersResponse(users []models.Users) []UserResponse {
	// If data not available, retun empty array
	if len(users) == 0 {
		return []UserResponse{}
	}

	// If no, iteration users and append into var UserResponse
	var userFormatter []UserResponse
	for _, user := range users {
		formatter := FormatUserResponse(user)
		userFormatter = append(userFormatter, formatter)
	}

	return userFormatter
}
