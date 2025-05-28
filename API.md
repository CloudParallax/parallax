# Parallax API Documentation

A clean architecture REST API built with Go and Fiber framework.

## Base URL

```
http://localhost:8080/api/v1
```

## Quick Start

Test the API with a simple hello world endpoint:

```bash
curl http://localhost:8080/api/v1/hello
curl "http://localhost:8080/api/v1/hello?name=Developer"
```

## Authentication

The API uses session-based authentication with CSRF protection for write operations. Authentication is optional for read operations and required for write operations.

### Getting Started
1. Get CSRF token: `GET /api/v1/auth/csrf`
2. Login: `POST /api/v1/auth/login` (with CSRF token)
3. Make authenticated requests with session cookie and CSRF token

## Response Format

All API responses follow a consistent JSON format:

### Success Response
```json
{
  "success": true,
  "data": {},
  "meta": {
    "page": 1,
    "limit": 10,
    "total": 100,
    "total_pages": 10
  }
}
```

### Error Response
```json
{
  "success": false,
  "error": {
    "code": "ERROR_CODE",
    "message": "Human readable error message",
    "details": "Additional error details"
  }
}
```

## Endpoints

### Health Check

#### GET /health
Check API health status.

**Response:**
```json
{
  "status": "ok",
  "service": "parallax-api",
  "version": "1.0.0"
}
```

---

## Blog Posts

### Create Blog Post

#### POST /blog/posts
Create a new blog post.

**Request Body:**
```json
{
  "title": "My Blog Post",
  "summary": "This is a summary of my blog post",
  "content": "Full content of the blog post in markdown",
  "author": "John Doe",
  "tags": ["go", "api", "clean-architecture"]
}
```

**Response:** `201 Created`
```json
{
  "success": true,
  "data": {
    "id": "post_1234567890",
    "title": "My Blog Post",
    "slug": "my-blog-post",
    "summary": "This is a summary of my blog post",
    "content": "Full content of the blog post in markdown",
    "author": "John Doe",
    "published_at": null,
    "created_at": "2023-12-01T10:00:00Z",
    "updated_at": "2023-12-01T10:00:00Z",
    "tags": ["go", "api", "clean-architecture"],
    "read_time": 5,
    "is_published": false
  }
}
```

### Get All Blog Posts

#### GET /blog/posts
Retrieve all blog posts with optional filtering and pagination.

**Query Parameters:**
- `page` (int): Page number (default: 1)
- `limit` (int): Items per page (default: 10, max: 100)
- `author_id` (string): Filter by author ID
- `tags` (string): Comma-separated list of tags to filter by
- `is_published` (boolean): Filter by publication status
- `from_date` (string): Filter posts from date (YYYY-MM-DD)
- `to_date` (string): Filter posts to date (YYYY-MM-DD)
- `q` (string): Search query for title, summary, or content
- `sort_by` (string): Sort by field (created_at, updated_at, published_at, title)
- `sort_order` (string): Sort order (asc, desc)

**Response:** `200 OK`
```json
{
  "success": true,
  "data": {
    "posts": [
      {
        "id": "post_1234567890",
        "title": "My Blog Post",
        "slug": "my-blog-post",
        "summary": "This is a summary of my blog post",
        "author": "John Doe",
        "published_at": "2023-12-01T10:00:00Z",
        "created_at": "2023-12-01T10:00:00Z",
        "tags": ["go", "api"],
        "read_time": 5,
        "is_published": true
      }
    ]
  },
  "meta": {
    "page": 1,
    "limit": 10,
    "total": 1,
    "total_pages": 1
  }
}
```

### Get Blog Post by ID

#### GET /blog/posts/{id}
Retrieve a specific blog post by its ID.

**Response:** `200 OK`
```json
{
  "success": true,
  "data": {
    "id": "post_1234567890",
    "title": "My Blog Post",
    "slug": "my-blog-post",
    "summary": "This is a summary of my blog post",
    "content": "Full content of the blog post in markdown",
    "author": "John Doe",
    "published_at": "2023-12-01T10:00:00Z",
    "created_at": "2023-12-01T10:00:00Z",
    "updated_at": "2023-12-01T10:00:00Z",
    "tags": ["go", "api", "clean-architecture"],
    "read_time": 5,
    "is_published": true
  }
}
```

### Get Blog Post by Slug

#### GET /blog/posts/slug/{slug}
Retrieve a specific blog post by its slug.

**Response:** Same as GET by ID

### Update Blog Post

#### PUT /blog/posts/{id}
Update an existing blog post.

**Request Body:**
```json
{
  "title": "Updated Blog Post Title",
  "summary": "Updated summary",
  "content": "Updated content",
  "tags": ["go", "api", "updated"]
}
```

**Response:** `200 OK` - Same format as Create

### Delete Blog Post

#### DELETE /blog/posts/{id}
Delete a blog post.

**Response:** `204 No Content`

### Publish Blog Post

#### POST /blog/posts/{id}/publish
Publish a blog post.

**Response:** `200 OK` - Returns updated blog post

### Unpublish Blog Post

#### POST /blog/posts/{id}/unpublish
Unpublish a blog post.

**Response:** `200 OK` - Returns updated blog post

### Get Published Posts

#### GET /blog/posts/published
Get only published blog posts.

**Query Parameters:**
- `page` (int): Page number (default: 1)
- `limit` (int): Items per page (default: 10, max: 100)

**Response:** Same format as Get All Posts

### Search Blog Posts

#### GET /blog/posts/search
Search blog posts by query.

