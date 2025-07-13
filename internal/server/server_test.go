package server_test

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"base-code-go-gin-clean/internal/config"
	"base-code-go-gin-clean/internal/handler/user"
	"base-code-go-gin-clean/internal/server"
	"base-code-go-gin-clean/internal/service/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestServer(t *testing.T) {
	// Setup test configuration
	cfg := &config.Config{
		Server: config.ServerConfig{
			Port:        "8080",
			Environment: "test",
		},
	}

	// Create a test logger
	log := newTestLogger()

	t.Run("new server", func(t *testing.T) {
		srv := server.New(cfg, log, &server.ServerOptions{})
		assert.NotNil(t, srv)
	})
	t.Run("start server", func(t *testing.T) {
		srv := server.New(cfg, log, &server.ServerOptions{})

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		go func() {
			err := srv.Start(ctx)
			assert.NoError(t, err)
		}()

		// Optionally, send real request to confirm server is up (skip if not needed)
		time.Sleep(1 * time.Second)
		resp, err := http.Get("http://localhost:" + cfg.Server.Port + "/swagger/index.html")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("with user handler", func(t *testing.T) {
		mockUserSvc := new(mocks.UserService)

		mockUserSvc.
			On("GetUserByID", mock.Anything, "123").
			Return(nil, fmt.Errorf("invalid user ID format"))

		userHandler := user.NewUserHandler(mockUserSvc)

		srv := server.New(cfg, log, &server.ServerOptions{
			UserHandler: userHandler,
		})
		assert.NotNil(t, srv)

		// Kirim request ke path yang sesuai
		req := httptest.NewRequest("GET", "/api/v1/users/123", nil)
		w := httptest.NewRecorder()
		srv.GetRouter().ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("swagger docs available", func(t *testing.T) {
		srv := server.New(cfg, log, &server.ServerOptions{})
		server := httptest.NewServer(srv.GetRouter())
		defer server.Close()

		resp, err := http.Get(server.URL + "/swagger/index.html")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

// newTestLogger creates a new test logger that discards output
func newTestLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
}
