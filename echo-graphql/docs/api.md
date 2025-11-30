# echo-graphql API Reference

## Endpoints

| Path       | Description        |
| ---------- | ------------------ |
| `/`        | GraphQL Playground |
| `/graphql` | GraphQL endpoint   |
| `/health`  | Health check       |

## Schema

### Types

#### Message

```graphql
type Message {
  id: ID!
  text: String!
  createdAt: String!
}
```

| Field       | Type    | Description               |
| ----------- | ------- | ------------------------- |
| `id`        | ID!     | Unique message identifier |
| `text`      | String! | Message content           |
| `createdAt` | String! | ISO 8601 timestamp        |

#### EchoResult

```graphql
type EchoResult {
  message: String
  error: String
}
```

| Field     | Type   | Description                     |
| --------- | ------ | ------------------------------- |
| `message` | String | Echoed message (null on error)  |
| `error`   | String | Error message (null on success) |

## Queries

### echo

Echo back the input message.

```graphql
query {
  echo(message: "hello")
}
```

**Response:**

```json
{
  "data": {
    "echo": "hello"
  }
}
```

**curl:**

```bash
curl -X POST http://localhost:8080/graphql \
  -H "Content-Type: application/json" \
  -d '{"query": "{ echo(message: \"hello\") }"}'
```

### echoWithDelay

Echo with delay for timeout testing.

| Argument  | Type    | Description           |
| --------- | ------- | --------------------- |
| `message` | String! | Message to echo       |
| `delayMs` | Int!    | Delay in milliseconds |

```graphql
query {
  echoWithDelay(message: "hello", delayMs: 5000)
}
```

**curl:**

```bash
curl -X POST http://localhost:8080/graphql \
  -H "Content-Type: application/json" \
  -d '{"query": "{ echoWithDelay(message: \"hello\", delayMs: 5000) }"}'
```

### echoError

Always returns a GraphQL error.

```graphql
query {
  echoError(message: "test")
}
```

**Response:**

```json
{
  "data": {
    "echoError": null
  },
  "errors": [
    {
      "message": "echo error: test",
      "path": ["echoError"]
    }
  ]
}
```

### echoPartialError

Returns partial data with errors. Messages containing "error" will fail.

```graphql
query {
  echoPartialError(messages: ["hello", "error", "world"]) {
    message
    error
  }
}
```

**Response:**

```json
{
  "data": {
    "echoPartialError": [
      { "message": "hello", "error": null },
      { "message": null, "error": "message contains 'error'" },
      { "message": "world", "error": null }
    ]
  }
}
```

### echoWithExtensions

Returns data with custom GraphQL extensions.

```graphql
query {
  echoWithExtensions(message: "hello")
}
```

**Response:**

```json
{
  "data": {
    "echoWithExtensions": "hello"
  },
  "extensions": {
    "timestamp": "2024-01-01T00:00:00Z",
    "requestId": "abc123"
  }
}
```

## Mutations

### createMessage

Create a new message.

```graphql
mutation {
  createMessage(text: "Hello, World!") {
    id
    text
    createdAt
  }
}
```

**Response:**

```json
{
  "data": {
    "createMessage": {
      "id": "1",
      "text": "Hello, World!",
      "createdAt": "2024-01-01T00:00:00Z"
    }
  }
}
```

**curl:**

```bash
curl -X POST http://localhost:8080/graphql \
  -H "Content-Type: application/json" \
  -d '{"query": "mutation { createMessage(text: \"Hello\") { id text createdAt } }"}'
```

### updateMessage

Update an existing message.

```graphql
mutation {
  updateMessage(id: "1", text: "Updated text") {
    id
    text
    createdAt
  }
}
```

**Response:**

```json
{
  "data": {
    "updateMessage": {
      "id": "1",
      "text": "Updated text",
      "createdAt": "2024-01-01T00:00:00Z"
    }
  }
}
```

### deleteMessage

Delete a message. Returns `true` if the message existed.

```graphql
mutation {
  deleteMessage(id: "1")
}
```

**Response:**

```json
{
  "data": {
    "deleteMessage": true
  }
}
```

## Subscriptions

Subscriptions use WebSocket protocol. Connect to `ws://localhost:8080/graphql`.

### messageCreated

Subscribe to new message events. Triggered when `createMessage` is called.

```graphql
subscription {
  messageCreated {
    id
    text
    createdAt
  }
}
```

**Using websocat:**

```bash
echo '{"type":"start","id":"1","payload":{"query":"subscription { messageCreated { id text } }"}}' | \
  websocat ws://localhost:8080/graphql -n
```

### countdown

Subscribe to countdown events. Streams integers from `from` down to 0.

```graphql
subscription {
  countdown(from: 5)
}
```

**Response stream:**

```json
{"data": {"countdown": 5}}
{"data": {"countdown": 4}}
{"data": {"countdown": 3}}
{"data": {"countdown": 2}}
{"data": {"countdown": 1}}
{"data": {"countdown": 0}}
```

## Introspection

GraphQL introspection is enabled. Query the schema:

```graphql
query {
  __schema {
    types {
      name
    }
  }
}
```

**curl:**

```bash
curl -X POST http://localhost:8080/graphql \
  -H "Content-Type: application/json" \
  -d '{"query": "{ __schema { types { name } } }"}'
```

## Error Handling

### Standard Errors

```json
{
  "data": null,
  "errors": [
    {
      "message": "error description",
      "path": ["fieldName"],
      "locations": [{ "line": 1, "column": 3 }]
    }
  ]
}
```

### Partial Errors

GraphQL allows partial success. Some fields may return data while others return errors:

```json
{
  "data": {
    "echo": "hello",
    "echoError": null
  },
  "errors": [
    {
      "message": "echo error: test",
      "path": ["echoError"]
    }
  ]
}
```

## Health Check

```bash
curl http://localhost:8080/health
```

**Response:**

```json
{
  "status": "ok"
}
```
