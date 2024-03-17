package domain

import (
	"context"
	"time"
)

// Order represents a shopping cart item.
type Cart struct {
	UUID        uuid.UUID  `db:"uuid"`
	UserUUID    uuid.UUID  `db:"user_uuid"`
	ProductUUID uuid.UUID  `db:"product_id"`
	Quantity    int64      `db:"quantity"`
	CreatedAt   *time.Time `db:"created_at"`
	UpdatedAt   *time.Time `db:"updated_at"`
}

// ICartRepository provides methods to interact with the cart data in the database.
type ICartRepository interface {
	ListProducts(ctx context.Context, userUUID uuid.UUID) ([]Cart, error)
	FindByUUID(ctx context.Context, cartUUID uuid.UUID) (Cart, error)
	InsertProductToCart(ctx context.Context, productUUID uuid.UUID, quantity int64, userUuid uuid.UUID) (Cart, error)
	UpdateProductQuantity(ctx context.Context, cartUUID uuid.UUID, quantity int64) error
	RemoveProduct(ctx context.Context, cartUUID uuid.UUID) error
}

// ICartService provides methods for managing shopping cart operations.
type ICartService interface {
	ListProducts(ctx context.Context, userUUID uuid.UUID) ([]Cart, error)
	AddProduct(ctx context.Context, productUUID uuid.UUID, quantity int64, userUuid uuid.UUID) (dto.CartResponse, error)
	UpdateProductQuantity(ctx context.Context, cartUUID uuid.UUID, quantity int64) error
	RemoveProduct(ctx context.Context, cartUUID uuid.UUID) error
}
