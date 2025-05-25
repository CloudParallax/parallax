# CORS Configuration

This application includes Cross-Origin Resource Sharing (CORS) middleware to handle cross-origin requests securely.

## Overview

CORS is implemented using Fiber's built-in CORS middleware and is configured through environment variables for flexibility across different deployment environments.

## Configuration

The CORS middleware can be configured using the following environment variables:

### Environment Variables

| Variable | Description | Default Value | Example |
|----------|-------------|---------------|---------|
| `CORS_ALLOW_ORIGINS` | Comma-separated list of allowed origins | `*` | `https://example.com,https://app.example.com` |
| `CORS_ALLOW_METHODS` | Comma-separated list of allowed HTTP methods | `GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS` | `GET,POST,PUT,DELETE` |
| `CORS_ALLOW_HEADERS` | Comma-separated list of allowed headers | `Origin,Content-Type,Accept,Authorization,X-Requested-With` | `Origin,Content-Type,Authorization` |
| `CORS_EXPOSE_HEADERS` | Comma-separated list of headers to expose to the client | `Content-Length` | `Content-Length,X-Total-Count` |
| `CORS_ALLOW_CREDENTIALS` | Whether to allow credentials (cookies, auth headers) | `false` | `true` |
| `CORS_MAX_AGE` | Max age for preflight requests cache (seconds) | `86400` | `3600` |

### Example Configuration

For development (allowing all origins):
```env
CORS_ALLOW_ORIGINS=*
CORS_ALLOW_CREDENTIALS=false
```

For production (specific origins):
```env
CORS_ALLOW_ORIGINS=https://yourdomain.com,https://app.yourdomain.com
CORS_ALLOW_CREDENTIALS=true
CORS_MAX_AGE=3600
```

## Security Considerations

1. **Never use `*` for origins in production** when `CORS_ALLOW_CREDENTIALS=true`
2. **Be specific with allowed origins** in production environments
3. **Limit allowed methods** to only what your API actually uses
4. **Review allowed headers** regularly to ensure they're necessary

## Implementation Details

The CORS middleware is loaded in the middleware chain before other application middleware to ensure all requests are properly handled.

Location: `internal/middleware/corsMiddleware.go`

The middleware uses helper functions to parse environment variables with fallback defaults, ensuring the application works even without explicit CORS configuration.

## Testing CORS

You can test CORS configuration using curl:

```bash
# Test preflight request
curl -X OPTIONS \
  -H "Origin: https://example.com" \
  -H "Access-Control-Request-Method: POST" \
  -H "Access-Control-Request-Headers: Content-Type" \
  http://localhost:8080/api/endpoint

# Test actual request
curl -X POST \
  -H "Origin: https://example.com" \
  -H "Content-Type: application/json" \
  http://localhost:8080/api/endpoint
```

The response should include appropriate `Access-Control-*` headers if CORS is configured correctly.