package api

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"user-service/domain"
	"user-service/dto/userDto"
	"user-service/internal/utils"
)

type userApi struct {
	userService domain.UserService
}

func NewUser(router *gin.Engine, userService domain.UserService) {
	handler := userApi{userService: userService}

	v1 := router.Group("/api/v1")
	{
		v1.GET("/users/:uuid", handler.Detail)
		v1.PUT("/users/:uuid", handler.Update)
	}

}

// Detail UsersDetailUUID /**
// UsersDetailUUID godoc
// @Summary      Users Detail by UUID
// @Description  Users Detail by UUID
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param uuid path int true "User uuid"
// @Success      200  	{object} userDto.DataUserResponse
// @Failure      400	{object} domain.ErrorResponse
// @Failure      404	{object} domain.ErrorResponse
// @Failure      500	{object} domain.ErrorResponse
// @Router       /users/{uuid} [get]
func (u userApi) Detail(ctx *gin.Context) {
	uuidParam := ctx.Param("uuid")
	uuidParse, err := uuid.Parse(uuidParam)
	if err != nil {
		return
	}
	userDetail, err := u.userService.Detail(ctx, uuidParse)

	if err != nil {
		utils.RespondWithErrorJSON(ctx, http.StatusInternalServerError, domain.ErrInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, userDto.DataUserResponse{
		Data: userDetail,
	})
}

// Update UsersDetailUUID /**
// UsersUpdateUUID godoc
// @Summary      Update user detail by UUID
// @Description  Update user detail by UUID
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param uuid path string true "User UUID"
// @Success      200  	{object} userDto.DataUserResponse
// @Failure      400	{object} domain.ErrorResponse
// @Failure      404	{object} domain.ErrorResponse
// @Failure      500	{object} domain.ErrorResponse
// @Router       /users/{uuid} [put]
func (u userApi) Update(ctx *gin.Context) {
	uuidParam := ctx.Param("uuid")
	uuidParse, err := uuid.Parse(uuidParam)
	if err != nil {
		utils.RespondWithErrorJSON(ctx, http.StatusBadRequest, domain.ErrBadRequest, err)
		return
	}

	// Parse the request body to extract updated user information
	var updatedUser userDto.UserUpdateRequest
	if err := ctx.BindJSON(&updatedUser); err != nil {
		utils.RespondWithErrorJSON(ctx, http.StatusBadRequest, domain.ErrBadRequest, err)
		return
	}

	updatedUser.Sanitize()
	validationErrors := utils.ValidateStruct(updatedUser)
	if len(validationErrors) > 0 {
		utils.RespondWithErrorJSON(ctx, http.StatusBadRequest, domain.ErrValidation, validationErrors)
		return
	}

	// Set the UUID of the updated user to match the UUID in the URL
	updatedUser.Uuid = uuidParse

	// Call the UserService to update the user
	updateResponse, err := u.userService.Update(ctx, updatedUser)
	if err != nil {
		return
	}
	if err != nil {
		utils.RespondWithErrorJSON(ctx, http.StatusInternalServerError, domain.ErrInternalServerError, err)
		return
	}

	// Return success response
	ctx.JSON(http.StatusOK, userDto.DataUserUpdateResponse{
		Message: "data user updated successfully",
		Data:    updateResponse,
	})
}
