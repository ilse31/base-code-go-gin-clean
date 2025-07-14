package test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

// SetupTestDB creates a test database connection
func SetupTestDB(t *testing.T) *bun.DB {
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		dsn = "postgresql://postgres:120579@127.0.0.1:5432/postgres?sslmode=disable"
	}

	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db := bun.NewDB(sqldb, pgdialect.New())

	// Verify the connection
	err := db.Ping()
	assert.NoError(t, err, "Failed to connect to test database")

	return db
}

// TeardownTestDB closes the database connection
func TeardownTestDB(t *testing.T, db *bun.DB) {
	if db != nil {
		err := db.Close()
		assert.NoError(t, err, "Failed to close test database connection")
	}
}

// SetupTestRouter creates a test Gin router
func SetupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

// MakeTestRequest creates and executes a test HTTP request
func MakeTestRequest(router *gin.Engine, method, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	return resp
}

// AssertJSONResponse asserts common JSON response patterns
func AssertJSONResponse(t *testing.T, resp *httptest.ResponseRecorder, expectedStatus int) {
	assert.Equal(t, expectedStatus, resp.Code)
	assert.Equal(t, "application/json; charset=utf-8", resp.Header().Get("Content-Type"))
}

// TestRequest contains both the request and response for testing
type TestRequest struct {
	Request  *http.Request
	Response *httptest.ResponseRecorder
}

// MakeTestRequestWithBody creates a test HTTP request with body
func MakeTestRequestWithBody(router *gin.Engine, method, path string, body *bytes.Reader) *TestRequest {
	req, _ := http.NewRequest(method, path, body)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	return &TestRequest{
		Request:  req,
		Response: resp,
	}
}

// MakeJSONBody creates a JSON-encoded request body from the given data
func MakeJSONBody(t *testing.T, data interface{}) *bytes.Reader {
	jsonData, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("Failed to marshal JSON: %v", err)
	}
	return bytes.NewReader(jsonData)
}
