package dto

type CreateOrUpdateBook struct {
	Name   string `json:"name" validate:"required"`
	Author string `json:"author" validate:"required"`
}
