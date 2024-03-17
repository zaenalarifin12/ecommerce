package utils

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetUserUUID(ctx *gin.Context) (uuid.UUID, error) {
	// Get user UUID string from context
	userUUIDString, ok := ctx.Get("user_uuid")
	if !ok {
		return uuid.Nil, fmt.Errorf("user_uuid not found in context")
	}

	// Parse UUID string to UUID
	userUUID, err := uuid.Parse(userUUIDString.(string))
	if err != nil {
		return uuid.Nil, err
	}

	return userUUID, nil
}

func ConvertStringToUUID(param any) uuid.UUID {
	// Get user UUID string from context

	// Parse UUID string to UUID
	result, err := uuid.Parse(param.(string))
	if err != nil {
		fmt.Println(err.Error())
		return uuid.Nil
	}

	return result
}
