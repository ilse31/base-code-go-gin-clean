package health

// HealthResponse represents the health check response structure
type HealthResponse struct {
	Version  string         `json:"version" example:"1.0.0"`
	Status   string         `json:"status" example:"ok"`
	Database DatabaseStatus `json:"database"`
}

// DatabaseStatus represents database health status
type DatabaseStatus struct {
	Status  string `json:"status" example:"ok"`
	Message string `json:"message,omitempty"`
}
