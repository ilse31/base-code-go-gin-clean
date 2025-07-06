# Authentication System

This document describes the token-based authentication system implemented in the application.

## Overview

The authentication system uses JWT (JSON Web Tokens) for stateless authentication with access and refresh tokens. The system is designed to be secure, scalable, and easy to integrate with the existing codebase.

## Components

### 1. Token Service (`internal/pkg/token`)

The token service handles the generation and validation of JWT tokens.

**Key Features:**

- Generates access tokens (short-lived JWTs)
- Generates refresh tokens (long-lived random strings)
- Validates access tokens
- Configurable token expiration times
- Secure token signing with HMAC-SHA256

### 2. Authentication Service (`internal/service`)

The authentication service implements the core authentication logic.

**Key Features:**

- User registration and login
- Token generation and refresh
- Secure password hashing and verification
- Session management with refresh tokens

### 3. Authentication Middleware (`internal/middleware`)

Middleware for protecting routes and validating tokens.

**Key Features:**

- Validates access tokens in protected routes
- Extracts user information from tokens
- Handles token expiration and renewal

## Configuration

The authentication system can be configured using environment variables:

```env
# JWT Configuration
ACCESS_TOKEN_SECRET=your-access-token-secret-key-32-chars-long
REFRESH_TOKEN_SECRET=your-refresh-token-secret-key-32-chars-long
ACCESS_TOKEN_EXPIRY=15        # in minutes
REFRESH_TOKEN_EXPIRY=10080    # in minutes (7 days)
```

## API Endpoints

### Authentication Endpoints

#### `POST /api/v1/auth/register`

Register a new user.

**Request:**

```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "securePassword123"
}
```

**Response:**

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "John Doe",
  "email": "john@example.com"
}
```

#### `POST /api/v1/auth/login`

Authenticate a user and get access and refresh tokens.

**Request:**

```json
{
  "email": "john@example.com",
  "password": "securePassword123"
}
```

**Response:**
Sets HTTP-only cookies:

- `access_token`: JWT access token
- `refresh_token`: Refresh token

```json
{
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "John Doe",
    "email": "john@example.com"
  },
  "token": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "refresh_token_here",
    "expires_in": 900
  }
}
```

#### `POST /api/v1/auth/refresh`

Refresh an expired access token using a refresh token.

**Request:**

- Requires `refresh_token` cookie

**Response:**
Sets new HTTP-only cookies:

- `access_token`: New JWT access token
- `refresh_token`: New refresh token (optional, if rotation is enabled)

#### `POST /api/v1/auth/logout`

Invalidate the current session.

**Response:**

- Clears authentication cookies
- Returns 200 OK on success

## Protecting Routes

To protect a route, use the `AuthMiddleware`:

```go
// In your route definitions
router.GET("/api/v1/protected", middleware.AuthMiddleware(tokenService), protectedHandler)
```

## Security Considerations

1. **Token Storage**:

   - Access tokens are stored in memory on the client side
   - Refresh tokens are stored in HTTP-only, secure, SameSite=Strict cookies

2. **Token Expiration**:

   - Access tokens expire after 15 minutes (configurable)
   - Refresh tokens expire after 7 days (configurable)

3. **Password Security**:

   - Passwords are hashed using bcrypt before storage
   - Minimum password length and complexity should be enforced

4. **HTTPS**:
   - Always use HTTPS in production to protect tokens in transit
   - The `Secure` flag is set on cookies when in production

## Testing

Authentication can be tested using the test suite in `internal/service/auth_service_test.go`. The test suite includes tests for:

- User registration
- Login with valid/invalid credentials
- Token refresh
- Session invalidation

## Future Improvements

1. Implement refresh token rotation
2. Add rate limiting for authentication endpoints
3. Add support for OAuth2 providers
4. Implement password reset functionality
5. Add account lockout after failed login attempts
