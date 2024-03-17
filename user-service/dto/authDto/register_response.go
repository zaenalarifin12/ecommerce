package authDto

import "github.com/google/uuid"

type RegisterResponse struct {
	Uuid     uuid.UUID `json:"uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	FullName string    `json:"full_name" example:"John Doe"`
	Phone    string    `json:"phone" example:"1234567890"`
	Email    string    `json:"email" example:"john@gmail.com"`
	Username string    `json:"username" example:"johndoe"`
}

// DataRegisterResponse represents a response to a registration request.
type DataRegisterResponse struct {
	Data RegisterResponse `json:"data"`
}
