package httplog

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"base-code-go-gin-clean/internal/domain/httplog"
)

// contextKey is a custom type for context keys
type contextKey string

// TraceIDKey is the key used to store the trace ID in the context
const TraceIDKey contextKey = "traceID"

// Config holds the configuration for the HTTP logger middleware
type Config struct {
	Service             httplog.Service
	SkipPaths           []string
	SkipHeaders         []string
	SkipBodyMethods     map[string]bool
	MaxBodySize         int64
	IncludeResponseBody bool
}

// DefaultConfig returns the default configuration
func DefaultConfig(service httplog.Service) Config {
	return Config{
		Service:     service,
		SkipPaths:   []string{"/health", "/metrics"},
		SkipHeaders: []string{"Authorization", "Cookie"},
		SkipBodyMethods: map[string]bool{
			"GET":     true,
			"HEAD":    true,
			"OPTIONS": true,
		},
		MaxBodySize:         1024 * 1024,
		IncludeResponseBody: true,
	}
}

// Middleware returns a new HTTP logger middleware
func Middleware(config Config) gin.HandlerFunc {
	if config.Service == nil {
		panic("httplog: Service is required")
	}

	skipPaths := make(map[string]struct{}, len(config.SkipPaths))
	for _, path := range config.SkipPaths {
		skipPaths[path] = struct{}{}
	}

	skipHeaders := make(map[string]struct{}, len(config.SkipHeaders))
	for _, header := range config.SkipHeaders {
		skipHeaders[strings.ToLower(header)] = struct{}{}
	}

	return func(c *gin.Context) {
		if _, ok := skipPaths[c.Request.URL.Path]; ok {
			c.Next()
			return
		}

		traceID := c.GetHeader("X-Trace-ID")
		if traceID == "" {
			traceID = uuid.New().String()
		}
		ctx := context.WithValue(c.Request.Context(), TraceIDKey, traceID)
		c.Request = c.Request.WithContext(ctx)
		c.Writer.Header().Set("X-Trace-ID", traceID)

		reqLog, err := logIncomingRequest(c, config, traceID)
		if err != nil {
			c.Next()
			return
		}

		var blw *bodyLogWriter
		if !config.SkipBodyMethods[c.Request.Method] && config.IncludeResponseBody {
			blw = &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
			c.Writer = blw
		}

		defer func() {
			logResponse(c, blw, reqLog, config, traceID)
		}()

		c.Next()

		if len(c.Errors) > 0 {
			errMsg := ""
			for _, e := range c.Errors {
				errMsg += e.Error() + "; "
			}

			logErr := &httplog.LogError{
				TraceID:    traceID,
				RequestID:  &reqLog.ID,
				StatusCode: c.Writer.Status(),
				Error:      strings.TrimSuffix(errMsg, "; "),
			}
			if err := config.Service.LogError(c.Request.Context(), logErr); err != nil {
				log.Printf("httplog: failed to log error: %v", err)
			}
		}
	}
}

// logIncomingRequest logs the incoming HTTP request
func logIncomingRequest(c *gin.Context, config Config, traceID string) (*httplog.LogIncomingRequest, error) {
	var requestBody interface{} = nil
	if !config.SkipBodyMethods[c.Request.Method] && c.Request.Body != nil {
		bodyBytes, _ := io.ReadAll(io.LimitReader(c.Request.Body, config.MaxBodySize))
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		if strings.Contains(c.ContentType(), "application/json") {
			var jsonBody interface{}
			if err := json.Unmarshal(bodyBytes, &jsonBody); err == nil {
				requestBody = jsonBody
			} else {
				requestBody = string(bodyBytes)
			}
		} else {
			requestBody = string(bodyBytes)
		}
	}

	headers := make(map[string]string)
	for name, values := range c.Request.Header {
		if !contains(config.SkipHeaders, strings.ToLower(name)) {
			headers[name] = strings.Join(values, ", ")
		}
	}

	logEntry := &httplog.LogIncomingRequest{
		TraceID:   traceID,
		EventName: c.FullPath(),
		Endpoint:  c.Request.URL.Path,
		Method:    httplog.HTTPMethod(c.Request.Method),
		Request: map[string]interface{}{
			"method":  c.Request.Method,
			"url":     c.Request.URL.String(),
			"headers": headers,
			"body":    requestBody,
		},
		IPAddress: c.ClientIP(),
		UserAgent: c.Request.UserAgent(),
	}

	logID, err := config.Service.LogIncomingRequest(c.Request.Context(), logEntry)
	if err != nil {
		return nil, err
	}

	logEntry.ID = logID
	return logEntry, nil
}

// logResponse logs the HTTP response
func logResponse(c *gin.Context, blw *bodyLogWriter, reqLog *httplog.LogIncomingRequest, config Config, traceID string) {
	var responseBody interface{} = nil

	if config.IncludeResponseBody && blw != nil && blw.body != nil {
		bodyBytes := blw.body.Bytes()
		if len(bodyBytes) > 0 {
			if strings.Contains(c.Writer.Header().Get("Content-Type"), "application/json") {
				var jsonBody interface{}
				if err := json.Unmarshal(bodyBytes, &jsonBody); err == nil {
					responseBody = jsonBody
				} else {
					responseBody = string(bodyBytes)
				}
			} else {
				responseBody = string(bodyBytes)
			}
		}
	}

	logEntry := &httplog.LogOutgoingRequest{
		TraceID:   traceID,
		EventName: reqLog.EventName,
		Endpoint:  reqLog.Endpoint,
		Method:    httplog.HTTPMethod(c.Request.Method),
		Request:   reqLog.Request,
		Response: map[string]interface{}{
			"status_code": c.Writer.Status(),
			"headers":     c.Writer.Header(),
			"body":        responseBody,
		},
		StatusCode: c.Writer.Status(),
	}

	if err := config.Service.LogOutgoingRequest(c.Request.Context(), logEntry); err != nil {
		log.Printf("httplog: failed to log outgoing request: %v", err)
	}
}

// bodyLogWriter is a custom ResponseWriter that captures the response body
type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// Write captures the response body
func (w *bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// contains checks if a string is present in a slice
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if strings.EqualFold(s, item) {
			return true
		}
	}
	return false
}
