package authDto

import "github.com/google/uuid"

type LoginResponse struct {
	Token    string    `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxMjM0NTY3ODkwLCJlbWFpbCI6ImpvaG5AZXhhbXBsZS5jb20iLCJpYXQiOjE1MTYyMzkwMjIsImV4cCI6MTUxNjI0NTQyMn0.QxegqwOfs8A8U6bgthmZG9y2AmWJ2t9UcYY9rKiu2I8"`
	FullName string    `json:"full_name" example:"John Doe"`
	Uuid     uuid.UUID `json:"uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	Phone    string    `json:"phone" example:"1234567890"`
	Email    string    `json:"email" example:"john@example.com"`
	Username string    `json:"username" example:"johndoe"`
}

// DataLoginResponse represents the response for a successful login.
// swagger:response DataLoginResponse
type DataLoginResponse struct {
	Data LoginResponse `json:"data"`
}
