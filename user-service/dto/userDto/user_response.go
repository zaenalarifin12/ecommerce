package userDto

import (
	"github.com/google/uuid"
	"time"
)

type UserResponse struct {
	Uuid          uuid.UUID  `json:"uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	FullName      string     `json:"full_name" example:"John Doe"`
	Phone         string     `json:"phone" example:"1234567890"`
	Email         string     `json:"email" example:"john@example.com"`
	Username      string     `json:"username" example:"johndoe"`
	EmailVerifyAt *time.Time `json:"email_verify_at"`
	CreatedAt     time.Time  `json:"created_at" example:"2024-02-29T12:00:00Z"`
	UpdatedAt     *time.Time `json:"updated_at" `
}

// DataUserResponse represents a response to a registration request.
// swagger:response DataUserResponse
type DataUserResponse struct {
	Data UserResponse `json:"data"`
}
