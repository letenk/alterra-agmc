package models

import "time"

type Users struct {
	ID        string    `json:"id" validate:"required"`
	Fullname  string    `json:"fullname" validate:"required"`
	Email     string    `json:"email" validate:"required"`
	Password  string    `json:"password" validate:"required"`
	CreatedAt time.Time `json:"created_at" binding:"required"`
	UpdatedAt time.Time `json:"updated_at" binding:"required"`
}
