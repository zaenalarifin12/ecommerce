package dto

// CartRequestUpdate represents the request body for adding a product to the cart.
// swagger:parameters CartRequestUpdate
type CartRequestUpdate struct {
	Quantity int64 `json:"quantity" example:"1"` // Quantity of the item to be added.
}
