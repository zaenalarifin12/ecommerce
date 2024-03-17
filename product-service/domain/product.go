package domain

import (
	"context"
	"github.com/google/uuid"
	"github.com/guregu/null"
	"github.com/zaenalarifin12/product-service/dto/productDto"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Product represents a product entity
type Product struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `bson:"name,omitempty"`
	Description string             `bson:"description,omitempty"`
	Price       float64            `bson:"price,omitempty"`
	Quantity    int                `bson:"quantity,omitempty"`
	UserUUID    uuid.UUID          `bson:"user_uuid,omitempty"`
	CreatedAt   null.Time          `bson:"created_at,omitempty"`
	UpdatedAt   null.Time          `bson:"updated_at,omitempty"`
	DeletedAt   null.Time          `bson:"deleted_at,omitempty"`
}

type IProductRepository interface {
	List(ctx context.Context, userUUID uuid.UUID, productIDs []primitive.ObjectID) ([]Product, error)
	FindById(ctx context.Context, id primitive.ObjectID, userUUID uuid.UUID) (Product, error)
	Insert(ctx context.Context, req productDto.ProductRequest, userUUID uuid.UUID) (Product, error)
	Update(ctx context.Context, id primitive.ObjectID, req productDto.ProductUpdateRequest, userUUID uuid.UUID) (Product, error)
	Delete(ctx context.Context, id primitive.ObjectID, userUUID uuid.UUID) error
}

type IProductService interface {
	ListProduct(ctx context.Context, userUUID uuid.UUID, productIDs []primitive.ObjectID) ([]productDto.ProductResponse, error)
	GetProductByID(ctx context.Context, id primitive.ObjectID, userUUID uuid.UUID) (productDto.ProductResponse, error)
	AddProduct(ctx context.Context, req productDto.ProductRequest, userUUID uuid.UUID) (productDto.ProductResponse, error)
	UpdateProduct(ctx context.Context, id primitive.ObjectID, req productDto.ProductUpdateRequest, userUUID uuid.UUID) (productDto.ProductUpdateResponse, error)
	RemoveProduct(ctx context.Context, id primitive.ObjectID, userUUID uuid.UUID) error
}
