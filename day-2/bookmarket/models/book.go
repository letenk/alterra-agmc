package models

import "time"

type Books struct {
	ID        string    `json:"id" validate:"required"`
	Name      string    `json:"name" binding:"required"`
	Author    string    `json:"author" binding:"required"`
	CreatedAt time.Time `json:"created_at" binding:"required"`
	UpdatedAt time.Time `json:"updated_at" binding:"required"`
}
