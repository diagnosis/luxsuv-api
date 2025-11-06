# Middleware Guide

## Overview

The `internal/middleware` package provides authentication and authorization middleware for protecting API endpoints.

## Available Middlewares

### 1. RequireJWT

Validates JWT access tokens and extracts user information into the request context.

**Features:**
- Validates Authorization header format (`Bearer <token>`)
- Parses and validates JWT tokens
- Verifies token signature, expiration, issuer, and audience
- Extracts user ID and role into context
- Logs authentication attempts

**Usage:**
```go
import customMiddleware "github.com/diagnosis/luxsuv-api-v2/internal/middleware"

// Protect a route group
api.Group(func(protected chi.Router) {
    protected.Use(customMiddleware.RequireJWT(app.Signer))

    // All routes in this group require valid JWT
    protected.Get("/profile", handler.GetProfile)
    protected.Put("/profile", handler.UpdateProfile)
})
```

**Request Example:**
```bash
curl -X GET http://localhost:8080/api/v1/profile \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

**Responses:**

Success: Continues to next handler with user context set

Error (401):
```json
{
  "error": {
    "code": "UNAUTHORIZED",
    "message": "Missing authorization header",
    "correlation_id": "abc-123",
    "timestamp": "2024-01-15T10:30:45Z"
  }
}
```

### 2. RequireRole

Checks if the authenticated user has one of the allowed roles.

**Features:**
- Validates user has required role(s)
- Supports multiple allowed roles
- Must be used AFTER RequireJWT
- Logs authorization attempts

**Usage:**
```go
// Admin-only routes
api.Group(func(adminOnly chi.Router) {
    adminOnly.Use(customMiddleware.RequireJWT(app.Signer))
    adminOnly.Use(customMiddleware.RequireRole("admin"))

    adminOnly.Get("/users", handler.ListAllUsers)
    adminOnly.Delete("/users/:id", handler.DeleteUser)
})

// Driver OR admin routes
api.Group(func(drivers chi.Router) {
    drivers.Use(customMiddleware.RequireJWT(app.Signer))
    drivers.Use(customMiddleware.RequireRole("driver", "admin"))

    drivers.Get("/rides", handler.GetRides)
})
```

**Responses:**

Success: Continues to next handler

Error (403):
```json
{
  "error": {
    "code": "FORBIDDEN",
    "message": "Insufficient permissions",
    "correlation_id": "abc-123",
    "timestamp": "2024-01-15T10:30:45Z"
  }
}
```

## Context Functions

### GetUserID

Retrieves the authenticated user's ID from the request context.

```go
import "github.com/diagnosis/luxsuv-api-v2/internal/middleware"

func MyHandler(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    userID, ok := middleware.GetUserID(ctx)
    if !ok {
        // User not authenticated
        return
    }

    // Use userID
    log.Printf("User ID: %s", userID)
}
```

### GetUserRole

Retrieves the authenticated user's role from the request context.

```go
import "github.com/diagnosis/luxsuv-api-v2/internal/middleware"

func MyHandler(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    role, ok := middleware.GetUserRole(ctx)
    if !ok {
        // User not authenticated
        return
    }

    // Use role
    log.Printf("User Role: %s", role)
}
```

## Complete Route Examples

### Public Routes
```go
r.Route("/api/v1", func(api chi.Router) {
    // No authentication required
    api.Post("/auth/login", app.UserHandler.HandleLogin)
    api.Post("/auth/register", app.UserHandler.HandleRegister)
})
```

### Protected Routes (Any Authenticated User)
```go
r.Route("/api/v1", func(api chi.Router) {
    api.Group(func(protected chi.Router) {
        protected.Use(customMiddleware.RequireJWT(app.Signer))

        // Any authenticated user can access these
        protected.Get("/profile", handler.GetProfile)
        protected.Put("/profile", handler.UpdateProfile)
        protected.Get("/rides/history", handler.GetRideHistory)
    })
})
```

### Admin-Only Routes
```go
r.Route("/api/v1", func(api chi.Router) {
    api.Group(func(adminOnly chi.Router) {
        adminOnly.Use(customMiddleware.RequireJWT(app.Signer))
        adminOnly.Use(customMiddleware.RequireRole("admin"))

        // Only admins can access these
        adminOnly.Get("/admin/users", handler.ListAllUsers)
        adminOnly.Delete("/admin/users/:id", handler.DeleteUser)
        adminOnly.Get("/admin/stats", handler.GetStats)
    })
})
```

### Driver-Only Routes
```go
r.Route("/api/v1", func(api chi.Router) {
    api.Group(func(driverOnly chi.Router) {
        driverOnly.Use(customMiddleware.RequireJWT(app.Signer))
        driverOnly.Use(customMiddleware.RequireRole("driver"))

        // Only drivers can access these
        driverOnly.Get("/driver/rides", handler.GetAvailableRides)
        driverOnly.Post("/driver/rides/:id/accept", handler.AcceptRide)
        driverOnly.Put("/driver/status", handler.UpdateDriverStatus)
    })
})
```

### Multi-Role Routes
```go
r.Route("/api/v1", func(api chi.Router) {
    api.Group(func(driversAndAdmins chi.Router) {
        driversAndAdmins.Use(customMiddleware.RequireJWT(app.Signer))
        driversAndAdmins.Use(customMiddleware.RequireRole("driver", "admin"))

        // Drivers and admins can access these
        driversAndAdmins.Get("/rides/:id", handler.GetRideDetails)
        driversAndAdmins.Get("/vehicles", handler.ListVehicles)
    })
})
```

## Handler Implementation Example

```go
package api

