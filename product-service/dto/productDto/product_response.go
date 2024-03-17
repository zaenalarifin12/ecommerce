package productDto

import (
	"github.com/guregu/null"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ProductResponse represents a request to update user information.
// swagger:parameters ProductResponse
type ProductResponse struct {
	ID          primitive.ObjectID `json:"product_id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name,omitempty" bson:"name,omitempty"`
	Description string             `json:"description,omitempty" bson:"description,omitempty"`
	Price       float64            `json:"price,omitempty" bson:"price,omitempty"`
	Quantity    int                `json:"quantity,omitempty" bson:"quantity,omitempty"`
	CreatedAt   null.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt   null.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}
