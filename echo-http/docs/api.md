# echo-http API Reference

## Base URL

```
http://localhost:80
```

## Endpoints

### GET /get

Echo request information including query parameters and headers.

**Request:**

```bash
curl "http://localhost:80/get?name=test&count=5"
```

**Response:**

```json
{
  "method": "GET",
  "url": "/get?name=test&count=5",
  "args": {
    "name": "test",
    "count": "5"
  },
  "headers": {
    "Accept": "*/*",
    "User-Agent": "curl/8.0.0"
  }
}
```

### POST /post

Echo request body with JSON or form data parsing.

**Request (JSON):**

```bash
curl -X POST http://localhost:80/post \
  -H "Content-Type: application/json" \
  -d '{"message": "hello", "count": 42}'
```

**Response:**

```json
{
  "method": "POST",
  "url": "/post",
  "args": {},
  "headers": {
    "Content-Type": "application/json"
  },
  "data": "{\"message\": \"hello\", \"count\": 42}",
  "json": {
    "message": "hello",
    "count": 42
  }
}
```

**Request (Form):**

```bash
curl -X POST http://localhost:80/post \
  -d "name=test" \
  -d "email=test@example.com"
```

**Response:**

```json
{
  "method": "POST",
  "url": "/post",
  "args": {},
  "headers": {
    "Content-Type": "application/x-www-form-urlencoded"
  },
  "data": "name=test&email=test@example.com",
  "form": {
    "name": "test",
    "email": "test@example.com"
  }
}
```

### PUT /put

Echo request body (same format as POST).

```bash
curl -X PUT http://localhost:80/put \
  -H "Content-Type: application/json" \
  -d '{"id": 1, "name": "updated"}'
```

### PATCH /patch

Echo request body (same format as POST).

```bash
curl -X PATCH http://localhost:80/patch \
  -H "Content-Type: application/json" \
  -d '{"name": "patched"}'
```

### DELETE /delete

Echo request information.

```bash
curl -X DELETE "http://localhost:80/delete?id=123"
```

### GET /headers

Return request headers only.

**Request:**

```bash
curl http://localhost:80/headers \
  -H "X-Custom-Header: custom-value" \
  -H "Authorization: Bearer token123"
```

**Response:**

```json
{
  "headers": {
    "Accept": "*/*",
    "Authorization": "Bearer token123",
    "User-Agent": "curl/8.0.0",
    "X-Custom-Header": "custom-value"
  }
}
```

### GET /status/{code}

Return the specified HTTP status code.

| Parameter | Type | Range   | Description      |
| --------- | ---- | ------- | ---------------- |
| `code`    | int  | 100-599 | HTTP status code |

**Examples:**

```bash
# 200 OK
curl -i http://localhost:80/status/200

# 404 Not Found
curl -i http://localhost:80/status/404

# 418 I'm a teapot
curl -i http://localhost:80/status/418

# 500 Internal Server Error
curl -i http://localhost:80/status/500
```

**Response:**

Returns an empty body with the specified status code.

### GET /delay/{seconds}

Echo after a specified delay. Useful for timeout testing.

| Parameter | Type  | Range | Description      |
| --------- | ----- | ----- | ---------------- |
| `seconds` | float | 0-300 | Delay in seconds |

**Examples:**

```bash
# 2 second delay
curl http://localhost:80/delay/2

# 0.5 second delay
curl http://localhost:80/delay/0.5

# Test client timeout (10 seconds)
curl --max-time 5 http://localhost:80/delay/10
```

**Response:**

Same format as `/get` but returned after the delay.

### GET /health

Health check endpoint.

**Request:**

```bash
curl http://localhost:80/health
```

**Response:**

```json
{
  "status": "ok"
}
```

## Response Format

All echo endpoints return a JSON object with the following structure:

| Field     | Type   | Description                              |
| --------- | ------ | ---------------------------------------- |
| `method`  | string | HTTP method used                         |
| `url`     | string | Request URL including query string       |
| `args`    | object | Parsed query parameters                  |
| `headers` | object | Request headers                          |
| `data`    | string | Raw request body (POST/PUT/PATCH only)   |
| `json`    | object | Parsed JSON body (if Content-Type: json) |
| `form`    | object | Parsed form body (if Content-Type: form) |

## Error Responses

| Status | Description           |
| ------ | --------------------- |
| 400    | Invalid request       |
| 404    | Endpoint not found    |
| 405    | Method not allowed    |
| 500    | Internal server error |
