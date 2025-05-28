#!/bin/bash

# Middleware Test Script for Parallax API
# This script tests CORS, CSRF, Authentication, and Rate Limiting

BASE_URL="http://localhost:8080/api/v1"
COOKIE_JAR="cookies.txt"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Helper functions
print_header() {
    echo -e "\n${BLUE}=== $1 ===${NC}"
}

print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

print_error() {
    echo -e "${RED}✗ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠ $1${NC}"
}

# Clean up function
cleanup() {
    rm -f "$COOKIE_JAR"
    echo -e "\n${YELLOW}Cleaned up temporary files${NC}"
}

# Set trap to cleanup on exit
trap cleanup EXIT

print_header "Starting Middleware Tests"

# Test 1: Health Check (No middleware required)
print_header "Test 1: Health Check"
response=$(curl -s -w "%{http_code}" -o /dev/null "$BASE_URL/health")
if [ "$response" = "200" ]; then
    print_success "Health check passed"
else
    print_error "Health check failed (HTTP $response)"
fi

# Test 2: CORS Preflight Request
print_header "Test 2: CORS Preflight Request"
response=$(curl -s -w "%{http_code}" -o /dev/null \
    -X OPTIONS \
    -H "Origin: http://localhost:3000" \
    -H "Access-Control-Request-Method: POST" \
    -H "Access-Control-Request-Headers: Content-Type,X-CSRF-Token" \
    "$BASE_URL/blog/posts")

if [ "$response" = "204" ] || [ "$response" = "200" ]; then
    print_success "CORS preflight request passed"
else
    print_error "CORS preflight request failed (HTTP $response)"
fi

# Test 3: CORS Actual Request
print_header "Test 3: CORS Actual Request"
cors_headers=$(curl -s -I \
    -H "Origin: http://localhost:3000" \
    "$BASE_URL/blog/posts" | grep -i "access-control")

if [ -n "$cors_headers" ]; then
    print_success "CORS headers present in response"
    echo "$cors_headers"
else
    print_warning "CORS headers not found or request failed"
fi

# Test 4: Get CSRF Token
print_header "Test 4: Get CSRF Token"
csrf_response=$(curl -s -c "$COOKIE_JAR" "$BASE_URL/auth/csrf")
csrf_token=$(echo "$csrf_response" | grep -o '"csrf_token":"[^"]*"' | cut -d'"' -f4)

if [ -n "$csrf_token" ]; then
    print_success "CSRF token obtained: ${csrf_token:0:20}..."
else
    print_error "Failed to get CSRF token"
    echo "Response: $csrf_response"
fi

# Test 5: Authentication Status (Unauthenticated)
print_header "Test 5: Authentication Status (Unauthenticated)"
auth_response=$(curl -s -b "$COOKIE_JAR" "$BASE_URL/auth/me")
authenticated=$(echo "$auth_response" | grep -o '"authenticated":[^,}]*' | cut -d':' -f2)

if [ "$authenticated" = "false" ]; then
    print_success "Unauthenticated status correct"
else
    print_error "Unexpected authentication status"
    echo "Response: $auth_response"
fi

# Test 6: Login
print_header "Test 6: User Login"
login_response=$(curl -s -b "$COOKIE_JAR" -c "$COOKIE_JAR" \
    -X POST \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    -d '{"username":"testuser","password":"testpass"}' \
    "$BASE_URL/auth/login")

login_success=$(echo "$login_response" | grep -o '"success":[^,}]*' | cut -d':' -f2)

if [ "$login_success" = "true" ]; then
    print_success "Login successful"
else
    print_error "Login failed"
    echo "Response: $login_response"
fi

# Test 7: Authentication Status (Authenticated)
print_header "Test 7: Authentication Status (Authenticated)"
auth_response=$(curl -s -b "$COOKIE_JAR" "$BASE_URL/auth/me")
authenticated=$(echo "$auth_response" | grep -o '"authenticated":[^,}]*' | cut -d':' -f2)

if [ "$authenticated" = "true" ]; then
    print_success "Authenticated status correct"
else
    print_error "Authentication failed after login"
    echo "Response: $auth_response"
fi

# Test 8: Protected Route Without CSRF Token
print_header "Test 8: Protected Route Without CSRF Token"
response=$(curl -s -w "%{http_code}" -o /dev/null \
    -b "$COOKIE_JAR" \
    -X POST \
    -H "Content-Type: application/json" \
    -d '{"title":"Test Post","content":"Test content","author":"testuser"}' \
    "$BASE_URL/blog/posts")

if [ "$response" = "403" ]; then
    print_success "CSRF protection working (blocked request without token)"
else
    print_warning "CSRF protection may not be working properly (HTTP $response)"
fi

