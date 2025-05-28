# Middleware Implementation in Clean Architecture

This document describes the complete middleware system implemented for the Parallax API following clean architecture principles.

## Overview

The middleware system has been successfully integrated into the clean architecture with the following components:

### Implemented Middlewares

1. **CORS Middleware** (`cors.go`)
   - Environment-based configuration
   - Support for wildcard subdomains
   - Preflight request handling
   - Credential support

2. **CSRF Middleware** (`csrf.go`)
   - Double-submit cookie pattern
   - Constant-time comparison for security
   - Configurable token lookup (header/form)
   - Automatic token generation and validation

3. **Authentication Middleware** (`auth.go`)
   - Session-based authentication using cookies
   - In-memory session store with cleanup
   - Role-based authorization
   - Login/logout functionality

4. **Rate Limiting Middleware** (`ratelimit.go`)
   - Token bucket algorithm
   - Per-client rate limiting
   - Configurable key extraction strategies
   - Rate limit headers in responses

5. **Middleware Manager** (`manager.go`)
   - Central management of all middlewares
   - Global and route-specific middleware setup
   - Security headers
   - Cleanup routines

## Architecture Structure

```
internal/adapters/http/middleware/
├── manager.go      # Central middleware manager
├── auth.go         # Session-based authentication
├── cors.go         # CORS protection
├── csrf.go         # CSRF protection
└── ratelimit.go    # Rate limiting
```

## Route Protection Levels

### Public Routes (No Authentication)
- Health check endpoints
- Authentication endpoints (login, logout, csrf)
- Public blog posts (read-only)
- Public counter views (read-only)

### Protected Routes (Authentication + CSRF Required)
- Blog post creation, editing, deletion
- Counter operations (create, update, delete)
- User profile operations

### Admin Routes (Admin Role Required)
- User management endpoints
- System administration endpoints

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

## API Endpoints

### Authentication Endpoints
```
POST /api/v1/auth/login    # User login
POST /api/v1/auth/logout   # User logout
GET  /api/v1/auth/csrf     # Get CSRF token
GET  /api/v1/auth/me       # Get authentication status
```

### Public Blog Endpoints
```
GET /api/v1/blog/posts               # List all posts
GET /api/v1/blog/posts/published     # Published posts only
GET /api/v1/blog/posts/search        # Search posts
GET /api/v1/blog/posts/:id           # Get post by ID
GET /api/v1/blog/posts/slug/:slug    # Get post by slug
```

### Protected Blog Endpoints (Auth + CSRF)
```
POST   /api/v1/blog/posts                # Create post
PUT    /api/v1/blog/posts/:id            # Update post
DELETE /api/v1/blog/posts/:id            # Delete post
POST   /api/v1/blog/posts/:id/publish    # Publish post
POST   /api/v1/blog/posts/:id/unpublish  # Unpublish post
```

### Protected Counter Endpoints (Auth + CSRF)
```
POST /api/v1/counter/              # Create counter
DELETE /api/v1/counter/:id         # Delete counter
PUT /api/v1/counter/:id/increment  # Increment counter
PUT /api/v1/counter/:id/decrement  # Decrement counter
PUT /api/v1/counter/:id/value      # Set counter value
PUT /api/v1/counter/:id/reset      # Reset counter
```

### Admin Endpoints (Admin Role Required)
```
GET /api/v1/admin/users      # List users
DELETE /api/v1/admin/users/:id  # Delete user
```

## Security Features

### Automatic Security Headers
- `X-Content-Type-Options: nosniff`
- `X-Frame-Options: DENY`
- `X-XSS-Protection: 1; mode=block`
- `Referrer-Policy: strict-origin-when-cross-origin`
- `Content-Security-Policy: default-src 'self'`

### Rate Limiting Headers
- `X-RateLimit-Limit`: Maximum requests allowed
- `X-RateLimit-Remaining`: Remaining requests in window
- `X-RateLimit-Reset`: Reset timestamp

## Frontend Integration

### Getting CSRF Token
```javascript
const response = await fetch('/api/v1/auth/csrf', {
  credentials: 'include'
});
const data = await response.json();
const csrfToken = data.data.csrf_token;
```

### Making Authenticated Requests
```javascript
// Login
await fetch('/api/v1/auth/login', {
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
await fetch('/api/v1/blog/posts', {
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

## Testing

A comprehensive test script is provided at `test_middleware.sh` that tests:

- CORS preflight and actual requests
- CSRF token generation and validation
- Authentication flow (login/logout)
- Rate limiting functionality
- Security headers presence
- Protected route access control
- Admin role authorization

Run the tests:
```bash
chmod +x test_middleware.sh
./test_middleware.sh
```

## Production Considerations

### HTTPS and Secure Cookies
Set these environment variables for production:
```bash
APP_ENV=production
CSRF_COOKIE_SECURE=true
```

### Session Storage
The current implementation uses in-memory session storage. For production, consider:
- Redis-based session store
- Database-backed sessions
- Distributed session management

### CORS Configuration
Restrict CORS origins in production:
```bash
CORS_ALLOW_ORIGINS=https://yourdomain.com,https://app.yourdomain.com
```

### Rate Limiting
For multiple server instances, implement distributed rate limiting:
- Redis-based rate limiting
- External rate limiting service
- Load balancer rate limiting

## Key Benefits

1. **Clean Architecture Compliance**: Middlewares are properly organized in the adapters layer
2. **Security by Design**: Multiple layers of protection (CORS, CSRF, Auth, Rate Limiting)
3. **Environment-Driven Configuration**: Easy to configure for different environments
4. **Session-Based Authentication**: Secure cookie-based sessions instead of stateless tokens
5. **Comprehensive Protection**: Covers all major web security concerns
6. **Easy Testing**: Comprehensive test suite for validation
7. **Production Ready**: Includes security headers and production considerations

## Migration Notes

- Old middleware files in `internal/middleware/` can be removed
- All middleware setup is now handled by the MiddlewareManager
- Routes are automatically protected based on their group configuration
- Session-based authentication replaces any JWT-based authentication

This implementation provides enterprise-level security while maintaining clean architecture principles and ease of use.