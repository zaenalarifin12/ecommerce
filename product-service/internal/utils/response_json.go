package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/zaenalarifin12/product-service/domain"
)

// RespondWithErrorJSON responds with an error message in JSON format.
func RespondWithErrorJSON(c *gin.Context, statusCode int, errorMessage error, errorDetail interface{}) {

	var errorDetailArray []interface{}
	switch v := errorDetail.(type) {
	case []interface{}:
		errorDetailArray = v
	default:
		errorDetailArray = []interface{}{errorDetail}
	}

	errorData := domain.ErrorResponse{
		Error: struct {
			Code    int         `json:"code"`
			Message string      `json:"message"`
			Detail  interface{} `json:"detail"`
		}{
			Code:    statusCode,
			Message: errorMessage.Error(),
			Detail:  errorDetailArray,
		},
	}
	c.JSON(statusCode, errorData)
}

// RespondWithDataJSON responds with data in JSON format.
func RespondWithDataJSON(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, gin.H{
		"message": message,
		"data":    data,
	})
	return
}

// RespondWithDataJSONSwagger responds with data in JSON format.
type RespondWithDataJSONSwagger struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
