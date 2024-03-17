package repository

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/guregu/null"
	"github.com/zaenalarifin12/product-service/domain"
	"github.com/zaenalarifin12/product-service/dto/productDto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type productRepository struct {
	collection *mongo.Collection
}

func NewProductRepository(db *mongo.Database, collectionName string) domain.IProductRepository {
	return &productRepository{
		collection: db.Collection(collectionName),
	}
}

func (p productRepository) List(ctx context.Context, userUUID uuid.UUID, productIDs []primitive.ObjectID) ([]domain.Product, error) {
	var products []domain.Product
	filter := bson.M{
		"user_uuid": userUUID,
	}

	// Add filter for product IDs only if productIDs is not nil
	if productIDs != nil {
		filter["_id"] = bson.M{"$in": productIDs}
	}

	cursor, err := p.collection.Find(ctx, filter, options.Find().SetSort(bson.D{{"created_at", -1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var product domain.Product
		if err := cursor.Decode(&product); err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func (p productRepository) FindById(ctx context.Context, id primitive.ObjectID, userUUID uuid.UUID) (domain.Product, error) {
	var product domain.Product
	filter := bson.M{"_id": id, "user_uuid": userUUID}
	err := p.collection.FindOne(ctx, filter).Decode(&product)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return product, domain.ErrNotFound
		}
		return product, err
	}
	return product, nil
}

func (p productRepository) Insert(ctx context.Context, req productDto.ProductRequest, userUUID uuid.UUID) (domain.Product, error) {
	product := domain.Product{
		ID:          primitive.NewObjectID(),
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Quantity:    req.Quantity,
		UserUUID:    userUUID,
		CreatedAt:   null.NewTime(time.Now(), true),
		UpdatedAt:   null.Time{},
		DeletedAt:   null.Time{},
	}
	result, err := p.collection.InsertOne(ctx, product)
	if err != nil {
		return domain.Product{}, err
	}
	insertedID := result.InsertedID.(primitive.ObjectID)
	return p.FindById(ctx, insertedID, userUUID)
}

func (p productRepository) Update(ctx context.Context, id primitive.ObjectID, req productDto.ProductUpdateRequest, userUUID uuid.UUID) (domain.Product, error) {

	updateTime := null.NewTime(time.Now(), true)
	filter := bson.M{"_id": id}
	update := bson.M{"$set": domain.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Quantity:    req.Quantity,
		UserUUID:    userUUID,
		UpdatedAt:   updateTime,
	}}
	_, err := p.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return domain.Product{}, err
	}

	// Return the updated product
	updatedProduct, err := p.FindById(ctx, id, userUUID)
	if err != nil {
		return domain.Product{}, err
	}

	return domain.Product{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Quantity:    req.Quantity,
		CreatedAt:   updatedProduct.CreatedAt,
		UpdatedAt:   updatedProduct.UpdatedAt,
	}, nil
}

func (p productRepository) UpdateProduct(ctx context.Context, id primitive.ObjectID, req productDto.ProductRequest, userUUID uuid.UUID) (domain.Product, error) {

	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": bson.M{
			"name":        req.Name,
			"description": req.Description,
			"price":       req.Price,
			"quantity":    req.Quantity,
			// Optionally, you can update timestamps if needed
			"updatedAt": time.Now(),
		},
	}
	_, err := p.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return domain.Product{}, err
	}

	// Return the updated product
	updatedProduct, err := p.FindById(ctx, id, userUUID)
	if err != nil {
		return domain.Product{}, err
	}

	return updatedProduct, nil
}

func (p productRepository) Delete(ctx context.Context, id primitive.ObjectID, userUUID uuid.UUID) error {
	filter := bson.M{"_id": id, "user_uuid": userUUID}
	result, err := p.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return domain.ErrNotFound
	}
	return nil
}
