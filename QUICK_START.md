# Quick Start Guide

Get the Parallax API running in minutes!

## Prerequisites

- Go 1.24+ installed
- Git

## 1. Clone and Setup

```bash
git clone <repository-url>
cd skopsgo
```

## 2. Install Dependencies

```bash
go mod tidy
```

## 3. Configure Environment

```bash
cp .env.example .env
```

Edit `.env` if needed (default values work for local development).

## 4. Start the Server

```bash
go run cmd/parallax/parallax.go
```

You should see:
```
🚀 Starting Parallax API server at :8080
```

## 5. Test the API

### Basic Hello World
```bash
curl http://localhost:8080/api/v1/hello
```

Expected response:
```json
{
  "success": true,
  "message": "Hello, World!",
  "api": "Parallax API",
  "version": "1.0.0",
  "middleware": {
    "cors": "enabled",
    "csrf": "ready", 
    "auth": "available",
    "rate_limit": "active"
  }
}
```

### Hello with Name
```bash
curl "http://localhost:8080/api/v1/hello?name=Developer"
```

### Health Check
```bash
curl http://localhost:8080/api/v1/health
```

## 6. Run Tests

```bash
# Test hello world endpoint
./test_hello.sh

# Test all middleware functionality
./test_middleware.sh
```

## Available Endpoints

### Public (No Auth Required)
- `GET /api/v1/hello` - Hello world
- `GET /api/v1/health` - Health check
- `GET /api/v1/blog/posts` - List blog posts
- `GET /api/v1/counter` - List counters

### Authentication
- `GET /api/v1/auth/csrf` - Get CSRF token
- `POST /api/v1/auth/login` - Login
- `POST /api/v1/auth/logout` - Logout
- `GET /api/v1/auth/me` - Auth status

### Protected (Auth + CSRF Required)
- `POST /api/v1/blog/posts` - Create blog post
- `PUT /api/v1/blog/posts/:id` - Update blog post
- `DELETE /api/v1/blog/posts/:id` - Delete blog post
- `POST /api/v1/counter` - Create counter
- `PUT /api/v1/counter/:id/increment` - Increment counter

### Admin Only
- `GET /api/v1/admin/users` - List users
- `DELETE /api/v1/admin/users/:id` - Delete user

## Authentication Flow

```bash
# 1. Get CSRF token
curl -c cookies.txt http://localhost:8080/api/v1/auth/csrf

# 2. Login (replace YOUR_CSRF_TOKEN)
curl -b cookies.txt -c cookies.txt \
  -X POST \
  -H "Content-Type: application/json" \
  -H "X-CSRF-Token: YOUR_CSRF_TOKEN" \
  -d '{"username":"testuser","password":"testpass"}' \
  http://localhost:8080/api/v1/auth/login

# 3. Make protected request
curl -b cookies.txt \
  -X POST \
  -H "Content-Type: application/json" \
  -H "X-CSRF-Token: YOUR_CSRF_TOKEN" \
  -d '{"title":"My Post","content":"Post content"}' \
  http://localhost:8080/api/v1/blog/posts
```

## Frontend Integration

```javascript
// Get CSRF token
const csrfResponse = await fetch('/api/v1/auth/csrf', {
  credentials: 'include'
});
const csrfData = await csrfResponse.json();
const csrfToken = csrfData.data.csrf_token;

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

// Make authenticated request
await fetch('/api/v1/blog/posts', {
  method: 'POST',
  credentials: 'include',
  headers: {
    'Content-Type': 'application/json',
    'X-CSRF-Token': csrfToken
  },
  body: JSON.stringify({
    title: 'My Post',
    content: 'Content here'
  })
});
```

## Troubleshooting

### Server won't start
- Check if port 8080 is already in use
- Verify Go is installed: `go version`
- Check for syntax errors: `go build ./cmd/parallax/`

### CORS issues
- Check `CORS_ALLOW_ORIGINS` in `.env`
- Ensure `credentials: 'include'` in frontend requests

### Authentication issues
- Always include CSRF token in headers
- Use `credentials: 'include'` for cookie-based auth
- Check session expiration (24h default)

### Rate limiting
- Default: 100 requests per minute
- Check `X-RateLimit-*` headers in responses
- Adjust `RATE_LIMIT_MAX_REQUESTS` in `.env`

## Next Steps

- Read [MIDDLEWARE.md](docs/MIDDLEWARE.md) for security details
- Check [API.md](API.md) for complete endpoint documentation
- See [MIDDLEWARE_IMPLEMENTATION.md](docs/MIDDLEWARE_IMPLEMENTATION.md) for architecture details

## Development

```bash
# Watch for changes (if air is installed)
air

# Build binary
go build -o parallax ./cmd/parallax/

# Run tests
go test ./...
```

That's it! Your Parallax API is now running with full middleware protection. 🚀