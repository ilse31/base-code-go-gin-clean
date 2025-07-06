package dto

// RegisterResponse represents the response after successful registration
type RegisterResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
