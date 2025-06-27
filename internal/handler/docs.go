package handler

// Common response structures for Swagger documentation

// SuccessResponse represents a successful API response
type SuccessResponse struct {
	Status string      `json:"status" example:"success"`
	Code   int         `json:"code" example:"200"`
	Data   interface{} `json:"data"`
}

// ErrorResponse represents an error API response
type ErrorResponse struct {
	Status  string      `json:"status" example:"error"`
	Code    int         `json:"code" example:"400"`
	Message string      `json:"message" example:"Bad Request"`
	Errors  interface{} `json:"errors,omitempty"`
}

// ValidationError represents a validation error response
type ValidationError struct {
	Status  string            `json:"status" example:"validation_error"`
	Code    int               `json:"code" example:"422"`
	Message string            `json:"message" example:"Validation failed"`
	Errors  map[string]string `json:"errors"`
}