**Query Parameters:**
- `q` (string, required): Search query
- `page` (int): Page number (default: 1)
- `limit` (int): Items per page (default: 10, max: 100)

**Response:** Same format as Get All Posts

### Get Posts by Tags

#### GET /blog/posts/tags
Get blog posts by tags.

**Query Parameters:**
- `tags` (string, required): Comma-separated list of tags

**Response:** Same format as Get All Posts

---

## Counters

### Create Counter

#### POST /counter
Create a new counter.

**Request Body:**
```json
{
  "id": "my-counter",
  "initial_value": 0,
  "min_value": 0,
  "max_value": 100
}
```

**Response:** `201 Created`
```json
{
  "success": true,
  "data": {
    "id": "my-counter",
    "value": 0,
    "min_value": 0,
    "max_value": 100,
    "is_at_minimum": true,
    "is_at_maximum": false,
    "created_at": "2023-12-01T10:00:00Z",
    "updated_at": "2023-12-01T10:00:00Z"
  }
}
```

### Get All Counters

#### GET /counter
Retrieve all counters.

**Response:** `200 OK`
```json
{
  "success": true,
  "data": {
    "counters": [
      {
        "id": "my-counter",
        "value": 5,
        "min_value": 0,
        "max_value": 100,
        "is_at_minimum": false,
        "is_at_maximum": false,
        "created_at": "2023-12-01T10:00:00Z",
        "updated_at": "2023-12-01T10:05:00Z"
      }
    ]
  }
}
```

### Get Counter by ID

#### GET /counter/{id}
Retrieve a specific counter.

**Response:** `200 OK`
```json
{
  "success": true,
  "data": {
    "id": "my-counter",
    "value": 5,
    "min_value": 0,
    "max_value": 100,
    "is_at_minimum": false,
    "is_at_maximum": false,
    "created_at": "2023-12-01T10:00:00Z",
    "updated_at": "2023-12-01T10:05:00Z"
  }
}
```

### Increment Counter

#### PUT /counter/{id}/increment
Increment counter value by 1.

**Response:** `200 OK`
```json
{
  "success": true,
  "data": {
    "counter": {
      "id": "my-counter",
      "value": 6,
      "min_value": 0,
      "max_value": 100,
      "is_at_minimum": false,
      "is_at_maximum": false,
      "created_at": "2023-12-01T10:00:00Z",
      "updated_at": "2023-12-01T10:06:00Z"
    },
    "operation": "increment",
    "success": true,
    "message": "Counter incremented successfully"
  }
}
```

### Decrement Counter

#### PUT /counter/{id}/decrement
Decrement counter value by 1.

**Response:** `200 OK` - Same format as Increment

### Set Counter Value

#### PUT /counter/{id}/value
Set counter to a specific value.

**Request Body:**
```json
{
  "value": 25
}
```

**Response:** `200 OK` - Same format as Increment with operation "set_value"

### Reset Counter

#### PUT /counter/{id}/reset
Reset counter to specified value or minimum value.

**Query Parameters:**
- `value` (int, optional): Value to reset to (defaults to minimum value)

**Response:** `200 OK` - Same format as Increment with operation "reset"

### Delete Counter

#### DELETE /counter/{id}
Delete a counter.

**Response:** `204 No Content`

---

## Error Codes

| Code | Description | HTTP Status |
|------|-------------|-------------|
| `VALIDATION_FAILED` | Request validation failed | 400 |
| `INVALID_INPUT` | Invalid input provided | 400 |
| `MISSING_FIELD` | Required field is missing | 400 |
| `INVALID_JSON` | Invalid JSON format | 400 |
| `NOT_FOUND` | Resource not found | 404 |
| `ALREADY_EXISTS` | Resource already exists | 409 |
| `CONFLICT` | Resource conflict | 409 |
| `UNPROCESSABLE_ENTITY` | Business rule violation | 422 |
| `INTERNAL_SERVER_ERROR` | Internal server error | 500 |

---

## Examples

### Create and Manage a Blog Post

```bash
# Create a blog post
curl -X POST http://localhost:3000/api/v1/blog/posts \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Getting Started with Go",
    "summary": "Learn the basics of Go programming",
    "content": "# Introduction\n\nGo is a powerful programming language...",
    "author": "Jane Doe",
    "tags": ["go", "programming", "tutorial"]
  }'

# Get the blog post
curl http://localhost:3000/api/v1/blog/posts/post_1234567890

# Publish the blog post
curl -X POST http://localhost:3000/api/v1/blog/posts/post_1234567890/publish

# Search blog posts
curl "http://localhost:3000/api/v1/blog/posts/search?q=go&limit=5"
```

### Create and Use a Counter

```bash
# Create a counter
curl -X POST http://localhost:3000/api/v1/counter \
  -H "Content-Type: application/json" \
  -d '{
    "id": "page-views",
    "initial_value": 0,
    "min_value": 0,
    "max_value": 1000000
  }'

# Increment the counter
curl -X PUT http://localhost:3000/api/v1/counter/page-views/increment

# Get counter value
curl http://localhost:3000/api/v1/counter/page-views

# Set specific value
curl -X PUT http://localhost:3000/api/v1/counter/page-views/value \
  -H "Content-Type: application/json" \
  -d '{"value": 100}'
```

---

## Development

### Running the API

```bash
# Set environment variables
export SERVER_PORT=3000

# Run the application
go run cmd/parallax/parallax.go
```

### Testing with curl

```bash
# Health check
curl http://localhost:3000/api/v1/health

# Get all blog posts
curl http://localhost:3000/api/v1/blog/posts

# Get all counters
curl http://localhost:3000/api/v1/counter
```