package user

import (
	"base-code-go-gin-clean/internal/handler/user/dto"
	"base-code-go-gin-clean/internal/pkg/http"
	"base-code-go-gin-clean/internal/service"

	"github.com/gin-gonic/gin"
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
	id := c.Param("id")
	if id == "" {
		http.BadRequest(c, "User ID is required", nil)
		return
	}

	userResponse, err := h.userService.GetUserByID(c.Request.Context(), id)
	if err != nil {
		http.NotFound(c, "User not found")
		return
	}

	// Map domain model to DTO
	response := dto.NewUserResponse(userResponse)
	http.Success(c, response)
}
