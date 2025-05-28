# Hello World API Endpoint

This document describes the simple Hello World endpoint added to verify the middleware system works correctly.

## Endpoint

```
GET /api/v1/hello
```

## Description

A simple Hello World endpoint that demonstrates:
- Basic API functionality
- Middleware integration (CORS, security headers, rate limiting)
- Query parameter handling
- JSON response format

## Parameters

| Parameter | Type   | Required | Default | Description                    |
|-----------|--------|----------|---------|--------------------------------|
| name      | string | No       | "World" | Name to include in the greeting |

## Response

```json
{
  "success": true,
  "message": "Hello, {name}!",
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

## Examples

### Basic Request

```bash
curl http://localhost:8080/api/v1/hello
```

Response:
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

### Request with Name Parameter

```bash
curl "http://localhost:8080/api/v1/hello?name=Developer"
```

Response:
```json
{
  "success": true,
  "message": "Hello, Developer!",
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

### JavaScript/Fetch Example

```javascript
// Basic request
const response = await fetch('/api/v1/hello');
const data = await response.json();
console.log(data.message); // "Hello, World!"

// With name parameter
const responseWithName = await fetch('/api/v1/hello?name=Alice');
const dataWithName = await responseWithName.json();
console.log(dataWithName.message); // "Hello, Alice!"
```

## Middleware Verification

The endpoint demonstrates that all middlewares are working:

1. **CORS**: Includes CORS headers for cross-origin requests
2. **Security Headers**: Adds security headers (X-Frame-Options, etc.)
3. **Rate Limiting**: Applies rate limiting rules
4. **Content-Type**: Sets JSON content type

## Testing

Use the provided test script to verify the endpoint:

```bash
./test_hello.sh
```

The script tests:
- Basic endpoint functionality
- Query parameter handling
- HTTP status codes
- CORS headers presence
- Security headers presence
- Rate limiting headers

## Use Cases

- **API Health Check**: Verify the API is running
- **Middleware Testing**: Confirm all middlewares are active
- **Development**: Quick test endpoint during development
- **Monitoring**: Simple endpoint for uptime checks
- **Demo**: Demonstrate API capabilities

## Security

This endpoint is public and requires no authentication. It's safe to expose as it:
- Returns only static information
- Doesn't access any sensitive data
- Includes all security headers
- Is subject to rate limiting

## Performance

The endpoint is lightweight and returns immediately with minimal processing, making it suitable for:
- Health checks
- Load balancer probes
- High-frequency monitoring