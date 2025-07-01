package server_test

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"base-code-go-gin-clean/internal/config"
	"base-code-go-gin-clean/internal/handler/user"
	"base-code-go-gin-clean/internal/server"
	"base-code-go-gin-clean/internal/service/mocks"

	"github.com/stretchr/testify/assert"
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
		srv := server.New(cfg, log)
		assert.NotNil(t, srv)
	})

	t.Run("start server", func(t *testing.T) {
		// Create a test server with a simple handler
		srv := server.New(cfg, log)
		assert.NotNil(t, srv)

		// Create a channel to signal when server is ready
		serverReady := make(chan bool)

		// Start the server in a goroutine
		go func() {
			if err := srv.Start(); err != nil && err != http.ErrServerClosed {
				t.Errorf("Server error: %v", err)
			}
			serverReady <- true
		}()

		// Wait for server to be ready
		<-serverReady

		// Get the server instance and close it
		httpSrv := srv.GetHTTPServer()
		assert.NotNil(t, httpSrv, "HTTP server should be initialized after Start()")

		// Stop the server
		err := httpSrv.Close()
		assert.NoError(t, err)
	})

	t.Run("with user handler", func(t *testing.T) {
		mockUserSvc := new(mocks.UserService)
		userHandler := user.NewUserHandler(mockUserSvc)

		srv := server.New(cfg, log, server.WithUserHandler(userHandler))
		assert.NotNil(t, srv)

		// Test the router directly
		req := httptest.NewRequest("GET", "/users/123", nil)
		w := httptest.NewRecorder()
		srv.GetRouter().ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("swagger docs available", func(t *testing.T) {
		srv := server.New(cfg, log)
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
