# Distributed In-Memory Cache

## Overview

This project implements a simple distributed in-memory cache with a registry server for node discovery.
- The registry server is responsible for registering and deregistering nodes.
- The node server is responsible for storing and retrieving key-value pairs.
- The registry server uses a hash ring to distribute the nodes.
- The node server uses a hash ring to distribute the key-value pairs.


## Design

Check the design documentation: [docs](./docs)

## Project Structure

```
master 
├── api
│   ├── v1
│   │   ├── master.openapi.yaml  # OpenAPI 3.0.3 specification
│   │   ├── gen.go                # Script to generate the API code
│   │   ├── service.config.yaml   # Configuration file for the API
│   │   └── api.gen.go            # Generated API code
├── service
│   ├── cmd
│   │   └── main.go               # Entrypoint of the service
│   ├── pkg                        # exportable packages that can be used by other services and defines dependencies
│   │   ├── app
│   │   │   └── app.go
│   ├── tests
│   │   ├── integration
│   │   │   ├── integration.go
│   │   │   └── e2e_test.go
│   │   ├── support
│   │   │   └── client
│   │   │       └── client.go
│   ├── internal                   # internal packages that are not exported and defines dependencies
│   │   ├── api                    # API layer
│   │   │   └── http
│   │   │       └── handler.go
│   │   ├── hash                  
│   │   │   └── hash_ring.go      # Hash ring implementation

node 
├── api
│   ├── the same as master
├── service
│   ├── cmd     # the same as master
│   ├── pkg     # the same as master
│   ├── tests  # the same as master
│   ├── internal                  
│   │   ├── api   # the same as master
│   │   ├── node                  
│   │   │   └── node.go      # Node implementation

```

## API

### Registry Server
 
link to the openapi specification: [register.openapi3.yaml](master/api/v1/register.openapi3.yaml)

### Node Server

link to the openapi specification: [node.openapi3.yaml](node/api/v1/cache.openapi3.yaml)


## Testing

### Unit Tests
```sh
make test
```

### Integration Tests
```sh
make test-integration
```

# Setup

## Prerequisites

- [Docker](https://docs.docker.com/get-docker/)
- [Kubectl](https://kubernetes.io/docs/tasks/tools/)
- [Minikube](https://minikube.sigs.k8s.io/docs/start/)
- [Skaffold](https://skaffold.dev/docs/install/) (optional)

### Step 1: Start Minikube

Start your Minikube cluster with the following command:

```sh
minikube start
```

### Step 2: Build the Docker Images

```sh
docker build -t master-image -f Dockerfile_master .
docker build -t node-image -f Dockerfile_node .
```


### Step 3: Apply Kubernetes Manifests

```sh
kubectl apply -f k8s/
```

### Step 4: Access the Master Service

```sh
minikube service master-service
```
take note of the IP and port of the master service

```sh
curl -X POST -H "Content-Type: application/json" -d '{"key":"value3", "value":"value3"}' http://[address]/keys 
curl -H "Content-Type: application/json" http://[address]/keys/value3   
```

# Skafold

If you have Skaffold installed, you can streamline the development process using:

```sh
skaffold dev
```

This will watch your source files and automatically rebuild and redeploy your application when changes are detected.


# TODO
- [ ] Add more tests
- [ ] Add more documentation
- [ ] Add metrics 
- [ ] Add tracing
- [ ] Add circuit breaker
- [ ] Add clustering for the master service(leader election)
- [ ] Add clustering for the node service
- [ ] master service should be able to handle node failures
- [ ] Add authentication and authorization
- [ ] Add TLS
- [ ] Add more configuration options


## License

MIT
