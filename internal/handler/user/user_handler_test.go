package user_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"base-code-go-gin-clean/internal/domain/user"
	userHandler "base-code-go-gin-clean/internal/handler/user"
	"base-code-go-gin-clean/internal/service/mocks"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	return r
}

func TestUserHandler_GetUserByID(t *testing.T) {
	mockUserSvc := new(mocks.UserService)
	handler := userHandler.NewUserHandler(mockUserSvc)

	t.Run("success", func(t *testing.T) {
		userID := uuid.New()
		expectedUser := &user.UserResponse{
			ID:        userID,
			Name:      "Test User",
			Email:     "test@example.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mockUserSvc.On("GetUserByID", mock.Anything, userID.String()).Return(expectedUser, nil)

		r := setupRouter()
		r.GET("/users/:id", handler.GetUserByID)

		req, _ := http.NewRequest("GET", "/users/"+userID.String(), nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response struct {
			Status string `json:"status"`
			Code   int    `json:"code"`
			Data   struct {
				ID string `json:"id"`
			} `json:"data"`
		}

		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "success", response.Status)
		assert.Equal(t, userID.String(), response.Data.ID)
		mockUserSvc.AssertExpectations(t)
	})

	t.Run("user not found", func(t *testing.T) {
		userID := uuid.New()
		mockUserSvc.On("GetUserByID", mock.Anything, userID.String()).Return(nil, assert.AnError)

		r := setupRouter()
		r.GET("/users/:id", handler.GetUserByID)

		req, _ := http.NewRequest("GET", "/users/"+userID.String(), nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}
