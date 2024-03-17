package dto

// CartRequest represents the request body for adding a product to the cart.
// swagger:parameters CartRequest
type CartRequest struct {
	ProductUuid string `json:"product_id" example:"6ba7b810-9dad-11d1-80b4-00c04fd430c8"` // Product ID of the item to be added to the cart.
	Quantity    int64  `json:"quantity" example:"1"`                                      // Quantity of the item to be added.
}
