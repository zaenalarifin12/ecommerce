package dto

import (
	"github.com/google/uuid"
	"time"
)

// CartResponse represents the response body for retrieving cart information.
// swagger:parameters CartResponse
type CartResponse struct {
	Uuid        uuid.UUID  `json:"uuid" example:"550e8400-e29b-41d4-a716-446655440000"`      // UUID of the cart.
	UserUuid    uuid.UUID  `json:"user_uuid" example:"550e8400-e29b-41d4-a716-446655440001"` // UUID of the user associated with the cart.
	ProductUuid string     `json:"product_id" example:"123"`                                 // ID of the product in the cart.
	Quantity    int64      `json:"quantity" example:"2"`                                     // Quantity of the product in the cart.
	CreatedAt   *time.Time `json:"created_at"`                                               // Date and time when the cart was created.
	UpdatedAt   *time.Time `json:"updated_at"`                                               // Date and time when the cart was last updated.
}
