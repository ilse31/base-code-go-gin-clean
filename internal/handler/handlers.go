package handler

import (
	emailDomain "base-code-go-gin-clean/internal/domain/email"
	"base-code-go-gin-clean/internal/handler/auth"
	email "base-code-go-gin-clean/internal/handler/email"
	"base-code-go-gin-clean/internal/handler/health"
	"base-code-go-gin-clean/internal/handler/roles"
	"base-code-go-gin-clean/internal/handler/user"
	"base-code-go-gin-clean/internal/service"
)

// UserHandler is an alias for user.UserHandler
type UserHandler = user.UserHandler

// NewUserHandler creates a new UserHandler
func NewUserHandler(userService service.UserService) *UserHandler {
	return user.NewUserHandler(userService)
}

// RolesHandler is an alias for roles.RolesHandler
type RolesHandler = roles.RolesHandler

// NewRolesHandler creates a new RolesHandler
func NewRolesHandler() *RolesHandler {
	return roles.NewRolesHandler()
}

// HealthHandler is an alias for health.HealthHandler
type HealthHandler = health.HealthHandler

// NewHealthHandler creates a new HealthHandler
func NewHealthHandler() *HealthHandler {
	return health.NewHealthHandler()
}

// EmailHandler is an alias for email.EmailHandler
type EmailHandler = email.EmailHandler

// NewEmailHandler creates a new EmailHandler
func NewEmailHandler(emailService emailDomain.EmailService) *EmailHandler {
	return email.NewEmailHandler(emailService)
}

// AuthHandler is an alias for auth.AuthHandler
type AuthHandler = auth.AuthHandler

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return auth.NewAuthHandler(authService)
}