import (
    "net/http"

    "github.com/diagnosis/luxsuv-api-v2/internal/helper"
    "github.com/diagnosis/luxsuv-api-v2/internal/logger"
    "github.com/diagnosis/luxsuv-api-v2/internal/middleware"
)

func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()

    // Get authenticated user ID from context
    userID, ok := middleware.GetUserID(ctx)
    if !ok {
        // This should never happen if RequireJWT is applied
        logger.Error(ctx, "user_id not found in context")
        helper.RespondError(w, r, apperror.Unauthorized("Authentication required"))
        return
    }

    logger.Debug(ctx, "fetching user profile", "user_id", userID)

    // Fetch user from database
    user, err := h.UserStore.GetByID(ctx, userID)
    if err != nil {
        helper.RespondError(w, r, err)
        return
    }

    helper.RespondJSON(w, r, http.StatusOK, user)
}
```

## Logging

All middleware operations are logged:

**Authentication Success:**
```
level=DEBUG msg="JWT authenticated" user_id=550e8400-e29b-41d4-a716-446655440000 role=driver correlation_id=abc-123
```

**Authentication Failure:**
```
level=WARN msg="invalid or expired token" error="token is expired" correlation_id=abc-123
```

**Authorization Success:**
```
level=DEBUG msg="role authorization passed" user_id=550e8400... role=admin correlation_id=abc-123
```

**Authorization Failure:**
```
level=WARN msg="insufficient permissions" user_id=550e8400... role=rider required_roles=[admin] correlation_id=abc-123
```

## Error Codes

| Error Code | HTTP Status | Description |
|------------|-------------|-------------|
| `UNAUTHORIZED` | 401 | Missing, invalid, or expired JWT token |
| `FORBIDDEN` | 403 | Valid token but insufficient permissions |

## Security Notes

1. **Always use HTTPS in production** to protect tokens in transit
2. **Token expiration:** Access tokens expire after 15 minutes by default
3. **RequireRole must be used AFTER RequireJWT** - it depends on JWT middleware setting the context
4. **Refresh tokens:** Use the refresh token endpoint to get new access tokens
5. **Audit logging:** All authentication and authorization attempts are logged for security auditing

## Testing

### Test with curl

```bash
# 1. Login to get token
TOKEN=$(curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}' \
  | jq -r '.data.access_token')

# 2. Use token to access protected endpoint
curl -X GET http://localhost:8080/api/v1/profile \
  -H "Authorization: Bearer $TOKEN"

# 3. Try admin endpoint (will fail if not admin)
curl -X GET http://localhost:8080/api/v1/admin/users \
  -H "Authorization: Bearer $TOKEN"
```

## Summary

- ✅ **RequireJWT** - Validates JWT tokens and sets user context
- ✅ **RequireRole** - Checks user has required role(s)
- ✅ **GetUserID** - Retrieves user ID from context
- ✅ **GetUserRole** - Retrieves user role from context
- ✅ **Comprehensive logging** - All auth attempts logged
- ✅ **Secure by default** - Fails closed on errors
- ✅ **Flexible** - Support multiple roles, easy to use

Use these middlewares to secure your API endpoints with JWT authentication and role-based authorization!
