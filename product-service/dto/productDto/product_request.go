package productDto

import "strings"

// ProductRequest represents a request to update user information.
// swagger:parameters ProductRequest
type ProductRequest struct {
	Name        string  `json:"name,omitempty" bson:"name,omitempty"`
	Description string  `json:"description,omitempty" bson:"description,omitempty"`
	Price       float64 `json:"price,omitempty" validate:"required,min=1000" bson:"price,omitempty"`
	Quantity    int     `json:"quantity,omitempty" validate:"required,min=1" bson:"quantity,omitempty"`
}

// Sanitize the fields of ProductRequest struct
func (req *ProductRequest) Sanitize() {
	// Trim leading and trailing spaces from each field
	req.Name = strings.TrimSpace(req.Name)
	req.Description = strings.TrimSpace(req.Description)
}
