# CSRF Protection for HTMX in Fiber

This document explains how CSRF (Cross-Site Request Forgery) protection is implemented for HTMX requests in this Fiber application.

## Overview

CSRF protection prevents malicious websites from making unauthorized requests on behalf of authenticated users. This implementation uses Fiber's built-in CSRF middleware with custom configuration for HTMX compatibility.

## Implementation Details

### 1. CSRF Middleware Configuration

Located in `internal/middleware/csrfMiddleware.go`:

- **Token Header**: `X-Csrf-Token`
- **Cookie Name**: `csrf_`
- **Cookie Settings**: Lax SameSite, HTTPOnly
- **Token Lifetime**: 30 minutes
- **Custom Error Handler**: Returns HTMX-friendly error responses

### 2. Template Integration

#### Layout Template (`layout.templ`)
- Includes CSRF token in `<meta name="csrf-token">` tag
- Automatically configures HTMX to send token with all requests
- JavaScript event listener adds token to `X-Csrf-Token` header

#### CSRF Helper Templates (`csrf.templ`)
- `CSRFForm()`: Creates forms with hidden CSRF token input
- `CSRFHiddenInput()`: Standalone hidden input for manual forms
- `CSRFMeta()`: Meta tag component (included in layout)

### 3. HTMX Configuration

The layout template includes JavaScript that:
```javascript
document.body.addEventListener('htmx:configRequest', function(evt) {
    evt.detail.headers['X-Csrf-Token'] = csrfToken;
});
```

This ensures all HTMX requests (GET, POST, PUT, DELETE, etc.) include the CSRF token.

### 4. Handler Integration

All handlers that render templates automatically receive the CSRF token:
- `csrf.TokenFromContext(c)` extracts the token from Fiber context
- Token is passed to templates via the `Layout()` function
- No manual token handling required in most cases

## Usage Examples

### HTMX Buttons (Automatic Protection)
```html
<button hx-put="/counter/increment" hx-target="#counter">
    Increment
</button>
```
CSRF token is automatically included in the request header.

### HTMX Forms (Automatic Protection)
```html
<form hx-post="/api/submit" hx-target="#result">
    <input type="text" name="message" required>
    <button type="submit">Submit</button>
</form>
```
CSRF token is automatically included in the request header.

### Traditional Forms (Manual Token Required)
```html
<form action="/submit" method="POST">
    <input type="hidden" name="_token" value="{{ .CSRFToken }}">
    <input type="text" name="message" required>
    <button type="submit">Submit</button>
</form>
```

## Security Features

### 1. Token Validation
- All non-safe HTTP methods (POST, PUT, DELETE, PATCH) are validated
- Tokens expire after 30 minutes of inactivity
- Double Submit Cookie pattern used for validation

### 2. Error Handling
- HTMX requests receive user-friendly error messages
- JSON requests receive structured error responses
- Automatic page refresh suggestions for expired tokens

### 3. Cookie Security
- HTTPOnly cookies prevent XSS token theft
- Lax SameSite prevents CSRF via third-party requests
- Secure flag should be enabled in production (HTTPS)

## Testing CSRF Protection

### Valid Request Test
1. Load the page normally
2. Use HTMX buttons or forms
3. Requests should succeed with automatic token inclusion

### Invalid Token Test
1. Open browser developer tools
2. Delete the `csrf_` cookie
3. Try to use HTMX functionality
4. Should receive CSRF error message

### Manual Token Test
1. Inspect page source for `<meta name="csrf-token">`
2. Verify token is present and non-empty
3. Check network requests include `X-Csrf-Token` header

## Production Considerations

### 1. HTTPS Requirements
```go
CookieSecure: true, // Enable in production
```

### 2. Token Storage
Consider using Redis or database storage for token persistence across server restarts:
```go
Storage: your_storage_implementation,
```

### 3. Domain Configuration
For multi-domain setups, configure trusted origins:
```go
TrustedOrigins: []string{"https://yourdomain.com"},
```

## Troubleshooting

### Common Issues

1. **403 Forbidden on HTMX requests**
   - Check if CSRF token meta tag is present
   - Verify JavaScript is loading and configuring HTMX
   - Ensure token hasn't expired

2. **Token not found errors**
   - Verify CSRF middleware is loaded before routes
   - Check cookie settings and domain configuration
   - Ensure handlers are getting token from context

3. **JavaScript errors**
   - Verify meta tag selector is correct
   - Check for JavaScript conflicts
   - Ensure HTMX library is loaded

### Debug Commands
```bash
# Check if templates are generated
templ generate

# Verify build succeeds
go build ./...

# Test server startup
make run
```

## Security Best Practices

1. **Always use HTTPS in production**
2. **Set secure cookie flags appropriately**
3. **Implement proper session management**
4. **Regular security audits and updates**
5. **Monitor for unusual request patterns**
6. **Use Content Security Policy (CSP) headers**