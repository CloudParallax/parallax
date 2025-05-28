# Middleware Documentation

This document describes the middleware system implemented in the Parallax API following clean architecture principles.

## Overview

The middleware system provides:
- **CORS** - Cross-Origin Resource Sharing protection
- **CSRF** - Cross-Site Request Forgery protection
- **Authentication** - Session-based authentication
- **Rate Limiting** - Request rate limiting
- **Security Headers** - Various security headers

## Architecture

```
internal/adapters/http/middleware/
├── manager.go      # Middleware manager
├── auth.go         # Session-based authentication
├── cors.go         # CORS protection
├── csrf.go         # CSRF protection
└── ratelimit.go    # Rate limiting
```

## Configuration

Environment variables for middleware configuration:

```bash
# Session Configuration
SESSION_COOKIE_NAME=session_id
SESSION_MAX_AGE=24h

# CORS Configuration
CORS_ALLOW_ORIGINS=http://localhost:3000,http://localhost:8080
CORS_ALLOW_METHODS=GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS
CORS_ALLOW_HEADERS=Origin,Content-Type,Accept,Authorization,X-Requested-With,X-CSRF-Token
CORS_EXPOSE_HEADERS=Content-Length,X-CSRF-Token
CORS_ALLOW_CREDENTIALS=true
CORS_MAX_AGE=86400

# CSRF Protection
CSRF_COOKIE_NAME=csrf_token
CSRF_HEADER_NAME=X-CSRF-Token
CSRF_COOKIE_SECURE=false
CSRF_COOKIE_HTTPONLY=true
CSRF_EXPIRATION=24h

# Rate Limiting
RATE_LIMIT_MAX_REQUESTS=100
RATE_LIMIT_WINDOW=1m
```

## Middleware Types

### 1. CORS Middleware

Handles Cross-Origin Resource Sharing with configurable origins, methods, and headers.

**Features:**
- Environment-based configuration
- Wildcard subdomain support (*.example.com)
- Credential support
- Preflight request handling

**Usage:**
```go
corsMiddleware := middleware.NewCORSMiddleware()
app.Use(corsMiddleware.Handler())
```

### 2. CSRF Middleware

Implements double-submit cookie pattern for CSRF protection.

**Features:**
- Double-submit cookie pattern
- Configurable token lookup (header/form)
- Constant-time comparison
- Automatic token generation

**Usage:**
```go
csrfMiddleware := middleware.NewCSRFMiddleware()
app.Use(csrfMiddleware.Handler())
```

### 3. Authentication Middleware

Session-based authentication with in-memory session store.

**Features:**
- Session management
- Cookie-based authentication
- Role-based authorization
- Session cleanup

**Usage:**
```go
authMiddleware := middleware.NewAuthMiddleware("session_id", 24*time.Hour)

// Require authentication
app.Use(authMiddleware.RequireAuth())

// Optional authentication
app.Use(authMiddleware.OptionalAuth())

// Require specific role
app.Use(authMiddleware.RequireRole("admin"))
```

### 4. Rate Limiting Middleware

Token bucket algorithm for rate limiting requests.

**Features:**
- Per-client rate limiting
- Configurable time windows
- Multiple key extraction strategies
- Rate limit headers

**Usage:**
```go
rateLimitMiddleware := middleware.NewRateLimitMiddleware(middleware.RateLimitConfig{
    MaxRequests: 100,
    Window:      time.Minute,
    KeyFunc:     middleware.KeyFuncByIP(),
})
app.Use(rateLimitMiddleware.Handler())
```

## Route Protection Levels

### Public Routes (No Authentication)
- Health check
- Authentication endpoints
- Public blog posts (read-only)

### Protected Routes (Authentication Required)
- Blog post creation/editing
- Counter operations
- User profile operations

### Admin Routes (Admin Role Required)
- User management
- System administration

## API Endpoints

### Authentication
```
POST /api/v1/auth/login    # Login user
POST /api/v1/auth/logout   # Logout user
GET  /api/v1/auth/csrf     # Get CSRF token
GET  /api/v1/auth/me       # Get auth status
```

### Public Blog Routes
```
GET /api/v1/blog/posts               # List posts
GET /api/v1/blog/posts/published     # Published posts
GET /api/v1/blog/posts/search        # Search posts
GET /api/v1/blog/posts/:id           # Get post by ID
GET /api/v1/blog/posts/slug/:slug    # Get post by slug
```

### Protected Blog Routes (Auth + CSRF Required)
```
POST   /api/v1/blog/posts                # Create post
PUT    /api/v1/blog/posts/:id            # Update post
DELETE /api/v1/blog/posts/:id            # Delete post
POST   /api/v1/blog/posts/:id/publish    # Publish post
POST   /api/v1/blog/posts/:id/unpublish  # Unpublish post
```

## Frontend Integration

### Getting CSRF Token
```javascript
// Get CSRF token
const response = await fetch('/api/v1/auth/csrf', {
    credentials: 'include'
});
const data = await response.json();
const csrfToken = data.data.csrf_token;
```

### Making Authenticated Requests
```javascript
// Login
const loginResponse = await fetch('/api/v1/auth/login', {
    method: 'POST',
    credentials: 'include',
    headers: {
        'Content-Type': 'application/json',
        'X-CSRF-Token': csrfToken
    },
    body: JSON.stringify({
        username: 'user',
        password: 'pass'
    })
});

// Protected request
const response = await fetch('/api/v1/blog/posts', {
    method: 'POST',
    credentials: 'include',
    headers: {
        'Content-Type': 'application/json',
        'X-CSRF-Token': csrfToken
    },
    body: JSON.stringify({
        title: 'New Post',
        content: 'Post content'
    })
});
```

## Security Features

### Security Headers
- `X-Content-Type-Options: nosniff`
- `X-Frame-Options: DENY`
- `X-XSS-Protection: 1; mode=block`
- `Referrer-Policy: strict-origin-when-cross-origin`
- `Content-Security-Policy: default-src 'self'`

### Rate Limiting Headers
- `X-RateLimit-Limit`: Maximum requests allowed
- `X-RateLimit-Remaining`: Remaining requests
- `X-RateLimit-Reset`: Reset timestamp

## Testing

Use the provided test script to verify middleware functionality:

```bash
./test_middleware.sh
```

The script tests:
- CORS preflight and actual requests
- CSRF token generation and validation
- Authentication flow
- Rate limiting
- Security headers
- Protected route access

## Production Considerations

### HTTPS
Enable secure cookies in production:
```bash
CSRF_COOKIE_SECURE=true
APP_ENV=production
```

### Session Store
Replace in-memory session store with persistent storage:
- Redis
- Database
- Distributed cache

### Rate Limiting
Consider distributed rate limiting for multiple server instances:
- Redis-based rate limiting
- External rate limiting service

### CORS Origins
Restrict CORS origins in production:
```bash
CORS_ALLOW_ORIGINS=https://yourdomain.com,https://app.yourdomain.com
```

## Troubleshooting

### CORS Issues
- Check `CORS_ALLOW_ORIGINS` configuration
- Verify credentials support if using cookies
- Check preflight request handling

### CSRF Issues
- Ensure token is included in requests
- Check cookie settings (HttpOnly, Secure, SameSite)
- Verify token generation and validation

### Authentication Issues
- Check session cookie configuration
- Verify session expiration settings
- Check cleanup routines

### Rate Limiting Issues
- Adjust limits based on usage patterns
- Consider different key functions
- Monitor rate limit headers