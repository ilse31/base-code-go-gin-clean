package dto

import (
	"base-code-go-gin-clean/internal/domain/user"
)

// UserResponse represents the user response structure
type UserResponse = user.UserResponse

// NewUserResponse creates a new UserResponse from domain model
func NewUserResponse(userResp *user.UserResponse) *UserResponse {
	return userResp
}

// UserListResponse represents a list of users response
type UserListResponse struct {
	Users []*UserResponse `json:"users"`
	Total int64           `json:"total"`
}

// NewUserListResponse creates a new UserListResponse from domain models
func NewUserListResponse(users []*user.UserResponse, total int64) *UserListResponse {
	resp := &UserListResponse{
		Total: total,
	}

	for _, userResp := range users {
		resp.Users = append(resp.Users, userResp)
	}

	return resp
}

// ErrorResponse represents a standard error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// NewErrorResponse creates a new error response
func NewErrorResponse(err error) *ErrorResponse {
	return &ErrorResponse{
		Error: err.Error(),
	}
}
