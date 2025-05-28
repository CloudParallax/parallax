#!/bin/bash

# Simple Hello World API Test Script

BASE_URL="http://localhost:8080/api/v1"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

print_header() {
    echo -e "\n${BLUE}=== $1 ===${NC}"
}

print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

print_error() {
    echo -e "${RED}✗ $1${NC}"
}

print_info() {
    echo -e "${YELLOW}ℹ $1${NC}"
}

print_header "Testing Hello World API Endpoint"

# Test 1: Basic Hello World
print_header "Test 1: Basic Hello World"
response=$(curl -s "$BASE_URL/hello")
status=$(curl -s -w "%{http_code}" -o /dev/null "$BASE_URL/hello")

if [ "$status" = "200" ]; then
    print_success "Hello endpoint returns 200 OK"
    echo "Response: $response"
else
    print_error "Hello endpoint failed (HTTP $status)"
fi

# Test 2: Hello World with name parameter
print_header "Test 2: Hello World with name parameter"
response=$(curl -s "$BASE_URL/hello?name=Developer")
status=$(curl -s -w "%{http_code}" -o /dev/null "$BASE_URL/hello?name=Developer")

if [ "$status" = "200" ]; then
    print_success "Hello endpoint with name parameter returns 200 OK"
    echo "Response: $response"
else
    print_error "Hello endpoint with name parameter failed (HTTP $status)"
fi

# Test 3: Health Check
print_header "Test 3: Health Check"
response=$(curl -s "$BASE_URL/health")
status=$(curl -s -w "%{http_code}" -o /dev/null "$BASE_URL/health")

if [ "$status" = "200" ]; then
    print_success "Health check returns 200 OK"
    echo "Response: $response"
else
    print_error "Health check failed (HTTP $status)"
fi

# Test 4: CORS Headers
print_header "Test 4: CORS Headers"
cors_headers=$(curl -s -I -H "Origin: http://localhost:3000" "$BASE_URL/hello" | grep -i "access-control")

if [ -n "$cors_headers" ]; then
    print_success "CORS headers present"
    echo "$cors_headers"
else
    print_error "CORS headers not found"
fi

# Test 5: Security Headers
print_header "Test 5: Security Headers"
security_headers=$(curl -s -I "$BASE_URL/hello" | grep -E "(X-Content-Type-Options|X-Frame-Options|X-XSS-Protection)")

if [ -n "$security_headers" ]; then
    print_success "Security headers present"
    echo "$security_headers"
else
    print_error "Security headers not found"
fi

# Test 6: Rate Limiting Headers
print_header "Test 6: Rate Limiting Headers"
rate_limit_headers=$(curl -s -I "$BASE_URL/hello" | grep -i "x-ratelimit")

if [ -n "$rate_limit_headers" ]; then
    print_success "Rate limiting headers present"
    echo "$rate_limit_headers"
else
    print_info "Rate limiting headers not found (may be normal)"
fi

print_header "Test Summary"
print_info "API Base URL: $BASE_URL"
print_info "All basic tests completed!"
print_info ""
print_info "To start the server: go run cmd/parallax/parallax.go"
print_info "To test manually:"
print_info "  curl $BASE_URL/hello"
print_info "  curl $BASE_URL/hello?name=YourName"
print_info "  curl $BASE_URL/health"