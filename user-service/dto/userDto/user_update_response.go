package userDto

import (
	"github.com/google/uuid"
	"time"
)

type UserUpdateResponse struct {
	Uuid          uuid.UUID  `json:"uuid" example:"f47ac10b-58cc-4372-a567-0e02b2c3d479"`
	FullName      string     `json:"full_name" example:"John Doe"`
	Phone         string     `json:"phone" example:"+1234567890"`
	Email         string     `json:"email" example:"john@example.com"`
	Username      string     `json:"username" example:"johndoe123"`
	EmailVerifyAt *time.Time `json:"email_verify_at"`
	CreatedAt     time.Time  `json:"created_at" `
	UpdatedAt     *time.Time `json:"updated_at" `
}

// DataUserUpdateResponse represents a response to a registration request.
// swagger:response DataUserUpdateResponse
type DataUserUpdateResponse struct {
	Message string             `json:"message"`
	Data    UserUpdateResponse `json:"data"`
}
