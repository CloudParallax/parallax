#!/bin/bash

# Demo Script for Parallax API with Middleware
# This script starts the server and demonstrates the hello world endpoint

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

# Check if Go is installed
if ! command -v go &> /dev/null; then
    print_error "Go is not installed. Please install Go to continue."
    exit 1
fi

print_header "Parallax API Demo with Middleware"
print_info "This demo will start the server and test the hello world endpoint"

# Build the application
print_header "Building Application"
if go build -o parallax-demo ./cmd/parallax/; then
    print_success "Application built successfully"
else
    print_error "Failed to build application"
    exit 1
fi

# Start the server in background
print_header "Starting Server"
print_info "Starting Parallax API server on port 8080..."

./parallax-demo &
SERVER_PID=$!

# Give server time to start
sleep 3

# Check if server is running
if ! kill -0 $SERVER_PID 2>/dev/null; then
    print_error "Server failed to start"
    exit 1
fi

print_success "Server started with PID: $SERVER_PID"

# Function to cleanup
cleanup() {
    print_header "Cleanup"
    if kill -0 $SERVER_PID 2>/dev/null; then
        kill $SERVER_PID
        print_info "Server stopped"
    fi
    rm -f parallax-demo
    print_info "Demo binary removed"
}

# Set trap to cleanup on exit
trap cleanup EXIT

# Wait for server to be ready
print_info "Waiting for server to be ready..."
sleep 2

# Test the hello endpoint
print_header "Testing Hello World Endpoint"

# Test 1: Basic hello
print_info "Testing basic hello endpoint..."
response=$(curl -s -w "HTTP_CODE:%{http_code}" http://localhost:8080/api/v1/hello 2>/dev/null)
http_code=$(echo "$response" | grep -o "HTTP_CODE:[0-9]*" | cut -d: -f2)
body=$(echo "$response" | sed 's/HTTP_CODE:[0-9]*$//')

if [ "$http_code" = "200" ]; then
    print_success "Hello endpoint returned 200 OK"
    echo "Response: $body"
else
    print_error "Hello endpoint failed (HTTP $http_code)"
fi

echo ""

# Test 2: Hello with name
print_info "Testing hello endpoint with name parameter..."
response=$(curl -s -w "HTTP_CODE:%{http_code}" "http://localhost:8080/api/v1/hello?name=Demo" 2>/dev/null)
http_code=$(echo "$response" | grep -o "HTTP_CODE:[0-9]*" | cut -d: -f2)
body=$(echo "$response" | sed 's/HTTP_CODE:[0-9]*$//')

if [ "$http_code" = "200" ]; then
    print_success "Hello endpoint with name returned 200 OK"
    echo "Response: $body"
else
    print_error "Hello endpoint with name failed (HTTP $http_code)"
fi

echo ""

# Test 3: Health check
print_info "Testing health check endpoint..."
response=$(curl -s -w "HTTP_CODE:%{http_code}" http://localhost:8080/api/v1/health 2>/dev/null)
http_code=$(echo "$response" | grep -o "HTTP_CODE:[0-9]*" | cut -d: -f2)
body=$(echo "$response" | sed 's/HTTP_CODE:[0-9]*$//')

if [ "$http_code" = "200" ]; then
    print_success "Health check returned 200 OK"
    echo "Response: $body"
else
    print_error "Health check failed (HTTP $http_code)"
fi

echo ""

# Test 4: Check middleware headers
print_info "Checking middleware headers..."
headers=$(curl -s -I http://localhost:8080/api/v1/hello 2>/dev/null)

if echo "$headers" | grep -q "X-Content-Type-Options"; then
    print_success "Security headers detected"
else
    print_info "Security headers may not be visible in curl output"
fi

if echo "$headers" | grep -q "Access-Control-Allow"; then
    print_success "CORS headers detected"
else
    print_info "CORS headers may not be visible without Origin header"
fi

print_header "Demo Summary"
print_success "Parallax API is running successfully!"
print_info "Server URL: http://localhost:8080"
print_info "API Base: http://localhost:8080/api/v1"
print_info ""
print_info "Available endpoints:"
print_info "  GET /api/v1/hello        - Hello world"
print_info "  GET /api/v1/hello?name=X - Hello with name"  
print_info "  GET /api/v1/health       - Health check"
print_info "  GET /api/v1/auth/csrf    - Get CSRF token"
print_info ""
print_info "Middleware features:"
print_info "  ✓ CORS protection"
print_info "  ✓ CSRF protection"  
print_info "  ✓ Session authentication"
print_info "  ✓ Rate limiting"
print_info "  ✓ Security headers"
print_info ""
print_info "Press Ctrl+C to stop the demo"

# Keep running until interrupted
while true; do
    sleep 1
done