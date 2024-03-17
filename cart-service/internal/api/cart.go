package api

import (
	"fmt"
	"github.com/cart-service/domain"
	"github.com/cart-service/dto"
	"github.com/cart-service/internal/client"
	"github.com/cart-service/internal/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

type CartApi struct {
	CartService domain.ICartService
}

func NewCart(router *gin.Engine, cartService domain.ICartService, middlewares ...gin.HandlerFunc) {

	api := CartApi{CartService: cartService}

	cartGroup := router.Group("/carts")
	cartGroup.Use(middlewares...)
	{
		cartGroup.GET("", api.ListProducts)
		cartGroup.POST("", api.AddProduct)
		cartGroup.PUT("/:uuid", api.UpdateProductQuantity)
		cartGroup.DELETE("/:uuid", api.RemoveProduct)
		cartGroup.POST("/tx", api.TxCart)
		cartGroup.POST("/tx-rollback", api.TxCartRollback)
	}
}

// ListProducts @Summary List products in cart
// @Description List all products in the cart
// @Tags Carts
// @Accept json
// @Produce json
// @Success 200 {array} domain.Cart
// @Failure 500 {object} domain.ErrorResponse
// @Security  BearerAuth
// @Router /carts [get]
func (c *CartApi) ListProducts(ctx *gin.Context) {

	userUUID, _ := util.GetUserUUID(ctx)

	products, err := c.CartService.ListProducts(ctx, userUUID)

	if products == nil {
		products = []domain.Cart{}
	}

	if err != nil {
		fmt.Println(err, "err")
		util.RespondWithErrorJSON(ctx, http.StatusInternalServerError, domain.ErrInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, products)
}

// AddProduct @Summary Add product to cart
// @Description Add a product to the cart
// @Tags Carts
// @Accept json
// @Produce json
// @Param request body dto.CartRequest true "Product details"
// @Success 201 {string} string "Product added successfully"
// @Failure 400 {object} domain.ErrorResponse
// @Security  BearerAuth
// @Router /carts [post]
func (c *CartApi) AddProduct(ctx *gin.Context) {

	userUUID, _ := util.GetUserUUID(ctx)
	token, _ := ctx.Get("token")

	var req dto.CartRequest

	if err := ctx.BindJSON(&req); err != nil {
		util.RespondWithErrorJSON(ctx, http.StatusBadRequest, domain.ErrBadRequest, err.Error())
		return
	}

	// check hit product
	response, err := client.APIProductClient.Get(client.ProductByUUID(req.ProductUuid), token.(string))
	if err != nil {
		util.RespondWithErrorJSON(ctx, http.StatusBadRequest, domain.ErrBadRequest, err.Error())
		return
	}

	// if product not exist print error
	if response.StatusCode != http.StatusOK {
		util.RespondWithErrorJSON(ctx, http.StatusBadRequest, domain.ErrBadRequest, domain.ErrProductNotFound)
		return
	}
	dataCart, err := c.CartService.AddProduct(ctx, req.ProductUuid, req.Quantity, userUUID)
	if err != nil {
		util.RespondWithErrorJSON(ctx, http.StatusBadRequest, domain.ErrBadRequest, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, dataCart)
}

// UpdateProductQuantity @Summary Update product quantity in cart
// @Description Update quantity of a product in the cart
// @Tags Carts
// @Accept json
// @Produce json
// @Param uuid path string true "Product UUID"
// @Param request body dto.CartRequest true "Product quantity"
// @Success 200 {string} string "Product quantity updated successfully"
// @Failure 400 {object} domain.ErrorResponse
// @Security  BearerAuth
// @Router /carts/{uuid} [put]
func (c *CartApi) UpdateProductQuantity(ctx *gin.Context) {
	uuidParam := ctx.Param("uuid")
	productUUID := uuidParam
	userUUID, _ := util.GetUserUUID(ctx)

	var req dto.CartRequestUpdate
	if err := ctx.BindJSON(&req); err != nil {
		util.RespondWithErrorJSON(ctx, http.StatusBadRequest, domain.ErrBadRequest, err.Error())
		return
	}

	err := c.CartService.UpdateProductQuantity(ctx, productUUID, req.Quantity, userUUID)
	if err != nil {
		util.RespondWithErrorJSON(ctx, http.StatusBadRequest, domain.ErrBadRequest, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, "Product quantity updated successfully")
}

// RemoveProduct @Summary Remove product from cart
// @Description Remove a product from the cart
// @Tags Carts
// @Accept json
// @Produce json
// @Param uuid path string true "Product UUID"
// @Success 200 {string} string "Product removed successfully"
// @Failure 400 {object} domain.ErrorResponse
// @Security  BearerAuth
// @Router /carts/{uuid} [delete]
func (c *CartApi) RemoveProduct(ctx *gin.Context) {
	uuidParam := ctx.Param("uuid")
	productUUID := uuidParam
	userUUID, _ := util.GetUserUUID(ctx)

	var err error
	if err != nil {
		util.RespondWithErrorJSON(ctx, http.StatusBadRequest, domain.ErrBadRequest, err)
		return
	}

	err = c.CartService.RemoveProduct(ctx, productUUID, userUUID)
	if err != nil {
		util.RespondWithErrorJSON(ctx, http.StatusInternalServerError, domain.ErrInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, "Product quantity removed successfully")
}

// TxCart @Summary Tx from cart
// @Description Tx cart
// @Tags TransactionCarts
// @Accept json
// @Produce json
// @Success 200 {string} string "Product tx successfully"
// @Failure 400 {object} domain.ErrorResponse
// @Security  BearerAuth
// @Router /carts/tx [post]
func (c *CartApi) TxCart(ctx *gin.Context) {
	userUUID, _ := util.GetUserUUID(ctx)

	err := c.CartService.TxCart(ctx, userUUID)
	if err != nil {
		util.RespondWithErrorJSON(ctx, http.StatusInternalServerError, domain.ErrInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, "Product tx successfully")
}

// TxCartRollback @Summary Tx from cart
// @Description Tx Cart Rollback
// @Tags TransactionCarts
// @Accept json
// @Produce json
// @Success 200 {string} string "Rollback Tx Successfully"
// @Failure 400 {object} domain.ErrorResponse
// @Security  BearerAuth
// @Router /carts/tx-rollback [post]
func (c *CartApi) TxCartRollback(ctx *gin.Context) {
	userUUID, _ := util.GetUserUUID(ctx)

	err := c.CartService.RollBackTxCart(ctx, userUUID)
	if err != nil {
		util.RespondWithErrorJSON(ctx, http.StatusInternalServerError, domain.ErrInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, "Product Tx Successfully")
}
