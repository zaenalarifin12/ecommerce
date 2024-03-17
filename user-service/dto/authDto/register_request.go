package authDto

import "strings"

// RegisterRequest represents a registration request.
type RegisterRequest struct {
	FullName string `json:"full_name" validate:"required" example:"Babahaha"`
	Phone    string `json:"phone" validate:"required,number" example:"0897654545"`
	Email    string `json:"email" validate:"required,email" example:"john@example.com"`
	Username string `json:"username" validate:"required,min=5" example:"johndoe"`
	Password string `json:"password" validate:"required,min=8" example:"password123"`
}

// Sanitize the fields of RegisterRequest struct
func (req *RegisterRequest) Sanitize() {
	// Trim leading and trailing spaces from each field
	req.FullName = strings.TrimSpace(req.FullName)
	req.Phone = strings.TrimSpace(req.Phone)
	req.Email = strings.TrimSpace(req.Email)
	req.Username = strings.TrimSpace(req.Username)
	req.Password = strings.TrimSpace(req.Password)

	// Optionally, convert email to lowercase
	req.Email = strings.ToLower(req.Email)
}
