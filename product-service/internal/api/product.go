package api

import (
	"github.com/gin-gonic/gin"
	"github.com/zaenalarifin12/product-service/domain"
	"github.com/zaenalarifin12/product-service/dto/productDto"
	"github.com/zaenalarifin12/product-service/internal/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

type productApi struct {
	productService domain.IProductService
}

func NewProduct(router *gin.Engine, productService domain.IProductService, middlewares ...gin.HandlerFunc) {
	api := productApi{productService: productService}

	productGroup := router.Group("/products")
	productGroup.Use(middlewares...)
	{
		productGroup.GET("", api.List)
		productGroup.POST("", api.Create)
		productGroup.GET("/:uuid", api.Detail)
		productGroup.PUT("/:uuid", api.Update)
		productGroup.DELETE("/:uuid", api.Delete)
	}
}

// List ProductsList /**
// ProductsList godoc
// @Summary      List all products or filter by product IDs
// @Description  List all products or filter by product IDs
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param        product_ids query []string false "Product IDs to filter by"
// @Success      200  	{object} utils.RespondWithDataJSONSwagger
// @Failure      400   {object} domain.ErrorResponse
// @Failure      500	{object} domain.ErrorResponse
// @Security     BearerAuth
// @Router       /products [get]
func (p productApi) List(ctx *gin.Context) {
	userUUID, _ := utils.GetUserUUID(ctx)

	// If there are any query parameters for product IDs, parse them
	ids, _ := ctx.GetQueryArray("product_ids")

	// Check if ids is nil, if so, proceed without adding any product IDs
	var productIDs []primitive.ObjectID
	if ids != nil {
		for _, id := range ids {
			uuid, err := primitive.ObjectIDFromHex(id)
			if err != nil {
				utils.RespondWithErrorJSON(ctx, http.StatusBadRequest, domain.ErrInvalidRequest, err)
				return
			}
			productIDs = append(productIDs, uuid)
		}
	}

	products, err := p.productService.ListProduct(ctx, userUUID, productIDs)
	if err != nil {
		utils.RespondWithErrorJSON(ctx, http.StatusInternalServerError, domain.ErrInternalServerError, err)
		return
	}

	utils.RespondWithDataJSON(ctx, http.StatusOK, "product list", products)
}

// Detail ProductDetailID /**
// ProductDetailID godoc
// @Summary      Product Detail by ID
// @Description  Product Detail by ID
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param uuid path string true "Product UUID"
// @Success      200  	{object} utils.RespondWithDataJSONSwagger
// @Failure      400	{object} domain.ErrorResponse
// @Failure      404	{object} domain.ErrorResponse
// @Failure      500	{object} domain.ErrorResponse
// @Security  BearerAuth
// @Router       /products/{uuid} [get]
func (p productApi) Detail(ctx *gin.Context) {

	userUUID, _ := utils.GetUserUUID(ctx)
	objectID, err := primitive.ObjectIDFromHex(ctx.Param("uuid"))
	if err != nil {
		utils.RespondWithErrorJSON(ctx, http.StatusBadRequest, domain.ErrBadRequest, err.Error())
		return
	}

	productDetail, err := p.productService.GetProductByID(ctx, objectID, userUUID)
	if err != nil {
		utils.RespondWithErrorJSON(ctx, http.StatusNotFound, domain.ErrProductNotFound, err.Error())
		return
	}

	utils.RespondWithDataJSON(ctx, http.StatusOK, "product created", productDetail)
}

// Update ProductDetailID /**
// ProductUpdateID godoc
// @Summary      Update product detail by ID
// @Description  Update product detail by ID
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param id path string true "Product ID"
// @Param product body productDto.ProductUpdateRequest true "Product details"
// @Success      200  	{object} utils.RespondWithDataJSONSwagger
// @Failure      400	{object} domain.ErrorResponse
// @Failure      404	{object} domain.ErrorResponse
// @Failure      500	{object} domain.ErrorResponse
// @Security  BearerAuth
// @Router       /products/{id} [put]
func (p productApi) Update(ctx *gin.Context) {

	userUUID, _ := utils.GetUserUUID(ctx)

	idParam := ctx.Param("uuid")
	objectID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		utils.RespondWithErrorJSON(ctx, http.StatusBadRequest, domain.ErrBadRequest, err.Error())
		return
	}

	var updatedProduct productDto.ProductUpdateRequest
	if err := ctx.BindJSON(&updatedProduct); err != nil {
		utils.RespondWithErrorJSON(ctx, http.StatusBadRequest, domain.ErrBadRequest, err.Error())
		return
	}

	_, err = p.productService.GetProductByID(ctx, objectID, userUUID)
	if err != nil {
		utils.RespondWithErrorJSON(ctx, http.StatusNotFound, domain.ErrProductNotFound, err.Error())
		return
	}

	product, err := p.productService.UpdateProduct(ctx, objectID, updatedProduct, userUUID)
	if err != nil {
		utils.RespondWithErrorJSON(ctx, http.StatusInternalServerError, domain.ErrInternalServerError, err.Error())
		return
	}

	utils.RespondWithDataJSON(ctx, http.StatusOK, "product updated", product)
}

// Create ProductCreate /**
// ProductCreate godoc
// @Summary      Create a new product
// @Description  Create a new product
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param product body productDto.ProductRequest true "Product details"
// @Success      200  	{object} utils.RespondWithDataJSONSwagger
// @Failure      400	{object} domain.ErrorResponse
// @Failure      500	{object} domain.ErrorResponse
// @Security  BearerAuth
// @Router       /products [post]
func (p productApi) Create(ctx *gin.Context) {

	userUUID, _ := utils.GetUserUUID(ctx)

	var newProduct productDto.ProductRequest
	if err := ctx.BindJSON(&newProduct); err != nil {
		utils.RespondWithErrorJSON(ctx, http.StatusBadRequest, domain.ErrBadRequest, err.Error())
		return
	}
	//
	product, err := p.productService.AddProduct(ctx, newProduct, userUUID)
	if err != nil {
		utils.RespondWithErrorJSON(ctx, http.StatusInternalServerError, domain.ErrInternalServerError, err)
		return
	}

	utils.RespondWithDataJSON(ctx, http.StatusCreated, "product created", product)
}

// Delete ProductDeleteID /**
// ProductDeleteID godoc
// @Summary      Delete product by ID
// @Description  Delete product by ID
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param id path string true "Product ID"
// @Success      200  	{object} productDto.ProductResponse
// @Failure      400	{object} domain.ErrorResponse
// @Failure      404	{object} domain.ErrorResponse
// @Failure      500	{object} domain.ErrorResponse
// @Security  BearerAuth
// @Router       /products/{id} [delete]
func (p productApi) Delete(ctx *gin.Context) {

	userUUID, _ := utils.GetUserUUID(ctx)

	idParam := ctx.Param("uuid")
	objectID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		utils.RespondWithErrorJSON(ctx, http.StatusBadRequest, domain.ErrBadRequest, err)
		return
	}

	err = p.productService.RemoveProduct(ctx, objectID, userUUID)
	if err != nil {
		if err == domain.ErrProductNotFound {
			utils.RespondWithErrorJSON(ctx, http.StatusNotFound, domain.ErrProductNotFound, err)
			return
		}
		utils.RespondWithErrorJSON(ctx, http.StatusInternalServerError, domain.ErrInternalServerError, err)
		return
	}

	utils.RespondWithDataJSON(ctx, http.StatusOK, "product deleted", nil)
}