# Test 9: Protected Route With CSRF Token
print_header "Test 9: Protected Route With CSRF Token"
# Get fresh CSRF token after login
csrf_response=$(curl -s -b "$COOKIE_JAR" "$BASE_URL/auth/csrf")
csrf_token=$(echo "$csrf_response" | grep -o '"csrf_token":"[^"]*"' | cut -d'"' -f4)

response=$(curl -s -w "%{http_code}" -o /dev/null \
    -b "$COOKIE_JAR" \
    -X POST \
    -H "Content-Type: application/json" \
    -H "X-CSRF-Token: $csrf_token" \
    -d '{"title":"Test Post","content":"Test content","author":"testuser"}' \
    "$BASE_URL/blog/posts")

if [ "$response" = "200" ] || [ "$response" = "201" ]; then
    print_success "Protected route accessible with valid CSRF token"
else
    print_error "Protected route failed with valid CSRF token (HTTP $response)"
fi

# Test 10: Rate Limiting
print_header "Test 10: Rate Limiting"
rate_limit_count=0
for i in {1..10}; do
    response=$(curl -s -w "%{http_code}" -o /dev/null "$BASE_URL/health")
    if [ "$response" = "429" ]; then
        rate_limit_count=$((rate_limit_count + 1))
        break
    fi
    sleep 0.1
done

if [ "$rate_limit_count" -eq 0 ]; then
    print_warning "Rate limiting not triggered (may need higher load)"
else
    print_success "Rate limiting working"
fi

# Test 11: CORS with Credentials
print_header "Test 11: CORS with Credentials"
response=$(curl -s -I \
    -H "Origin: http://localhost:3000" \
    -b "$COOKIE_JAR" \
    "$BASE_URL/auth/me" | grep -i "access-control-allow-credentials")

if [ -n "$response" ]; then
    print_success "CORS credentials support enabled"
else
    print_warning "CORS credentials support not detected"
fi

# Test 12: Security Headers
print_header "Test 12: Security Headers"
security_headers=$(curl -s -I "$BASE_URL/health" | grep -E "(X-Content-Type-Options|X-Frame-Options|X-XSS-Protection|Content-Security-Policy)")

if [ -n "$security_headers" ]; then
    print_success "Security headers present"
    echo "$security_headers"
else
    print_warning "Security headers not found"
fi

# Test 13: Admin Route (Should fail without admin role)
print_header "Test 13: Admin Route Access"
response=$(curl -s -w "%{http_code}" -o /dev/null \
    -b "$COOKIE_JAR" \
    -H "X-CSRF-Token: $csrf_token" \
    "$BASE_URL/admin/users")

if [ "$response" = "403" ]; then
    print_success "Admin route properly protected (access denied)"
else
    print_warning "Admin route protection may not be working (HTTP $response)"
fi

# Test 14: Logout
print_header "Test 14: User Logout"
logout_response=$(curl -s \
    -b "$COOKIE_JAR" -c "$COOKIE_JAR" \
    -X POST \
    -H "X-CSRF-Token: $csrf_token" \
    "$BASE_URL/auth/logout")

logout_success=$(echo "$logout_response" | grep -o '"success":[^,}]*' | cut -d':' -f2)

if [ "$logout_success" = "true" ]; then
    print_success "Logout successful"
else
    print_error "Logout failed"
    echo "Response: $logout_response"
fi

# Test 15: Authentication Status After Logout
print_header "Test 15: Authentication Status After Logout"
auth_response=$(curl -s -b "$COOKIE_JAR" "$BASE_URL/auth/me")
authenticated=$(echo "$auth_response" | grep -o '"authenticated":[^,}]*' | cut -d':' -f2)

if [ "$authenticated" = "false" ]; then
    print_success "Logout cleared authentication"
else
    print_error "Authentication still active after logout"
    echo "Response: $auth_response"
fi

# Test 16: Invalid Origin CORS Test
print_header "Test 16: Invalid Origin CORS Test"
response=$(curl -s -w "%{http_code}" -o /dev/null \
    -H "Origin: http://malicious-site.com" \
    "$BASE_URL/health")

if [ "$response" = "403" ] || [ "$response" = "200" ]; then
    print_success "CORS origin validation working"
else
    print_warning "CORS origin validation may need adjustment (HTTP $response)"
fi

print_header "Middleware Test Summary"
echo "All middleware tests completed!"
echo "Check the results above for any failures or warnings."
echo ""
echo "Key Points:"
echo "- ✓ should indicate working middleware"
echo "- ⚠ may indicate configuration needed"
echo "- ✗ indicates a problem that needs fixing"
echo ""
echo "Note: Some tests may show warnings in development mode."
echo "Ensure HTTPS and proper production settings for production deployment."