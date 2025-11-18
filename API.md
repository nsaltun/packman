# Packman API Documentation

## Overview

The Packman API is a RESTful service that optimizes package fulfillment by calculating the most efficient combination of pack sizes for any given order quantity. The system aims to minimize the number of packs while ensuring all items are shipped.

**Base URL:** `http://localhost:{PORT}`  
**API Version:** v1  
**Content-Type:** `application/json`

### Available Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/calculate` | Calculate optimal pack combination for a given quantity |
| GET | `/api/v1/pack-sizes` | Retrieve current pack size configuration |
| PUT | `/api/v1/pack-sizes` | Update pack size configuration |
| GET | `/health` | Check service and database health status |

---

## Authentication

Currently, the API does not require authentication. All endpoints are publicly accessible.

---

## Response Format

All API responses follow a standardized JSON structure:

### Success Response

```json
{
  "data": { ... },
  "request_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

### Error Response

```json
{
  "error": {
    "code": "ERROR_CODE",
    "message": "Human-readable error message",
    "details": { ... }
  },
  "request_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

### Standard Error Codes

| Code | HTTP Status | Description |
|------|-------------|-------------|
| `VALIDATION_ERROR` | 400 | Request validation failed |
| `BAD_REQUEST` | 400 | Malformed request body |
| `NOT_FOUND` | 404 | Resource not found |
| `CONFLICT` | 409 | Resource conflict |
| `INTERNAL_ERROR` | 500 | Internal server error |
| `SERVICE_UNAVAILABLE` | 503 | Service temporarily unavailable |

---

## Endpoints

### 1. Calculate Packs

Calculates the optimal combination of pack sizes for a given order quantity.

**Endpoint:** `POST /api/v1/calculate`

#### Description

This endpoint implements an intelligent pack calculation algorithm that:
- Minimizes the total number of packs used
- Ensures all items are fulfilled (may send slightly more than requested)
- Uses the currently configured pack sizes from the database

#### Request

**Headers:**
- `Content-Type: application/json`

**Body:**
```json
{
  "quantity": 250
}
```

**Parameters:**

| Field | Type | Required | Constraints | Description |
|-------|------|----------|-------------|-------------|
| `quantity` | integer | Yes | > 0, ≤ 10,000,000 | Number of items to pack |

#### Response

**Status Code:** `200 OK`

**Body:**
```json
{
  "data": {
    "quantity": 250,
    "packs": {
      "250": 1
    }
  },
  "request_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

**Response Fields:**

| Field | Type | Description |
|-------|------|-------------|
| `quantity` | integer | The original requested quantity |
| `packs` | object | Map of pack sizes to quantities (e.g., `{"500": 1, "250": 2}` means 1 pack of 500 and 2 packs of 250) |

#### Examples

**Example 1: Simple calculation**
```bash
curl -X POST http://localhost:8080/api/v1/calculate \
  -H "Content-Type: application/json" \
  -d '{"quantity": 251}'
```

Response:
```json
{
  "data": {
    "quantity": 251,
    "packs": {
      "500": 1
    }
  },
  "request_id": "..."
}
```

**Example 2: Multiple pack sizes**
```bash
curl -X POST http://localhost:8080/api/v1/calculate \
  -H "Content-Type: application/json" \
  -d '{"quantity": 12001}'
```

Response:
```json
{
  "data": {
    "quantity": 12001,
    "packs": {
      "5000": 2,
      "2000": 1,
      "250": 1
    }
  },
  "request_id": "..."
}
```

#### Error Responses

**Validation Error (400):**
```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "quantity must be greater than zero"
  },
  "request_id": "..."
}
```

**Bad Request (400):**
```json
{
  "error": {
    "code": "BAD_REQUEST",
    "message": "Invalid request format"
  },
  "request_id": "..."
}
```

**Internal Error (500):**
```json
{
  "error": {
    "code": "INTERNAL_ERROR",
    "message": "An internal error occurred"
  },
  "request_id": "..."
}
```

---

### 2. Get Pack Sizes

Retrieves the current pack size configuration.

**Endpoint:** `GET /api/v1/pack-sizes`

#### Description

Returns the active pack sizes configuration used by the pack calculation algorithm. The configuration includes version information and audit metadata.

#### Request

**Headers:**
- None required

**Query Parameters:**
- None

#### Response

**Status Code:** `200 OK`

**Body:**
```json
{
  "data": {
    "pack_sizes": [250, 500, 1000, 2000, 5000],
    "version": 1,
    "updated_at": "2025-11-17T10:30:00Z",
    "updated_by": "admin@example.com"
  },
  "request_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

**Response Fields:**

| Field | Type | Description |
|-------|------|-------------|
| `pack_sizes` | array[integer] | List of available pack sizes in ascending order |
| `version` | integer | Configuration version number (increments with each update) |
| `updated_at` | string (ISO 8601) | Timestamp of the last configuration update |
| `updated_by` | string | Identifier of the user/system that last updated the configuration (optional) |

#### Example

```bash
curl -X GET http://localhost:8080/api/v1/pack-sizes
```

#### Error Responses

**Not Found (404):**
```json
{
  "error": {
    "code": "NOT_FOUND",
    "message": "Pack configuration not found"
  },
  "request_id": "..."
}
```

**Internal Error (500):**
```json
{
  "error": {
    "code": "INTERNAL_ERROR",
    "message": "An internal error occurred"
  },
  "request_id": "..."
}
```

---

### 3. Update Pack Sizes

Updates the pack size configuration.

**Endpoint:** `PUT /api/v1/pack-sizes`

#### Description

Updates the available pack sizes used for order fulfillment calculations. This endpoint:
- Automatically deduplicates pack sizes
- Validates all pack sizes are positive integers
- Implements optimistic locking to prevent concurrent update conflicts
- Maintains an audit trail with version history

#### Request

**Headers:**
- `Content-Type: application/json`

**Body:**
```json
{
  "pack_sizes": [250, 500, 1000, 2000, 5000],
  "updated_by": "admin@example.com"
}
```

**Parameters:**

| Field | Type | Required | Constraints | Description |
|-------|------|----------|-------------|-------------|
| `pack_sizes` | array[integer] | Yes | Non-empty, each > 0, each ≤ 1,000,000 | New pack sizes to use |
| `updated_by` | string | No | ≤ 100 characters | Identifier of who is making the update |

#### Response

**Status Code:** `200 OK`

**Body:**
```json
{
  "data": {
    "message": "Pack sizes updated successfully"
  },
  "request_id": "550e8400-e29b-41d4-a716-446655440000"
}
```

#### Examples

**Basic update**
```bash
curl -X PUT http://localhost:8080/api/v1/pack-sizes \
  -H "Content-Type: application/json" \
  -d '{
    "pack_sizes": [250, 500, 1000, 2000, 5000],
    "updated_by": "admin@example.com"
  }'
```


#### Error Responses

**Validation Error (400):**
```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "pack_sizes cannot be empty"
  },
  "request_id": "..."
}
```

**Validation Error - Invalid Size (400):**
```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "pack sizes must be greater than zero"
  },
  "request_id": "..."
}
```

**Conflict Error (409):**
```json
{
  "error": {
    "code": "CONFLICT",
    "message": "Pack configuration has been modified by another process"
  },
  "request_id": "..."
}
```

**Internal Error (500):**
```json
{
  "error": {
    "code": "INTERNAL_ERROR",
    "message": "An internal error occurred"
  },
  "request_id": "..."
}
```

---

### 4. Health Check

Checks the health status of the API and its dependencies.

**Endpoint:** `GET /health`

#### Description

Returns the operational status of the service and its database connection. This endpoint:
- Verifies database connectivity with a ping operation
- Reports connection pool statistics
- Provides response time metrics
- Returns appropriate HTTP status codes for monitoring systems

#### Request

**Headers:**
- None required

**Query Parameters:**
- None

#### Response

**Status Code:** 
- `200 OK` - Service is healthy
- `503 Service Unavailable` - Service or dependencies are unhealthy

**Body (Healthy):**
```json
{
  "status": "healthy",
  "database": {
    "status": "healthy",
    "response_time_ms": 2,
    "error": ""
  },
  "connection_pool": {
    "total_conns": 10,
    "acquired_conns": 1,
    "idle_conns": 9,
    "max_conns": 25
  }
}
```

**Body (Unhealthy):**
```json
{
  "status": "unhealthy",
  "database": {
    "status": "unhealthy",
    "response_time_ms": 5000,
    "error": "connection timeout"
  },
  "connection_pool": {
    "total_conns": 0,
    "acquired_conns": 0,
    "idle_conns": 0,
    "max_conns": 25
  }
}
```

**Response Fields:**

| Field | Type | Description |
|-------|------|-------------|
| `status` | string | Overall service status: `healthy` or `unhealthy` |
| `database.status` | string | Database connectivity status |
| `database.response_time_ms` | integer | Database ping response time in milliseconds |
| `database.error` | string | Error message if database is unhealthy (empty string if healthy) |
| `connection_pool.total_conns` | integer | Total number of connections in the pool |
| `connection_pool.acquired_conns` | integer | Number of connections currently in use |
| `connection_pool.idle_conns` | integer | Number of idle connections available |
| `connection_pool.max_conns` | integer | Maximum number of connections allowed |

#### Example

```bash
curl -X GET http://localhost:8080/health
```


## Versioning

The API uses URL path versioning (e.g., `/api/v1/`). Breaking changes will result in a new version number.

---

## For Clients

1. **Handle Errors Gracefully:** Always check for the `error` field in responses
2. **Use Request IDs:** Include the `request_id` when reporting issues
3. **Retry Logic:** Implement exponential backoff for 5xx errors
4. **Timeout Configuration:** Set appropriate request timeouts (recommended: 30s)
5. **Connection Pooling:** Reuse HTTP connections for better performance

### Example Error Handling (JavaScript)

```javascript
try {
  const response = await fetch('http://localhost:8080/api/v1/calculate', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ quantity: 250 })
  });
  
  const result = await response.json();
  
  if (result.error) {
    console.error(`Error [${result.error.code}]: ${result.error.message}`);
    console.error(`Request ID: ${result.request_id}`);
    // Handle error based on error code
  } else {
    console.log('Packs:', result.data.packs);
  }
} catch (error) {
  console.error('Network error:', error);
}
```