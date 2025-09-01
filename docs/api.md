# API Documentation

## Base URL
```
http://localhost:8080
```

## Authentication
All API endpoints (except health check) require authentication via the `X-Extension-Key` header.

```bash
X-Extension-Key: <your-secret-key>
```

## Endpoints

### Health Check
Check if the server is running and healthy.

**Endpoint:** `GET /health`

**Authentication:** Not required

**Response:**
```json
{
  "status": "ok"
}
```

**Example:**
```bash
curl http://localhost:8080/health
```

---

### Receive Walmart Orders
Receive order data from the Chrome extension.

**Endpoint:** `POST /api/walmart/orders`

**Authentication:** Required

**Request Body:**
```json
{
  "orderNumber": "123456789",
  "orderDate": "2024-01-15",
  "orderTotal": 150.00,
  "items": [
    {
      "name": "Great Value Milk",
      "price": 3.99,
      "quantity": 1,
      "category": null
    },
    {
      "name": "Bounty Paper Towels",
      "price": 12.99,
      "quantity": 2,
      "category": null
    }
  ]
}
```

**Success Response (200):**
```json
{
  "status": "success",
  "message": "Order received successfully",
  "orderId": "123456789",
  "timestamp": "2024-01-15T10:30:00Z"
}
```

**Error Responses:**

**400 Bad Request - Invalid JSON:**
```json
{
  "status": "error",
  "message": "Invalid JSON or validation error: <details>"
}
```

**400 Bad Request - Invalid Order Total:**
```json
{
  "status": "error",
  "message": "Invalid order total: must be positive"
}
```

**401 Unauthorized:**
```json
{
  "status": "error",
  "message": "Unauthorized: Missing or invalid extension key"
}
```

**Example:**
```bash
curl -X POST http://localhost:8080/api/walmart/orders \
  -H "Content-Type: application/json" \
  -H "X-Extension-Key: your-secret-key" \
  -d @sample-order.json
```

## Future Endpoints (Phase 2-3)

### List Monarch Categories
`GET /api/categories` - Fetch available Monarch Money categories

### Categorize Items
`POST /api/categorize` - Use LLM to categorize Walmart items

### Split Transaction
`POST /api/transactions/split` - Split a Monarch transaction based on Walmart order items

### Get Audit Trail
`GET /api/transactions/{id}/audit` - Get the split history for a transaction

## Rate Limiting
All endpoints are rate-limited to prevent abuse:
- 100 requests per minute per IP address
- 1000 requests per hour per API key

## Error Handling
All errors follow a consistent format:
```json
{
  "status": "error",
  "message": "Human-readable error message",
  "details": {} // Optional additional error details
}
```

## Status Codes
- `200 OK` - Request successful
- `400 Bad Request` - Invalid request data
- `401 Unauthorized` - Missing or invalid authentication
- `404 Not Found` - Resource not found
- `429 Too Many Requests` - Rate limit exceeded
- `500 Internal Server Error` - Server error