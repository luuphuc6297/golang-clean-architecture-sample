# API Testing Guide

This document provides cURL commands for testing API endpoints. You can copy and paste these commands into Postman or terminal.

## Base URL
```
http://localhost:8080
```

## 1. Health Check

### Check server status
```bash
curl -X GET http://localhost:8080/health
```

**Expected Response:**
```json
{
  "status": "ok"
}
```

## 2. Authentication Endpoints

### 2.1 User Registration
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123",
    "first_name": "John",
    "last_name": "Doe"
  }'
```

**Expected Response:**
```json
{
  "success": true,
  "data": {
    "message": "User registered successfully",
    "user": {
      "id": "uuid-here",
      "email": "user@example.com",
      "first_name": "John",
      "last_name": "Doe",
      "role": "user",
      "is_active": true,
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  }
}
```

### 2.2 User Login
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

**Expected Response:**
```json
{
  "success": true,
  "data": {
    "message": "Login successful",
    "tokens": {
      "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
      "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
      "expires_in": 1703123456
    }
  }
}
```

### 2.3 Refresh Token
```bash
curl -X POST http://localhost:8080/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "your-refresh-token-here"
  }'
```

## 3. Product Endpoints

### 3.1 List Products (Public)
```bash
curl -X GET http://localhost:8080/api/v1/products
```

### 3.2 Get Product by ID (Public)
```bash
curl -X GET http://localhost:8080/api/v1/products/{product-id}
```

### 3.3 Create Product (Protected)
```bash
curl -X POST http://localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -d '{
    "name": "Sample Product",
    "description": "Product description",
    "price": 29.99,
    "stock": 100,
    "category": "electronics"
  }'
```

### 3.4 Update Product (Protected)
```bash
curl -X PUT http://localhost:8080/api/v1/products/{product-id} \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -d '{
    "name": "Updated Product",
    "description": "Updated description",
    "price": 39.99,
    "stock": 50,
    "category": "electronics"
  }'
```

### 3.5 Delete Product (Protected)
```bash
curl -X DELETE http://localhost:8080/api/v1/products/{product-id} \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

## 4. User Management (Admin Only)

### 4.1 List Users
```bash
curl -X GET http://localhost:8080/api/v1/users \
  -H "Authorization: Bearer ADMIN_ACCESS_TOKEN"
```

### 4.2 Get User by ID
```bash
curl -X GET http://localhost:8080/api/v1/users/{user-id} \
  -H "Authorization: Bearer ADMIN_ACCESS_TOKEN"
```

### 4.3 Update User
```bash
curl -X PUT http://localhost:8080/api/v1/users/{user-id} \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer ADMIN_ACCESS_TOKEN" \
  -d '{
    "first_name": "Updated",
    "last_name": "Name",
    "role": "admin",
    "is_active": true
  }'
```

### 4.4 Delete User
```bash
curl -X DELETE http://localhost:8080/api/v1/users/{user-id} \
  -H "Authorization: Bearer ADMIN_ACCESS_TOKEN"
```

## 5. Automated Testing

Use the optimized test script:
```bash
make test-api
```

Or run directly:
```bash
chmod +x scripts/optimized-test.sh
./scripts/optimized-test.sh
```

## 6. Error Responses

All endpoints return consistent error responses:
```json
{
  "success": false,
  "error": "Error message",
  "code": "ERROR_CODE"
}
```

## 7. Authentication

- Public endpoints: No authentication required
- Protected endpoints: Require valid JWT access token
- Admin endpoints: Require admin role in addition to authentication

## 8. Rate Limiting

