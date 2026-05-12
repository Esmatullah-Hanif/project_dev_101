#!/bin/bash

# Auth Service Examples
# Make sure services are running on ports 8000-8002

BASE_URL="http://localhost:8000"

echo "========================================="
echo "Auth Service Examples"
echo "========================================="

# Sign Up
echo ""
echo "1. Sign Up"
echo "POST /api/v1/auth/signup"
SIGNUP_RESPONSE=$(curl -s -X POST "$BASE_URL/api/v1/auth/signup" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "testuser@example.com",
    "password": "password123"
  }')
echo "$SIGNUP_RESPONSE" | jq .

# Extract tokens
ACCESS_TOKEN=$(echo "$SIGNUP_RESPONSE" | jq -r '.data.access_token')
REFRESH_TOKEN=$(echo "$SIGNUP_RESPONSE" | jq -r '.data.refresh_token')
USER_ID=$(echo "$SIGNUP_RESPONSE" | jq -r '.data.user.id')

echo ""
echo "Extracted:"
echo "ACCESS_TOKEN: $ACCESS_TOKEN"
echo "REFRESH_TOKEN: $REFRESH_TOKEN"
echo "USER_ID: $USER_ID"

# Sign In
echo ""
echo "2. Sign In"
echo "POST /api/v1/auth/signin"
curl -s -X POST "$BASE_URL/api/v1/auth/signin" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "testuser@example.com",
    "password": "password123"
  }' | jq .

# Refresh Token
echo ""
echo "3. Refresh Token"
echo "POST /api/v1/auth/refresh"
curl -s -X POST "$BASE_URL/api/v1/auth/refresh" \
  -H "Content-Type: application/json" \
  -d "{
    \"refresh_token\": \"$REFRESH_TOKEN\"
  }" | jq .

# Invalid Credentials
echo ""
echo "4. Sign In with Invalid Credentials (Error Case)"
echo "POST /api/v1/auth/signin"
curl -s -X POST "$BASE_URL/api/v1/auth/signin" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "testuser@example.com",
    "password": "wrongpassword"
  }' | jq .

# Duplicate Email
echo ""
echo "5. Sign Up with Duplicate Email (Error Case)"
echo "POST /api/v1/auth/signup"
curl -s -X POST "$BASE_URL/api/v1/auth/signup" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "testuser@example.com",
    "password": "password456"
  }' | jq .

echo ""
echo "========================================="
echo "Done!"
echo "========================================="
