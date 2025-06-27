package health

// HealthResponse represents the health check response structure
type HealthResponse struct {
	Version string `json:"version" example:"1.0.0"`
	Status  string `json:"status" example:"ok"`
}
