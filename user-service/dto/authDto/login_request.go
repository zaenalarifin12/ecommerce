package authDto

import "strings"

// LoginRequest represents a request to log in.
// swagger:parameters LoginRequest
type LoginRequest struct {
	Email    string `json:"email" validate:"email,required" example:"john@example.com"`
	Password string `json:"password" validate:"required" example:"password123"`
}

// Sanitize the fields of LoginRequest struct
func (req *LoginRequest) Sanitize() {
	// Trim leading and trailing spaces from each field
	req.Email = strings.TrimSpace(req.Email)
	req.Password = strings.TrimSpace(req.Password)

	// Optionally, convert email to lowercase
	req.Email = strings.ToLower(req.Email)
}
