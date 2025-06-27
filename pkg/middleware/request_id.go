package middleware

import (
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RequestID adds a request ID to the context of each request
func RequestID() gin.HandlerFunc {
	return requestid.New(
		requestid.WithGenerator(func() string {
			return uuid.New().String()
		}),
	)
}
