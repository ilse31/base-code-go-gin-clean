package auth

import (
	"base-code-go-gin-clean/internal/handler/auth/dto"
	httpPkg "base-code-go-gin-clean/internal/pkg/http"
	"base-code-go-gin-clean/internal/service"
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	accessTokenCookieName  = "access_token"
	refreshTokenCookieName = "refresh_token"
	maxAgeAccessToken      = 15 * 60       // 15 minutes in seconds
	maxAgeRefreshToken     = 7 * 24 * 3600 // 7 days in seconds
)

// Context keys for storing values in the request context
const (
	userIDKey = "userID"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register handles user registration
// @Summary Register a new user
// @Description Register a new user with name, email, and password. Password must be at least 8 characters long.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Registration details"
// @Success 201 {object} handler.SuccessResponse{data=dto.RegisterResponse} "User registered successfully"
// @Failure 400 {object} handler.ErrorResponse "Bad Request: Invalid input format or validation failed"
// @Failure 409 {object} handler.ErrorResponse "Conflict: Email already exists"
// @Failure 500 {object} handler.ErrorResponse "Internal Server Error: Failed to process registration"
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpPkg.BadRequest(c, "Invalid request body: "+err.Error(), nil)
		return
	}

	user, err := h.authService.Register(c.Request.Context(), req.Name, req.Email, req.Password)
	if err != nil {
		if err.Error() == "user with this email already exists" {
			httpPkg.ErrorResponse(c, 409, "Email already exists", nil)
		} else {
			httpPkg.InternalServerError(c, "Failed to register user")
		}
		return
	}

	// Convert user to response DTO
	response := &dto.RegisterResponse{
		ID:    user.ID.String(),
		Name:  user.Name,
		Email: user.Email,
	}

	httpPkg.Created(c, response)
}

// Login handles user login
// @Summary Authenticate a user
// @Description Authenticate user with email and password. Returns user details on success.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login credentials"
// @Success 200 {object} handler.SuccessResponse{data=dto.LoginResponse} "Login successful"
// @Failure 400 {object} handler.ErrorResponse "Bad Request: Invalid input format"
// @Failure 401 {object} handler.ErrorResponse "Unauthorized: Invalid email or password"
// @Failure 500 {object} handler.ErrorResponse "Internal Server Error: Failed to process login"
// @Router /auth/login [post]
// setAuthCookies sets the authentication cookies in the response
func setAuthCookies(c *gin.Context, accessToken, refreshToken string) {
	host := c.Request.Host
	// Remove port if present
	if host != "" {
		// First try net.SplitHostPort
		if h, _, err := net.SplitHostPort(host); err == nil {
			host = h
		} else if i := strings.Index(host, ":"); i != -1 {
			// Fallback to simple string split if net.SplitHostPort fails
			host = host[:i]
		}
	}

	// Set access token cookie (HTTP-only, Secure, SameSite=Strict)
	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie(
		accessTokenCookieName,
		accessToken,
		maxAgeAccessToken,
		"/",
		host,
		false, // Secure - set to true in production with HTTPS
		true,  // httpOnly
	)

	// Set refresh token cookie (HTTP-only, Secure, SameSite=Strict, Path=/auth/refresh)
	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie(
		refreshTokenCookieName,
		refreshToken,
		maxAgeRefreshToken,
		"/auth/refresh",
		host,
		false, // Secure - set to true in production with HTTPS
		true,  // httpOnly
	)
}

// clearAuthCookies clears the authentication cookies
func clearAuthCookies(c *gin.Context) {
	// Clear access token cookie
	c.SetCookie(
		accessTokenCookieName,
		"",
		-1, // MaxAge<0 means delete cookie
		"/",
		"",
		false,
		true,
	)

	// Clear refresh token cookie
	c.SetCookie(
		refreshTokenCookieName,
		"",
		-1, // MaxAge<0 means delete cookie
		"/auth/refresh",
		"",
		false,
		true,
	)
}

// Login handles user login
// @Summary Authenticate a user
// @Description Authenticate user with email and password. Returns user details and sets HTTP-only cookies with access and refresh tokens.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login credentials"
// @Success 200 {object} handler.SuccessResponse{data=dto.LoginResponse} "Login successful"
// @Failure 400 {object} handler.ErrorResponse "Bad Request: Invalid input format"
// @Failure 401 {object} handler.ErrorResponse "Unauthorized: Invalid email or password"
// @Failure 500 {object} handler.ErrorResponse "Internal Server Error: Failed to process login"
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpPkg.BadRequest(c, "Invalid request body: "+err.Error(), nil)
		return
	}

	loginResponse, err := h.authService.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		httpPkg.Unauthorized(c, "Invalid email or password")
		return
	}

	fmt.Println(loginResponse)

	// Set HTTP-only cookies
	setAuthCookies(c, loginResponse.Token.AccessToken, loginResponse.Token.RefreshToken)

	// Prepare response (don't include tokens in the response body when using cookies)
	response := &dto.LoginResponse{
		User: &dto.UserInfo{
			ID:    loginResponse.User.ID.String(),
			Name:  loginResponse.User.Name,
			Email: loginResponse.User.Email,
		},
		Token: dto.TokenResponse{
			AccessToken:  "", // Not including in response when using cookies
			RefreshToken: "", // Not including in response when using cookies
			ExpiresIn:    loginResponse.Token.ExpiresIn,
		},
	}

	httpPkg.Success(c, response)
}

// RefreshToken handles access token refresh using a refresh token
// @Summary Refresh access token
// @Description Refresh access token using a refresh token from HTTP-only cookie
// @Tags Authentication
// @Accept json
// @Produce json
// @Success 200 {object} handler.SuccessResponse{data=dto.TokenResponse} "Token refreshed successfully"
// @Failure 401 {object} handler.ErrorResponse "Unauthorized: Invalid or expired refresh token"
// @Failure 500 {object} handler.ErrorResponse "Internal Server Error: Failed to refresh token"
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	// Get refresh token from cookie
	refreshToken, err := c.Cookie(refreshTokenCookieName)
	if err != nil {
		httpPkg.Unauthorized(c, "Refresh token is required")
		return
	}

	// Call service to refresh token
	tokenResponse, err := h.authService.RefreshToken(c.Request.Context(), refreshToken)
	if err != nil {
		httpPkg.Unauthorized(c, "Invalid or expired refresh token")
		return
	}

	// Set new cookies
	setAuthCookies(c, tokenResponse.AccessToken, tokenResponse.RefreshToken)

	// Don't include tokens in the response body
	tokenResponse.AccessToken = ""
	tokenResponse.RefreshToken = ""

	httpPkg.Success(c, tokenResponse)
}

// Logout handles user logout
// @Summary Logout user
// @Description Logout user by clearing authentication cookies
// @Tags Authentication
// @Accept json
// @Produce json
// @Success 200 {object} handler.SuccessResponse{} "Logout successful"
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get(userIDKey)
	if exists && userID != nil {
		// Invalidate the refresh token
		_ = h.authService.Logout(c.Request.Context(), userID.(string))
	}

	// Clear auth cookies
	clearAuthCookies(c)

	httpPkg.Success(c, nil)
}
