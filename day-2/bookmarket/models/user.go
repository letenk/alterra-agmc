package models

import "time"

type Users struct {
	ID        string    `json:"id" gorm:"type:uuid;primary_key" validate:"required"`
	Fullname  string    `json:"fullname" validate:"required"`
	Password  string    `json:"password" validate:"required"`
	CreatedAt time.Time `json:"created_at" binding:"required"`
	UpdatedAt time.Time `json:"updated_at" binding:"required"`
}
