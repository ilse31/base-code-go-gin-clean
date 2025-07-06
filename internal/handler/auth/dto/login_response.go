package dto

// TokenResponse represents the token information in the response
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
	ExpiresIn    int64  `json:"expires_in"`
}

// LoginResponse represents the response after successful login
type LoginResponse struct {
	User  *UserInfo     `json:"user"`
	Token TokenResponse `json:"token"`
}

// UserInfo represents the user information in the login response
type UserInfo struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
