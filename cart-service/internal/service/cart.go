package service

import (
	"context"
	"errors"
	"github.com/cart-service/domain"
	"github.com/cart-service/dto"
	"github.com/google/uuid"
)

type CartService struct {
	Repo domain.ICartRepository
}

func NewCartService(repo domain.ICartRepository) domain.ICartService {
	return &CartService{
		Repo: repo,
	}
}

func (cs *CartService) ListProducts(ctx context.Context, userUUID uuid.UUID) ([]domain.Cart, error) {
	return cs.Repo.ListProducts(ctx, userUUID)
}

func (cs *CartService) AddProduct(ctx context.Context, productUUID string, quantity int64, userUuid uuid.UUID) (dto.CartResponse, error) {
	// Check if product already exists in the cart
	_, err := cs.Repo.FindByUUIDProduct(ctx, productUUID, userUuid)
	if err == nil {
		return dto.CartResponse{}, errors.New("product already exists in the cart")
	}
	if err != domain.ErrCartNotFound {
		return dto.CartResponse{}, err
	}

	// Product does not exist, so add it to the cart
	result, err := cs.Repo.InsertProductToCart(ctx, productUUID, quantity, userUuid)
	if err != nil {
		return dto.CartResponse{}, err
	}

	// Update product quantity if needed
	if quantity > 0 {
		err = cs.Repo.UpdateProductQuantity(ctx, productUUID, quantity, userUuid)
		if err != nil {
			return dto.CartResponse{}, err
		}
	}

	return dto.CartResponse{
		Uuid:        result.UUID,
		UserUuid:    result.UserUUID,
		ProductUuid: result.ProductUUID,
		Quantity:    result.Quantity,
		CreatedAt:   result.CreatedAt,
	}, nil
}

func (cs *CartService) UpdateProductQuantity(ctx context.Context, productUUID string, quantity int64, userUUID uuid.UUID) error {
	// Check if product exists in the cart
	_, err := cs.Repo.FindByUUIDProduct(ctx, productUUID, userUUID)
	if err != nil {
		return err
	}

	// Update product quantity
	err = cs.Repo.UpdateProductQuantity(ctx, productUUID, quantity, userUUID)
	if err != nil {
		return err
	}

	return nil
}

func (cs *CartService) RemoveProduct(ctx context.Context, productUUID string, userUUID uuid.UUID) error {
	// Check if product exists in the cart
	_, err := cs.Repo.FindByUUIDProduct(ctx, productUUID, userUUID)
	if err != nil {
		return err
	}

	// Remove product from the cart
	err = cs.Repo.RemoveProduct(ctx, productUUID, userUUID)
	if err != nil {
		return err
	}

	return nil
}

func (cs *CartService) TxCart(ctx context.Context, userUUID uuid.UUID) error {
	err := cs.Repo.TxCart(ctx, userUUID)
	if err != nil {
		return err
	}

	return nil
}

func (cs *CartService) RollBackTxCart(ctx context.Context, userUUID uuid.UUID) error {
	err := cs.Repo.RollBackTxCart(ctx, userUUID)
	if err != nil {
		return err
	}

	return nil
}
