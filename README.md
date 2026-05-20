# Go REST API Server

A clean, well-structured REST API server built with Go and the Gin web framework. This project demonstrates best practices for API development with a layered architecture pattern.

## 📋 Table of Contents

- [Overview](#overview)
- [Architecture](#architecture)
- [Project Structure](#project-structure)
- [Prerequisites](#prerequisites)
- [Installation & Setup](#installation--setup)
- [Running the Application](#running-the-application)
- [API Endpoints](#api-endpoints)
- [Authentication](#authentication)
- [Testing](#testing)
- [Directory Overview](#directory-overview)
- [Dependencies](#dependencies)

## Overview

This is a simple yet well-organized REST API server that provides user management functionality. It's built with:

- **Language**: Go 1.26.2
- **Web Framework**: Gin v1.12.0
- **Architecture Pattern**: Layered/Clean Architecture
- **Data Storage**: In-memory repository

The application runs on port `8080` and provides both health checks and user management endpoints.

## Architecture

The project follows a **layered architecture pattern** with clear separation of concerns:

```
API Request
    ↓
Routes (Entry point, dependency injection)
    ↓
Controllers (Request/response handling)
    ↓
Services (Business logic)
    ↓
Repositories (Data access)
    ↓
Models (Data structures)
```

### Key Principles:

- **Separation of Concerns**: Each layer has a specific responsibility
- **Dependency Injection**: Dependencies are injected through constructors
- **Interface-based Design**: Services use interfaces for loose coupling
- **DTOs**: Data Transfer Objects separate API contracts from models
- **Error Handling**: Consistent error responses through utility functions

## Project Structure

```
go-server/
├── main.go                          # Application entry point
├── go.mod                           # Go module definition
├── README.md                        # This file
├── api_load_test.go                 # API performance/load tests
├── controllers/                     # HTTP request handlers
│   ├── health_controller.go         # Health check endpoints
│   └── user_controller.go           # User CRUD endpoints
├── models/                          # Data structures
│   └── user.go                      # User model
├── services/                        # Business logic layer
│   └── user_service.go              # User service interface & implementation
├── repositories/                    # Data access layer
│   └── user_repository.go           # User data repository (in-memory)
├── routes/                          # Route registration & DI setup
│   ├── health_routes.go             # Health check route setup
│   └── user_routes.go               # User routes setup
├── dto/                             # Data Transfer Objects
│   └── user_dto.go                  # User request/response DTOs
└── utils/                           # Utility functions
    └── response.go                  # API response helpers
```

## Prerequisites

- **Go 1.26.2** or higher
- **Git** (for version control)

## Installation & Setup

1. **Clone or navigate to the project directory:**
   ```bash
   cd /Users/wasiulislam/Documents/Go
   ```

2. **Install dependencies:**
   ```bash
   go mod download
   ```

3. **Tidy up dependencies (remove unused):**
   ```bash
   go mod tidy
   ```

## Running the Application

Start the server:

```bash
go run main.go
```

The server will start on `http://localhost:8080`

You should see output like:
```
[GIN-debug] Loaded HTML Templates (0): 
[GIN-debug] Loaded HTML Templates (0): 
[GIN-debug] Listening and serving HTTP on :8080
```

## API Endpoints

The `GET /users` endpoint requires this request header:

```http
Authorization: Bearer Gib me access
```

The `/health`, `POST /users`, `GET /users/:id`, and `DELETE /users/:id` endpoints do not require authentication.

### Health Check

### User Management

| Method | Endpoint | Auth Required | Description | Request Body |
|--------|----------|---------------|-------------|--------------|
| GET | `/users` | Yes | Fetch all users | - |
| POST | `/users` | No | Create a new user | `{ "name": string, "age": number }` |
| GET | `/users/:id` | No | Fetch user by ID | - |
| DELETE | `/users/:id` | No | Delete user by ID | - |

### Example Requests

**Get all users:**
```bash
curl -X GET http://localhost:8080/users \
  -H "Authorization: Bearer Gib me access"
```

**Create a user:**
```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"name": "John Doe", "age": 28}'
```

**Get user by ID:**
```bash
curl -X GET http://localhost:8080/users/1
```

**Delete user:**
```bash
curl -X DELETE http://localhost:8080/users/1
```

### API Response Format

All endpoints return a standardized JSON response:

**Success Response:**
```json
{
  "success": true,
  "message": "Users fetched successfully",
  "data": [
    {
      "id": 1,
      "name": "John Doe",
      "age": 28
    }
  ]
}
```

**Error Response:**
```json
{
  "success": false,
  "message": "Bad Request",
  "error": "Invalid request body"
}
```

## Authentication

Only the `GET /users` endpoint is protected by auth middleware. Requests to this endpoint must include the following header:

```http
Authorization: Bearer Gib me access
```

Requests to `GET /users` without this header, with a different token, or with an invalid bearer format will be rejected.

Example unauthorized request:

```bash
curl -X GET http://localhost:8080/users
```

Example authorized request:

```bash
curl -X GET http://localhost:8080/users \
  -H "Authorization: Bearer Gib me access"
```

## Testing

This project supports both regular Go tests and API performance/load tests.

### Test Types

| Test Type | Purpose | Command |
|----------|---------|---------|
| Unit/API tests | Test project behavior with Go test runner | `go test ./...` |
| Verbose tests | Show each test name and output | `go test ./... -v` |
| Race detection | Detect unsafe concurrent access | `go test ./... -race` |
| Benchmarks | Measure Go benchmark performance | `go test ./... -bench=. -benchmem` |

### API Performance Test

The project includes an API load test file:

```text
api_load_test.go
```

This file sends multiple HTTP requests to the running server and prints performance information such as:

```text
Total Requests
Success Requests
Failed Requests
Total Duration
Requests/sec
```

Example output:

```text
========================================
API Load Test: GET /users
========================================
Total Requests:   1000
Success Requests: 1000
Failed Requests:  0
Total Duration:   82.893208ms
Requests/sec:     12063.71
========================================
```

### Running API Performance Tests

First, start the API server:

```bash
go run main.go
```

Then open another terminal and run:

```bash
go test -v
```

The API load tests call the real server at:

```text
http://localhost:8080
```

### Current API Load Test Coverage

Because `GET /users` is protected, API tests should include this header for requests to `GET /users`:

```http
Authorization: Bearer Gib me access
```

The create-user tests track users created during testing and attempt to clean them up afterward using:

```text
DELETE /users/:id
```

```bash
go test ./... -race -v
```

### Recommended API Test Cases

| Endpoint | Test Case | Expected Result |
|----------|-----------|-----------------|
| GET `/health` | Server health check | HTTP 200 |
| GET `/users` | Missing auth header | HTTP 401 |
| GET `/users` | Invalid auth token | HTTP 401 |
| GET `/users` | Valid auth header | HTTP 200 with user list |
| POST `/users` | Create valid user | HTTP 201 or 200 with created user |
| POST `/users` | Invalid JSON body | HTTP 400 |
| POST `/users` | Missing required fields | HTTP 400 |
| GET `/users/:id` | Existing user ID | HTTP 200 |
| GET `/users/:id` | Invalid ID format | HTTP 400 |
| GET `/users/:id` | Non-existing user ID | HTTP 404 |
| DELETE `/users/:id` | Existing user ID | HTTP 200 or 204 |
| DELETE `/users/:id` | Non-existing user ID | HTTP 404 |

## Directory Overview

### `main.go`
The application entry point that:
- Initializes a new Gin engine
- Registers all routes (health and user)
- Starts the HTTP server on port 8080

### `controllers/`
Contains HTTP request handlers that:
- Validate incoming requests
- Call the appropriate service methods
- Return formatted responses
- Handle errors gracefully

**Files:**
- `health_controller.go`: Simple health check endpoint
- `user_controller.go`: User CRUD operations (GetUsers, CreateUser, GetUserByID, DeleteUser)

### `services/`
Contains business logic that:
- Implements core application logic
- Validates business rules
- Coordinates between controllers and repositories
- Uses interfaces for loose coupling

**Files:**
- `user_service.go`: Defines UserService interface and implements business logic for user operations

### `repositories/`
Contains data access layer that:
- Abstracts data storage (currently in-memory)
- Provides clean data access interface
- Can be easily replaced with database implementations

**Files:**
- `user_repository.go`: In-memory implementation of user data access

### `models/`
Contains data structure definitions:
- `user.go`: User model with JSON tags for serialization

### `routes/`
Sets up API routes and dependency injection:
- Instantiates repositories, services, and controllers
- Registers routes with the Gin engine
- Manages dependency wiring

**Files:**
- `health_routes.go`: Health check endpoint setup
- `user_routes.go`: User CRUD endpoints setup with full DI

### `dto/`
Data Transfer Objects for API contracts:
- `user_dto.go`: Request and response DTOs for user operations
- Decouples API contracts from internal models

### `utils/`
Utility functions and helpers:
- `response.go`: Standardized API response helper functions (Success, Error, BadRequest, UnauthorizedAccess)

## Dependencies

The project uses the following Go packages:

| Package | Version | Purpose |
|---------|---------|---------|
| `github.com/gin-gonic/gin` | v1.12.0 | Web framework for building HTTP APIs |
| `github.com/gin-contrib/sse` | v1.1.0 | Server-Sent Events support for Gin |
| Various indirect dependencies | - | Supporting libraries for JSON serialization, validation, and utilities |

Run `go mod tidy` to ensure all dependencies are properly resolved.

## Best Practices Demonstrated

✅ **Layered Architecture**: Clear separation between controllers, services, and repositories  
✅ **Dependency Injection**: Dependencies passed through constructors, not created internally  
✅ **Interface-based Design**: Services defined as interfaces for testability and flexibility  
✅ **Error Handling**: Consistent error responses across all endpoints  
✅ **DTOs**: Request/response objects separate from internal models  
✅ **Code Organization**: Files grouped by functionality  
✅ **Standardized Responses**: All endpoints follow a consistent response format  
✅ **API Performance Testing**: Load tests measure request throughput and success/failure counts  
✅ **Clean Code**: Clear naming, single responsibility principle  

## Future Enhancements

- Replace static bearer-token authentication with JWT-based authentication and authorization
- Add input validation and sanitization
- Add logging and monitoring
- Expand automated testing with deeper unit, integration, and database-backed tests
- Add API documentation with Swagger
- Add middleware for CORS, rate limiting, etc.

---

**Author**: Go Server Project  
**Last Updated**: May 16, 2026
