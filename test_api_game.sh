#!/bin/bash

# Test script for API Game Tester

echo "=== API Game Test Script ==="
echo ""

# Check if backend is running
if ! curl -s http://localhost:8080/healthz > /dev/null; then
    echo "❌ Backend is not running at localhost:8080"
    echo "Please start the backend first with: cd backend && go run main.go"
    exit 1
fi

echo "✅ Backend is running"

# Create test user and get token
echo ""
echo "Creating test user..."

# Register user
REGISTER_RESPONSE=$(curl -s -X POST http://localhost:8080/api/auth/register \
    -H "Content-Type: application/json" \
    -d '{"username":"testuser'$(date +%s)'","password":"testpass"}')

if [[ $REGISTER_RESPONSE == *"error"* ]]; then
    echo "Registration might have failed, trying login..."
fi

# Login to get token
LOGIN_RESPONSE=$(curl -s -X POST http://localhost:8080/api/auth/login \
    -H "Content-Type: application/json" \
    -d '{"username":"testuser'$(date +%s)'","password":"testpass"}')

# Extract token using grep and sed
TOKEN=$(echo "$LOGIN_RESPONSE" | grep -o '"token":"[^"]*' | sed 's/"token":"//')

if [ -z "$TOKEN" ]; then
    echo "❌ Failed to get auth token"
    echo "Response: $LOGIN_RESPONSE"
    exit 1
fi

echo "✅ Got auth token: ${TOKEN:0:20}..."

# Run the API test
echo ""
echo "Starting API game test..."
echo ""

cd backend/test
go run run_api_test.go -token "$TOKEN" -verbose

echo ""
echo "=== Test Complete ==="