# echo-grpc API Reference

## Service Definition

```protobuf
package echo.v1;

service Echo {
  // Unary RPCs
  rpc Echo (EchoRequest) returns (EchoResponse);
  rpc EchoWithDelay (EchoWithDelayRequest) returns (EchoResponse);
  rpc EchoError (EchoErrorRequest) returns (EchoResponse);

  // Streaming RPCs
  rpc ServerStream (ServerStreamRequest) returns (stream EchoResponse);
  rpc ClientStream (stream EchoRequest) returns (EchoResponse);
  rpc BidirectionalStream (stream EchoRequest) returns (stream EchoResponse);
}
```

## Messages

### EchoRequest

```protobuf
message EchoRequest {
  string message = 1;
}
```

| Field     | Type   | Description          |
| --------- | ------ | -------------------- |
| `message` | string | Message to echo back |

### EchoResponse

```protobuf
message EchoResponse {
  string message = 1;
  map<string, string> metadata = 2;
}
```

| Field      | Type               | Description               |
| ---------- | ------------------ | ------------------------- |
| `message`  | string             | Echoed message            |
| `metadata` | map<string,string> | Request metadata (echoed) |

### EchoWithDelayRequest

```protobuf
message EchoWithDelayRequest {
  string message = 1;
  int32 delay_ms = 2;
}
```

| Field      | Type   | Description           |
| ---------- | ------ | --------------------- |
| `message`  | string | Message to echo back  |
| `delay_ms` | int32  | Delay before response |

### EchoErrorRequest

```protobuf
message EchoErrorRequest {
  string message = 1;
  int32 code = 2;
  string details = 3;
}
```

| Field     | Type   | Description               |
| --------- | ------ | ------------------------- |
| `message` | string | Message (unused in error) |
| `code`    | int32  | gRPC status code (0-16)   |
| `details` | string | Error details message     |

### ServerStreamRequest

```protobuf
message ServerStreamRequest {
  string message = 1;
  int32 count = 2;
  int32 interval_ms = 3;
}
```

| Field         | Type   | Description                      |
| ------------- | ------ | -------------------------------- |
| `message`     | string | Message to echo in each response |
| `count`       | int32  | Number of responses to stream    |
| `interval_ms` | int32  | Interval between responses       |

## RPCs

### Echo (Unary)

Simple echo - returns the input message.

```bash
grpcurl -plaintext -d '{"message": "hello"}' \
  localhost:50051 echo.v1.Echo/Echo
```

**Response:**

```json
{
  "message": "hello",
  "metadata": {
    "content-type": "application/grpc"
  }
}
```

### EchoWithDelay (Unary)

Echo with delay for timeout testing.

```bash
grpcurl -plaintext -d '{"message": "hello", "delay_ms": 5000}' \
  localhost:50051 echo.v1.Echo/EchoWithDelay
```

**Response:** Same as Echo, returned after specified delay.

### EchoError (Unary)

Returns a gRPC error with the specified status code.

| Code | Name                |
| ---- | ------------------- |
| 0    | OK                  |
| 1    | CANCELLED           |
| 2    | UNKNOWN             |
| 3    | INVALID_ARGUMENT    |
| 4    | DEADLINE_EXCEEDED   |
| 5    | NOT_FOUND           |
| 6    | ALREADY_EXISTS      |
| 7    | PERMISSION_DENIED   |
| 8    | RESOURCE_EXHAUSTED  |
| 9    | FAILED_PRECONDITION |
| 10   | ABORTED             |
| 11   | OUT_OF_RANGE        |
| 12   | UNIMPLEMENTED       |
| 13   | INTERNAL            |
| 14   | UNAVAILABLE         |
| 15   | DATA_LOSS           |
| 16   | UNAUTHENTICATED     |

```bash
grpcurl -plaintext -d '{"message": "test", "code": 5, "details": "resource not found"}' \
  localhost:50051 echo.v1.Echo/EchoError
```

**Response:**

```
ERROR:
  Code: NotFound
  Message: resource not found
```

### ServerStream (Server Streaming)

Server sends multiple responses over time.

```bash
grpcurl -plaintext -d '{"message": "ping", "count": 5, "interval_ms": 1000}' \
  localhost:50051 echo.v1.Echo/ServerStream
```

**Response:** Streams `count` responses with `interval_ms` delay between each:

```json
{"message": "ping [1/5]", "metadata": {...}}
{"message": "ping [2/5]", "metadata": {...}}
{"message": "ping [3/5]", "metadata": {...}}
{"message": "ping [4/5]", "metadata": {...}}
{"message": "ping [5/5]", "metadata": {...}}
```

### ClientStream (Client Streaming)

Client sends multiple messages, server responds once with aggregated result.

```bash
echo '{"message": "hello"}
{"message": "world"}
{"message": "!"}' | grpcurl -plaintext -d @ \
  localhost:50051 echo.v1.Echo/ClientStream
```

**Response:**

```json
{
  "message": "Received 3 messages: hello, world, !",
  "metadata": {...}
}
```

### BidirectionalStream (Bidirectional Streaming)

Both client and server stream messages simultaneously. Each client message is echoed back immediately.

```bash
echo '{"message": "one"}
{"message": "two"}
{"message": "three"}' | grpcurl -plaintext -d @ \
  localhost:50051 echo.v1.Echo/BidirectionalStream
```

**Response:**

```json
{"message": "one", "metadata": {...}}
{"message": "two", "metadata": {...}}
{"message": "three", "metadata": {...}}
```

## Server Reflection

The server supports gRPC server reflection for service discovery.

```bash
# List all services
grpcurl -plaintext localhost:50051 list

# Describe service
grpcurl -plaintext localhost:50051 describe echo.v1.Echo

# Describe message
grpcurl -plaintext localhost:50051 describe echo.v1.EchoRequest
```

## Metadata

Request metadata is echoed back in the `metadata` field of every response. Custom metadata can be sent using grpcurl's `-H` flag:

```bash
grpcurl -plaintext \
  -H "x-request-id: abc123" \
  -H "x-custom-header: value" \
  -d '{"message": "hello"}' \
  localhost:50051 echo.v1.Echo/Echo
```

**Response:**

```json
{
  "message": "hello",
  "metadata": {
    "content-type": "application/grpc",
    "x-custom-header": "value",
    "x-request-id": "abc123"
  }
}
```
