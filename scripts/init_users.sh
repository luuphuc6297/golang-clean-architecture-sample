#!/bin/bash

# Initialize default users for testing API
# Usage: ./scripts/init_users.sh <API_BASE_URL>

API_BASE_URL=${1:-"http://localhost:8080"}

echo "🚀 Initializing default users at $API_BASE_URL..."

# Function to create user and handle response
create_user() {
    local email=$1
    local password=$2
    local first_name=$3
    local last_name=$4
    local role_type=$5

    echo "📝 Creating $role_type user: $email"
    
    response=$(curl -s -w "HTTPSTATUS:%{http_code}" -X POST "$API_BASE_URL/api/v1/auth/register" \
      -H "Content-Type: application/json" \
      -d "{
        \"email\": \"$email\",
        \"password\": \"$password\",
        \"first_name\": \"$first_name\",
        \"last_name\": \"$last_name\"
      }")
    
    # Extract HTTP status code
    http_code=$(echo $response | tr -d '\n' | sed -e 's/.*HTTPSTATUS://')
    # Extract response body
    body=$(echo $response | sed -e 's/HTTPSTATUS\:.*//g')
    
    if [ $http_code -eq 201 ]; then
        echo "✅ $role_type user created successfully"
        echo "📧 Email: $email"
        echo "🔒 Password: $password"
        echo ""
    else
        echo "❌ Failed to create $role_type user (HTTP: $http_code)"
        echo "Response: $body"
        echo ""
    fi
}

# Wait for API to be ready
echo "⏳ Waiting for API to be ready..."
for i in {1..30}; do
    if curl -s "$API_BASE_URL/health" > /dev/null 2>&1; then
        echo "✅ API is ready!"
        break
    fi
    echo "   Attempt $i/30 - waiting 2 seconds..."
    sleep 2
done

# Check if API is actually ready
if ! curl -s "$API_BASE_URL/health" > /dev/null 2>&1; then
    echo "❌ API is not responding. Please check if the service is running."
    exit 1
fi

echo ""
echo "🔧 Creating default users..."
echo "================================"

# Create normal user
create_user "user@example.com" "password123" "Normal" "User" "Normal"

# Create admin user  
create_user "admin@example.com" "adminpassword" "Admin" "User" "Admin"

# Create test users
create_user "test1@example.com" "testpass123" "Test" "User1" "Test"
create_user "test2@example.com" "testpass456" "Test" "User2" "Test"

echo "================================"
echo "🎉 User initialization completed!"
echo ""
echo "📋 Summary of created users:"
echo "1. Normal User: user@example.com / password123"
echo "2. Admin User: admin@example.com / adminpassword" 
echo "3. Test User 1: test1@example.com / testpass123"
echo "4. Test User 2: test2@example.com / testpass456"
echo ""
echo "🔗 API Base URL: $API_BASE_URL"
echo "🏥 Health Check: $API_BASE_URL/health"
echo "📖 API Documentation: Check README.md for endpoint details"
