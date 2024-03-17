package domain

import "errors"

var (
	ErrInternalServerError = errors.New("internal server error")
	ErrCartNotFound        = errors.New("cart not found")
	ErrProductNotFound     = errors.New("product not found")
	ErrUsernameIsTaken     = errors.New("Username is taken")
	ErrNotFound            = errors.New("Data not found")
	ErrListError           = errors.New("Error occurred while listing data")
	ErrPermissionDenied    = errors.New("Permission denied")
	ErrInvalidInput        = errors.New("Invalid input provided")
	ErrDatabaseError       = errors.New("Database error occurred")
	ErrUnauthorizedAccess  = errors.New("Unauthorized access")
	ErrInvalidCredentials  = errors.New("Invalid credentials provided")
	ErrValidationFailed    = errors.New("Validation failed")
	ErrTimeout             = errors.New("Operation timed out")
	ErrDuplicateEntry      = errors.New("Duplicate entry")
	ErrorRateLimit         = errors.New("Rate limit exceeded")
	ErrBadRequest          = errors.New("Bad Request")
	ErrValidation          = errors.New("Validation Error")
	ErrFailedInsert        = errors.New("Failed Insert Data")
)
