package api

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"net/http"
	"user-service/domain"
	"user-service/dto/authDto"
	"user-service/internal/utils"
)

type authApi struct {
	userService domain.UserService
	limiter     *rate.Limiter
}

func NewAuth(router *gin.Engine, userService domain.UserService, r rate.Limit, b int) {
	handler := authApi{userService: userService,
		limiter: rate.NewLimiter(r, b),
	}

	v1 := router.Group("/api/v1")
	{
		v1.POST("/auth/register", handler.rateLimitMiddleware, handler.Register)
		v1.POST("/auth/login", handler.rateLimitMiddleware, handler.Authenticate)
	}

}

func (a *authApi) rateLimitMiddleware(ctx *gin.Context) {
	if !a.limiter.Allow() {
		utils.RespondWithErrorJSON(ctx, http.StatusTooManyRequests, domain.ErrorRateLimit, domain.ErrorRateLimit)
		ctx.Abort()
		return
	}
}

// Register /**
// Register godoc
// @Summary      Register User
// @Description  Register User
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param register body authDto.RegisterRequest true "ExampleValue"
// @Success      200  	{object} authDto.DataRegisterResponse
// @Failure      400	{object} domain.ErrorResponse
// @Failure      404	{object} domain.ErrorResponse
// @Failure      500	{object} domain.ErrorResponse
// @Router       /auth/register [post]
func (a authApi) Register(ctx *gin.Context) {
	var req authDto.RegisterRequest

	if err := ctx.BindJSON(&req); err != nil {
		utils.RespondWithErrorJSON(ctx, http.StatusBadRequest, domain.ErrBadRequest, domain.ErrBadRequest)
		return
	}

	req.Sanitize()
	validationErrors := utils.ValidateStruct(req)
	if len(validationErrors) > 0 {
		utils.RespondWithErrorJSON(ctx, http.StatusBadRequest, domain.ErrValidation, validationErrors)
		return
	}

	res, err := a.userService.Register(ctx, req)
	if err != nil {
		utils.RespondWithErrorJSON(ctx, http.StatusInternalServerError, domain.ErrFailedInsert, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, authDto.DataRegisterResponse{Data: res})
}

// Authenticate /**
// Authenticate godoc
// @Summary      Authenticate User
// @Description  Authenticate User
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param register body authDto.LoginRequest true "ExampleValue"
// @Success      200  	{object} authDto.DataLoginResponse
// @Failure      400	{object} domain.ErrorResponse
// @Failure      404	{object} domain.ErrorResponse
// @Failure      500	{object} domain.ErrorResponse
// @Router       /auth/login [post]
func (a authApi) Authenticate(ctx *gin.Context) {
	var req authDto.LoginRequest

	if err := ctx.BindJSON(&req); err != nil {
		utils.RespondWithErrorJSON(ctx, http.StatusBadRequest, domain.ErrBadRequest, domain.ErrBadRequest)
		return
	}

	req.Sanitize()
	validationErrors := utils.ValidateStruct(req)
	if len(validationErrors) > 0 {
		utils.RespondWithErrorJSON(ctx, http.StatusBadRequest, domain.ErrValidation, validationErrors)
		return
	}

	res, err := a.userService.Authenticate(ctx, req)
	if err != nil {
		utils.RespondWithErrorJSON(ctx, http.StatusInternalServerError, domain.ErrInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": res})
}
