package user

import (
	"context"
	"errors"

	"base-code-go-gin-clean/internal/handler/user/dto"
	"base-code-go-gin-clean/internal/pkg/http"
	"base-code-go-gin-clean/internal/pkg/telemetry"
	"base-code-go-gin-clean/internal/service"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// GetUserByID handles user retrieval by ID
// @Summary Get user by ID
// @Description Get user by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} handler.SuccessResponse{data=dto.UserResponse} "Success response with user data"
// @Failure 400 {object} handler.ErrorResponse "Bad Request: Invalid user ID"
// @Failure 401 {object} handler.ErrorResponse "Unauthorized: Authentication required"
// @Failure 403 {object} handler.ErrorResponse "Forbidden: Insufficient permissions"
// @Failure 404 {object} handler.ErrorResponse "Not Found: User not found"
// @Failure 500 {object} handler.ErrorResponse "Internal Server Error"
// @Router /users/{id} [get]
func (h *UserHandler) GetUserByID(c *gin.Context) {
	ctx := c.Request.Context()

	// Get user ID from path
	id := c.Param("id")
	if id == "" {
		err := errors.New("user ID is required")
		telemetry.RecordError(ctx, err)
		http.BadRequest(c, "User ID is required", nil)
		return
	}

	// Use WithSpan to automatically handle span creation and error recording
	err := telemetry.WithSpan(ctx, "GetUserByID", func(ctx context.Context) error {
		// Add user ID to span
		telemetry.AddSpanAttributes(ctx, attribute.String("user.id", id))

		// Call service
		userResponse, err := h.userService.GetUserByID(ctx, id)
		if err != nil {
			telemetry.AddSpanAttributes(ctx, attribute.String("error.type", "user_not_found"))
			http.NotFound(c, "User not found")
			return err
		}

		// Map domain model to DTO and return success response
		http.Success(c, dto.NewUserResponse(userResponse))
		return nil
	})

	// Error is already recorded in the span by WithSpan
	if err != nil {
		// Error response already sent in the WithSpan callback
		return
	}
}
