package http

import "github.com/gin-gonic/gin"

// HTTP status code constants
const (
	StatusSuccess    = "success"
	StatusError      = "error"
	StatusValidation = "validation_error"
)

// Response represents the standard API response structure
type Response struct {
	Status  string      `json:"status"`
	Code    int         `json:"code"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

// SuccessResponse creates a success response
func SuccessResponse(c *gin.Context, code int, data interface{}) {
	c.JSON(code, Response{
		Status: StatusSuccess,
		Code:   code,
		Data:   data,
	})
}

// ErrorResponse creates an error response
func ErrorResponse(c *gin.Context, code int, message string, errors interface{}) {
	status := StatusError
	if code == 422 { // Validation error
		status = StatusValidation
	}

	c.JSON(code, Response{
		Status:  status,
		Code:    code,
		Message: message,
		Errors:  errors,
	})
}

// Success creates a success response with status code 200
func Success(c *gin.Context, data interface{}) {
	SuccessResponse(c, 200, data)
}

// Created creates a success response with status code 201
func Created(c *gin.Context, data interface{}) {
	SuccessResponse(c, 201, data)
}

// BadRequest creates an error response with status code 400
func BadRequest(c *gin.Context, message string, errors interface{}) {
	ErrorResponse(c, 400, message, errors)
}

// Unauthorized creates an error response with status code 401
func Unauthorized(c *gin.Context, message string) {
	ErrorResponse(c, 401, message, nil)
}

// Forbidden creates an error response with status code 403
func Forbidden(c *gin.Context, message string) {
	ErrorResponse(c, 403, message, nil)
}

// NotFound creates an error response with status code 404
func NotFound(c *gin.Context, message string) {
	ErrorResponse(c, 404, message, nil)
}

// ValidationError creates a validation error response with status code 422
func ValidationError(c *gin.Context, message string, errors interface{}) {
	ErrorResponse(c, 422, message, errors)
}

// InternalServerError creates an error response with status code 500
func InternalServerError(c *gin.Context, message string) {
	ErrorResponse(c, 500, message, nil)
}
