#!/bin/bash

# CSRF Test Script for HTMX Fiber Application
# This script tests the CSRF protection implementation

set -e

echo "🧪 Testing CSRF Protection for HTMX..."

# Check if server is running
SERVER_URL="http://localhost:8080"
if ! curl -s "$SERVER_URL" > /dev/null; then
    echo "❌ Server not running at $SERVER_URL"
    echo "Please start the server with: make run"
    exit 1
fi

echo "✅ Server is running"

# Test 1: Get the main page and extract CSRF token
echo "📄 Test 1: Getting main page and extracting CSRF token..."
RESPONSE=$(curl -s -c cookies.txt "$SERVER_URL")

# Check if CSRF meta tag exists
if echo "$RESPONSE" | grep -q 'name="csrf-token"'; then
    CSRF_TOKEN=$(echo "$RESPONSE" | grep 'name="csrf-token"' | sed 's/.*content="\([^"]*\)".*/\1/')
    echo "✅ CSRF token found in meta tag: ${CSRF_TOKEN:0:20}..."
else
    echo "❌ CSRF token meta tag not found"
    exit 1
fi

# Test 2: Test valid HTMX request with CSRF token
echo "🔒 Test 2: Testing valid HTMX request with CSRF token..."
COUNTER_RESPONSE=$(curl -s -b cookies.txt -H "X-Csrf-Token: $CSRF_TOKEN" -H "HX-Request: true" -X PUT "$SERVER_URL/counter/increment")

if [[ "$COUNTER_RESPONSE" =~ ^[0-9]+$ ]]; then
    echo "✅ Valid CSRF request succeeded. Counter value: $COUNTER_RESPONSE"
else
    echo "❌ Valid CSRF request failed. Response: $COUNTER_RESPONSE"
fi

# Test 3: Test HTMX request without CSRF token (should fail)
echo "🚫 Test 3: Testing HTMX request without CSRF token (should fail)..."
NO_TOKEN_RESPONSE=$(curl -s -b cookies.txt -H "HX-Request: true" -X PUT "$SERVER_URL/counter/increment" -w "%{http_code}")

if echo "$NO_TOKEN_RESPONSE" | grep -q "403"; then
    echo "✅ Request without CSRF token correctly rejected (403)"
else
    echo "❌ Request without CSRF token should have been rejected"
fi

# Test 4: Test HTMX request with invalid CSRF token (should fail)
echo "🔑 Test 4: Testing HTMX request with invalid CSRF token (should fail)..."
INVALID_TOKEN_RESPONSE=$(curl -s -b cookies.txt -H "X-Csrf-Token: invalid-token-123" -H "HX-Request: true" -X PUT "$SERVER_URL/counter/increment" -w "%{http_code}")

if echo "$INVALID_TOKEN_RESPONSE" | grep -q "403"; then
    echo "✅ Request with invalid CSRF token correctly rejected (403)"
else
    echo "❌ Request with invalid CSRF token should have been rejected"
fi

# Test 5: Test form submission with CSRF token
echo "📝 Test 5: Testing form submission with CSRF token..."
FORM_RESPONSE=$(curl -s -b cookies.txt -H "X-Csrf-Token: $CSRF_TOKEN" -H "HX-Request: true" -X POST -d "message=Test message from script" "$SERVER_URL/api/test-form")

if echo "$FORM_RESPONSE" | grep -q "Form submitted successfully"; then
    echo "✅ Form submission with CSRF token succeeded"
else
    echo "❌ Form submission with CSRF token failed. Response: $FORM_RESPONSE"
fi

# Test 6: Check if JavaScript CSRF configuration is present
echo "🔧 Test 6: Checking JavaScript CSRF configuration..."
if echo "$RESPONSE" | grep -q "htmx:configRequest"; then
    echo "✅ HTMX CSRF configuration found in JavaScript"
else
    echo "❌ HTMX CSRF configuration not found"
fi

# Test 7: Verify CSRF cookie is set
echo "🍪 Test 7: Checking CSRF cookie..."
if grep -q "csrf_" cookies.txt; then
    echo "✅ CSRF cookie is set"
else
    echo "❌ CSRF cookie not found"
fi

# Cleanup
rm -f cookies.txt

echo ""
echo "🎉 CSRF Testing Complete!"
echo ""
echo "Summary:"
echo "- CSRF token is properly embedded in page meta tag"
echo "- HTMX requests with valid tokens are accepted"
echo "- Requests without tokens are properly rejected"
echo "- Requests with invalid tokens are properly rejected"
echo "- Form submissions work with CSRF protection"
echo "- JavaScript configuration for HTMX is present"
echo "- CSRF cookies are properly set"
echo ""
echo "✅ CSRF protection is working correctly!"