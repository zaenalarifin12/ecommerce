package repository

import (
	"context"
	"fmt"
	"github.com/cart-service/domain"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type cartRepository struct {
	DB          *pgxpool.Pool
	RedisClient *redis.Client
}

func NewCart(db *pgxpool.Pool, redisClient *redis.Client) domain.ICartRepository {
	return &cartRepository{
		DB:          db,
		RedisClient: redisClient,
	}
}

func (cr *cartRepository) ListProducts(ctx context.Context, userUUID uuid.UUID) ([]domain.Cart, error) {
	rows, err := cr.DB.Query(ctx, "SELECT uuid, user_uuid, product_uuid, quantity, created_at, updated_at FROM carts WHERE user_uuid = $1", userUUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var carts []domain.Cart
	for rows.Next() {
		var cart domain.Cart
		if err := rows.Scan(&cart.UUID, &cart.UserUUID, &cart.ProductUUID, &cart.Quantity, &cart.CreatedAt, &cart.UpdatedAt); err != nil {
			return nil, err
		}
		carts = append(carts, cart)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return carts, nil
}

func (cr *cartRepository) FindByUUIDProduct(ctx context.Context, productUUID string, userUUID uuid.UUID) (domain.Cart, error) {
	var cart domain.Cart
	row := cr.DB.QueryRow(ctx, "SELECT uuid, user_uuid, product_uuid, quantity, created_at, updated_at FROM carts WHERE product_uuid=$1 AND user_uuid=$2", productUUID, userUUID)
	err := row.Scan(&cart.UUID, &cart.UserUUID, &cart.ProductUUID, &cart.Quantity, &cart.CreatedAt, &cart.UpdatedAt)
	if err == pgx.ErrNoRows {
		return cart, domain.ErrCartNotFound
	}
	if err != nil {
		return cart, err
	}
	return cart, nil
}

func (cr *cartRepository) InsertProductToCart(ctx context.Context, productUUID string, quantity int64, userUuid uuid.UUID) (domain.Cart, error) {
	u := uuid.New()
	timeNow := time.Now()
	_, err := cr.DB.Exec(ctx, "INSERT INTO carts (uuid, product_uuid, quantity, user_uuid, created_at) VALUES ($1, $2,$3,$4, $5)", u, productUUID, quantity, userUuid, timeNow)
	if err != nil {
		return domain.Cart{}, err
	}
	return domain.Cart{
		UUID:        u,
		UserUUID:    userUuid,
		ProductUUID: productUUID,
		Quantity:    quantity,
		CreatedAt:   &timeNow,
	}, nil
}

func (cr *cartRepository) UpdateProductQuantity(ctx context.Context, productUUID string, quantity int64, userUUID uuid.UUID) error {
	_, err := cr.DB.Exec(ctx, "UPDATE carts SET quantity=$1 WHERE product_uuid=$2 AND user_uuid=$3", quantity, productUUID, userUUID)
	if err != nil {
		return err
	}
	return nil
}

func (cr *cartRepository) RemoveProduct(ctx context.Context, productUUID string, userUUID uuid.UUID) error {
	_, err := cr.DB.Exec(ctx, "DELETE FROM carts WHERE product_uuid=$1 AND user_uuid=$2", productUUID, userUUID)
	if err != nil {
		return err
	}
	return nil
}

func (cr *cartRepository) TxCart(ctx context.Context, userUUID uuid.UUID) error {
	tx, err := cr.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback(ctx)
			panic(p) // Re-throw panic after rolling back the transaction
		} else if err != nil {
			tx.Rollback(ctx)
		} else {
			err = tx.Commit(ctx)
		}
	}()

	products, err := cr.ListProducts(ctx, userUUID)
	if err != nil {
		return err
	}

	// save to redis
	// Initialize a map to store field-value pairs
	fieldValueMap := make(map[string]interface{})

	// Iterate over products and populate fieldValueMap
	for _, product := range products {
		field := product.ProductUUID
		value := product.Quantity

		// Print for debugging
		fmt.Println(product.ProductUUID, product.Quantity)

		// Add field-value pair to the map
		fieldValueMap[field] = value
	}

	// Construct the key for the hash
	key := fmt.Sprintf("user:%s:products", userUUID)

	// Use HMSET command to set multiple field-value pairs in the hash
	err = cr.RedisClient.HMSet(ctx, key, fieldValueMap).Err()
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, "DELETE FROM carts WHERE user_uuid=$1", userUUID)
	if err != nil {
		return err
	}
	return nil
}

func (cr *cartRepository) RollBackTxCart(ctx context.Context, userUUID uuid.UUID) error {
	key := fmt.Sprintf("user:%s:products", userUUID)

	// key is the name of the hash you want to retrieve
	//values, err := cr.RedisClient.HGetAll(ctx, key).Result()
	//if err != nil {
	//	return err
	//} else {
	//	// values is a map[string]string containing field-value pairs
	//	for field, value := range values {
	//		fmt.Printf("Field: %s, Value: %s\n", field, value)
	//	}
	//}

	//====
	//Retrieve field-value pairs from Redis
	values, err := cr.RedisClient.HGetAll(ctx, key).Result()
	if err != nil {
		return err
	}

	// Prepare the bulk insert statement
	stmt := `INSERT INTO carts (uuid, product_uuid, quantity, user_uuid, created_at) VALUES `
	var args []interface{}
	i := 0
	for productUUID, quantity := range values {
		if i > 0 {
			stmt += ","
		}
		stmt += fmt.Sprintf("($%d, $%d, $%d,$%d, $%d)", i*5+1, i*5+2, i*5+3, i*5+4, i*5+5)

		args = append(args, uuid.New(), productUUID, quantity, userUUID, time.Now())
		i++
	}

	// Execute the bulk insert statement
	fmt.Println(args)
	_, err = cr.DB.Exec(ctx, stmt, args...)
	if err != nil {
		return err
	}

	// Remove the entire hash from Redis after successful insertion
	_, err = cr.RedisClient.Del(ctx, key).Result()
	if err != nil {
		return err
	}

	return nil
}
