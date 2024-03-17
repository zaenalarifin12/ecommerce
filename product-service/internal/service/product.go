package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/zaenalarifin12/product-service/domain"
	"github.com/zaenalarifin12/product-service/dto/productDto"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
)

type productService struct {
	productRepo domain.IProductRepository
}

func NewProductService(repo domain.IProductRepository) domain.IProductService {
	return &productService{
		productRepo: repo,
	}
}

// ListProduct fetches all products or products with specific IDs.
func (s *productService) ListProduct(ctx context.Context, userUUID uuid.UUID, productIDs []primitive.ObjectID) ([]productDto.ProductResponse, error) {
	products, err := s.productRepo.List(ctx, userUUID, productIDs) // Pass productIDs to List function
	if err != nil {
		log.Printf("Error fetching products: %v", err)
		return nil, err
	}

	var productResponses []productDto.ProductResponse
	for _, p := range products {
		productResponses = append(productResponses, productDto.ProductResponse{
			ID:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			Quantity:    p.Quantity,
			CreatedAt:   p.CreatedAt,
			UpdatedAt:   p.UpdatedAt,
		})
	}

	return productResponses, nil
}

// GetProductByID fetches a product by ID.
func (s *productService) GetProductByID(ctx context.Context, id primitive.ObjectID, userUUID uuid.UUID) (productDto.ProductResponse, error) {
	product, err := s.productRepo.FindById(ctx, id, userUUID)
	if err != nil {
		log.Printf("Error fetching product by ID %v: %v", id, err)
		return productDto.ProductResponse{}, err
	}

	return productDto.ProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Quantity:    product.Quantity,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}, nil
}

// AddProduct adds a new product.
func (s *productService) AddProduct(ctx context.Context, req productDto.ProductRequest, userUUID uuid.UUID) (productDto.ProductResponse, error) {
	product, err := s.productRepo.Insert(ctx, req, userUUID)
	if err != nil {
		log.Printf("Error adding product: %v", err)
		return productDto.ProductResponse{}, err
	}
	return productDto.ProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Quantity:    product.Quantity,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}, nil
}

func (s *productService) UpdateProduct(ctx context.Context, id primitive.ObjectID, req productDto.ProductUpdateRequest, userUUID uuid.UUID) (productDto.ProductUpdateResponse, error) {
	// Fetch the existing product by ID
	_, err := s.productRepo.FindById(ctx, id, userUUID)
	if err != nil {
		log.Printf("Error fetching product by ID %v: %v", id, err)
		return productDto.ProductUpdateResponse{}, err
	}

	// Update the existing product with the new details
	updateData := productDto.ProductUpdateRequest{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Quantity:    req.Quantity,
	}
	// Call repository's Update method to update the product in the database
	updatedProduct, err := s.productRepo.Update(ctx, id, updateData, userUUID)
	if err != nil {
		log.Printf("Error updating product: %v", err)
		return productDto.ProductUpdateResponse{}, err
	}

	return productDto.ProductUpdateResponse{
		ID:          updatedProduct.ID,
		Name:        updatedProduct.Name,
		Description: updatedProduct.Description,
		Price:       updatedProduct.Price,
		Quantity:    updatedProduct.Quantity,
		CreatedAt:   updatedProduct.CreatedAt,
		UpdatedAt:   updatedProduct.UpdatedAt,
	}, nil
}

// RemoveProduct removes a product by ID.
func (s *productService) RemoveProduct(ctx context.Context, id primitive.ObjectID, userUUID uuid.UUID) error {
	err := s.productRepo.Delete(ctx, id, userUUID)
	if err != nil {
		log.Printf("Error removing product by ID %v: %v", id, err)
		return domain.ErrProductNotFound
	}
	return nil
}
