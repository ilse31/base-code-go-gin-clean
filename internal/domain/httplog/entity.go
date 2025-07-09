package httplog

import (
	"time"
)

// HTTPMethod represents the HTTP method used in the request
type HTTPMethod string

const (
	GET     HTTPMethod = "GET"
	POST    HTTPMethod = "POST"
	PUT     HTTPMethod = "PUT"
	DELETE  HTTPMethod = "DELETE"
	PATCH   HTTPMethod = "PATCH"
	HEAD    HTTPMethod = "HEAD"
	OPTIONS HTTPMethod = "OPTIONS"
	TRACE   HTTPMethod = "TRACE"
)

// LogOutgoingRequest represents an outgoing HTTP request made by the application
type LogOutgoingRequest struct {
	ID         int64       `bun:",pk,autoincrement"` // Gunakan autoincrement
	TraceID    string      `pg:",notnull" json:"trace_id"`
	EventName  string      `pg:",notnull" json:"event_name"`
	Endpoint   string      `pg:",notnull" json:"endpoint"`
	Method     HTTPMethod  `pg:",notnull,type:http_method" json:"method"`
	Request    interface{} `pg:",type:jsonb,notnull" json:"request"`
	StatusCode int         `pg:",notnull" json:"status_code"`
	Response   interface{} `pg:",type:jsonb" json:"response,omitempty"`
	CreatedAt  time.Time   `pg:",notnull,default:now()" json:"created_at"`
	UpdatedAt  time.Time   `pg:",notnull,default:now()" json:"updated_at"`
}

// LogIncomingRequest represents an incoming HTTP request to the application
type LogIncomingRequest struct {
	ID        string      `bun:",pk,type:uuid,default:uuid_generate_v4()"`
	TraceID   string      `pg:",notnull" json:"trace_id"`
	EventName string      `pg:",notnull" json:"event_name"`
	Endpoint  string      `pg:",notnull" json:"endpoint"`
	Method    HTTPMethod  `pg:",notnull,type:http_method" json:"method"`
	Request   interface{} `pg:",type:jsonb,notnull" json:"request"`
	IPAddress string      `pg:"ip_address" json:"ip_address,omitempty"`
	UserAgent string      `pg:"user_agent" json:"user_agent,omitempty"`
	CreatedAt time.Time   `pg:",notnull,default:now()" json:"created_at"`
	UpdatedAt time.Time   `pg:",notnull,default:now()" json:"updated_at"`
}

// LogError represents an error that occurred during request processing
type LogError struct {
	ID         string      `bun:",pk,type:uuid,default:uuid_generate_v4()"`
	TraceID    string      `pg:",notnull" json:"trace_id"`
	StatusCode int         `pg:",notnull" json:"status_code"`
	Error      string      `pg:",notnull" json:"error"`
	RequestID  *string     `pg:",type:uuid" json:"request_id,omitempty"`
	Traces     interface{} `pg:",type:jsonb" json:"traces,omitempty"`
	CreatedAt  time.Time   `pg:",notnull,default:now()" json:"created_at"`
	UpdatedAt  time.Time   `pg:",notnull,default:now()" json:"updated_at"`
}
