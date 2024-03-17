package domain

import (
	"context"
	"github.com/cart-service/dto"
	"github.com/google/uuid"
	"time"
)

// Cart represents a shopping cart item.
type Cart struct {
	UUID        uuid.UUID  `db:"uuid"`
	UserUUID    uuid.UUID  `db:"user_uuid"`
	ProductUUID string     `db:"product_id"`
	Quantity    int64      `db:"quantity"`
	CreatedAt   *time.Time `db:"created_at"`
	UpdatedAt   *time.Time `db:"updated_at"`
}

// ICartRepository provides methods to interact with the cart data in the database.
type ICartRepository interface {
	ListProducts(ctx context.Context, userUUID uuid.UUID) ([]Cart, error)
	FindByUUIDProduct(ctx context.Context, productUUID string, userUUID uuid.UUID) (Cart, error)
	InsertProductToCart(ctx context.Context, productUUID string, quantity int64, userUuid uuid.UUID) (Cart, error)
	UpdateProductQuantity(ctx context.Context, productUUID string, quantity int64, userUUID uuid.UUID) error
	RemoveProduct(ctx context.Context, productUUID string, userUUID uuid.UUID) error
	TxCart(ctx context.Context, userUUID uuid.UUID) error
	RollBackTxCart(ctx context.Context, userUUID uuid.UUID) error
}

// ICartService provides methods for managing shopping cart operations.
type ICartService interface {
	ListProducts(ctx context.Context, userUUID uuid.UUID) ([]Cart, error)
	AddProduct(ctx context.Context, productUUID string, quantity int64, userUuid uuid.UUID) (dto.CartResponse, error)
	UpdateProductQuantity(ctx context.Context, productUUID string, quantity int64, userUUID uuid.UUID) error
	RemoveProduct(ctx context.Context, productUUID string, userUUID uuid.UUID) error
	TxCart(ctx context.Context, userUUID uuid.UUID) error
	RollBackTxCart(ctx context.Context, userUUID uuid.UUID) error
}
