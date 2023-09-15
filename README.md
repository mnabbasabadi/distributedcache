# Distributed In-Memory Cache

## Overview

This project implements a simple distributed in-memory cache with a registry server for node discovery.

## Project Structure

```
project-root/
├── cmd/
│   ├── registry/
│   │   └── main.go
│   └── node/
│       └── main.go
├── pkg/
│   ├── cache/
│   │   └── cache.go
│   └── hash/
│       └── hash_ring.go
├── internal/
│   └── registry/
│       └── registry.go
├── api/
│   └── handler.go
└── README.md

```


## Setup

1. Set up your Go environment.
2. Clone this repository.
3. Run the registry server: `go run cmd/registry/main.go`.
4. Run node servers: `go run cmd/node/main.go`.

## API

### Registry Server

- **Register Node**: `POST /register` with JSON body `{ "IpAddr": "127.0.0.1:5000" }`.
- **Deregister Node**: `POST /deregister` with JSON body `{ "IpAddr": "127.0.0.1:5000" }`.
- **Get Nodes**: `GET /nodes` to get a list of all registered nodes.

### Node Server

- **Set Value**: `GET /set?key=<key>&value=<value>`.
- **Get Value**: `GET /get?key=<key>`.
- **Delete Value**: `GET /delete?key=<key>`.

## License

MIT
