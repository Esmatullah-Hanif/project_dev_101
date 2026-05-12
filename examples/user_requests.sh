#!/bin/bash

# User Service Examples
# First, run auth_requests.sh to get tokens

BASE_URL="http://localhost:8000"

# You need to set these from auth_requests.sh output
# export ACCESS_TOKEN="your_token_here"
# export USER_ID="your_user_id_here"

if [ -z "$ACCESS_TOKEN" ]; then
  echo "ERROR: ACCESS_TOKEN not set"
  echo "Run auth_requests.sh first and set tokens:"
  echo "  export ACCESS_TOKEN=<token>"
  echo "  export USER_ID=<user_id>"
  exit 1
fi

echo "========================================="
echo "User Service Examples"
echo "========================================="
echo "Using ACCESS_TOKEN: ${ACCESS_TOKEN:0:20}..."
echo "Using USER_ID: $USER_ID"

# Get User Profile
echo ""
echo "1. Get User Profile"
echo "GET /api/v1/users/$USER_ID"
curl -s -X GET "$BASE_URL/api/v1/users/$USER_ID" \
  -H "Authorization: Bearer $ACCESS_TOKEN" | jq .

# List Users
echo ""
echo "2. List Users"
echo "GET /api/v1/users?page=1&page_size=10"
curl -s -X GET "$BASE_URL/api/v1/users?page=1&page_size=10" \
  -H "Authorization: Bearer $ACCESS_TOKEN" | jq .

# Update User Profile
echo ""
echo "3. Update User Profile"
echo "PUT /api/v1/users/$USER_ID"
curl -s -X PUT "$BASE_URL/api/v1/users/$USER_ID" \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "John",
    "last_name": "Doe",
    "bio": "Software Engineer",
    "avatar_url": "https://example.com/avatar.jpg"
  }' | jq .

# Get Updated Profile
echo ""
echo "4. Get Updated User Profile"
echo "GET /api/v1/users/$USER_ID"
curl -s -X GET "$BASE_URL/api/v1/users/$USER_ID" \
  -H "Authorization: Bearer $ACCESS_TOKEN" | jq .

# Unauthorized Request (without token)
echo ""
echo "5. Unauthorized Request (Error Case)"
echo "PUT /api/v1/users/$USER_ID (without token)"
curl -s -X PUT "$BASE_URL/api/v1/users/$USER_ID" \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "Jane"
  }' | jq .

# Not Found (with invalid ID)
echo ""
echo "6. Not Found (Error Case)"
echo "GET /api/v1/users/invalid-id"
curl -s -X GET "$BASE_URL/api/v1/users/invalid-id" \
  -H "Authorization: Bearer $ACCESS_TOKEN" | jq .

# List with pagination
echo ""
echo "7. List Users with Pagination"
echo "GET /api/v1/users?page=2&page_size=5"
curl -s -X GET "$BASE_URL/api/v1/users?page=2&page_size=5" \
  -H "Authorization: Bearer $ACCESS_TOKEN" | jq .

# Delete User
echo ""
echo "8. Delete User"
echo "DELETE /api/v1/users/$USER_ID"
curl -s -X DELETE "$BASE_URL/api/v1/users/$USER_ID" \
  -H "Authorization: Bearer $ACCESS_TOKEN"

echo ""
echo ""
echo "========================================="
echo "Done!"
echo "========================================="
