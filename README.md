[![Go Report Card](https://goreportcard.com/badge/github.com/nccapo/kv-store)](https://goreportcard.com/report/github.com/nccapo/kv-store)
![Go Version](https://img.shields.io/github/go-mod/go-version/nccapo/kv-store)
[![Go Reference](https://pkg.go.dev/badge/github.com/nccapo/kv-store.svg)](https://pkg.go.dev/github.com/nccapo/kv-store)
![License](https://img.shields.io/github/license/nccapo/kv-store)
![Issues Welcome](https://img.shields.io/badge/contributions-welcome-brightgreen.svg?style=flat)
![Version](https://img.shields.io/github/v/release/nccapo/kv-store)

# Distributed Key-Value Store

A distributed key-value store implementation inspired by etcd, built with Go. This project provides a simple, secure, fast, and reliable distributed key-value store with the following features:

- Distributed consensus using the Raft protocol
- HTTP API for key-value operations
- TTL (Time-To-Live) support for keys
- Thread-safe operations
- Graceful shutdown

## Features

- **Distributed Consensus**: Uses the Raft protocol to ensure consistency across all nodes in the cluster
- **HTTP API**: RESTful API for key-value operations
- **TTL Support**: Keys can be set with an expiration time
- **Thread Safety**: All operations are thread-safe using Go's sync package
- **Persistence**: Log entries are persisted to disk
- **Graceful Shutdown**: Clean shutdown of all components

## Getting Started

### Prerequisites

- Go 1.21 or later

### Building

```bash
go build
```

### Running

To start a node in the cluster:

```bash
./kv-store -id node1 -addr :8080 -peers "node2=localhost:8081,node3=localhost:8082"
```

### API Endpoints

- `GET /v1/kv/{key}` - Get a value by key
- `PUT /v1/kv/{key}` - Set a value with optional TTL
- `DELETE /v1/kv/{key}` - Delete a key-value pair

Example request to set a value with TTL:

```bash
curl -X PUT http://localhost:8080/v1/kv/mykey \
  -H "Content-Type: application/json" \
  -d '{"value": "myvalue", "ttl": "5s"}'
```

## Architecture

The project is organized into the following components:

- `internal/raft`: Raft consensus protocol implementation
- `internal/api`: HTTP API server
- `store`: Key-value store implementation

### Raft Protocol

The Raft protocol ensures that all nodes in the cluster maintain a consistent state. The implementation includes:

- Leader election
- Log replication
- Membership changes
- Snapshotting

### Key-Value Store

The key-value store provides:

- Thread-safe operations using sync.Map
- TTL support for keys
- Persistence of data
- Clean error handling

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
