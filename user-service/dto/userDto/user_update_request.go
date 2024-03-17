package userDto

import (
	"github.com/google/uuid"
	"strings"
)

// UserUpdateRequest represents a request to update user information.
// swagger:parameters UserUpdateRequest
type UserUpdateRequest struct {
	Uuid     uuid.UUID `json:"uuid" example:"f47ac10b-58cc-4372-a567-0e02b2c3d479"`
	FullName string    `json:"full_name" validate:"required" example:"John Doe"`
	Phone    string    `json:"phone" validate:"required,number" example:"1234567890"`
	Username string    `json:"username" validate:"required,min=5" example:"johndoe123"`
	Password string    `json:"password" validate:"required,min=8" example:"password123"`
}

// Sanitize the fields of RegisterRequest struct
func (req *UserUpdateRequest) Sanitize() {
	// Trim leading and trailing spaces from each field
	req.FullName = strings.TrimSpace(req.FullName)
	req.Phone = strings.TrimSpace(req.Phone)
	req.Username = strings.TrimSpace(req.Username)
	req.Password = strings.TrimSpace(req.Password)
}
