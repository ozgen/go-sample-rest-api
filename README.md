
# go-sample-rest-api

## Project Overview

This repository demonstrates the implementation of a RESTful API in Go with integration to Azure Blob Storage for handling image uploads and downloads. It also includes Prometheus metrics to monitor the API performance and health. The API leverages popular Go libraries for web routing, environment management, database interactions, and more, following best practices to provide a robust backend service.

## Prerequisites

- **Go Version**: Requires [Go 1.23.0](https://golang.org/dl/) or higher.

## Dependencies

Here are the primary libraries and tools used:

- **[Gorilla Mux](https://github.com/gorilla/mux)** - For HTTP routing.
- **[godotenv](https://github.com/joho/godotenv)** - For environment variable management.
- **[lib/pq](https://github.com/lib/pq)** - PostgreSQL driver.
- **[Go Playground Validator](https://github.com/go-playground/validator)** - For data validation.
- **[Go JWT](https://github.com/golang-jwt/jwt)** - For authentication.
- **[golang-migrate](https://github.com/golang-migrate/migrate)** - For database migrations.
- **[Azure Storage Blob Go](https://github.com/Azure/azure-storage-blob-go)** - For managing Azure Blob Storage.
- **[Prometheus Go Client](https://github.com/prometheus/client_golang)** - For exposing custom metrics collected from the API.

## Installation

Install all dependencies at once:

```bash
go get -u github.com/gorilla/mux github.com/joho/godotenv github.com/lib/pq github.com/go-playground/validator/v10 github.com/golang-jwt/jwt/v5 github.com/DATA-DOG/go-sqlmock github.com/stretchr/testify github.com/google/uuid github.com/sirupsen/logrus github.com/Azure/azure-storage-blob-go/azblob github.com/prometheus/client_golang/prometheus github.com/prometheus/client_golang/prometheus/promhttp
```

### Database Migration Tool

Install golang-migrate on macOS using Homebrew:

```bash
brew install golang-migrate
```

## Usage with Makefile

Common tasks can be executed using the Makefile:

```bash
make build        # Build the application
make test         # Run tests
make run          # Run the application
make migration    # Create a database migration
make migrate-up   # Apply migrations
make migrate-down # Revert migrations
make swagger      # Generate Swagger documentation
```

## Debugging

Set up Delve for sophisticated debugging:

```bash
go install github.com/go-delve/delve/cmd/dlv@latest
make debug        # Start the application in debug mode
```

### API and Metrics Documentation

Swagger UI is accessible at:
[Swagger API Documentation](http://localhost:8080/api/v1/documentation/index.html)

Swagger is served using `swaggo/http-swagger` integrated into the Gorilla Mux setup.

Prometheus metrics are exposed at:
[Prometheus Metrics Endpoint](http://localhost:8080/metrics)

This endpoint is used by Prometheus to collect metrics about the application's performance and health, leveraging the Prometheus Go client.

## Deployment Instructions

Deploy using Kubernetes and Kustomize:

1. **Create Namespace**:
    ```bash
    kubectl create namespace ozgen
    kubectl create namespace logging
    kubectl create namespace monitoring
    ```

2. **Deploy Resources**:
    ```bash
    kustomize build k8s | kubectl apply -f -
    ```

3. **Initialize Database**:
    ```bash
    kustomize build k8s/ozgen/postgres | kubectl apply -f - -n ozgen
    kubectl port-forward svc/my-postgres 5432:5432 -n ozgen
    make migration up
    ```

Check pod status and logs using:
```bash
kubectl get pods -n ozgen
kubectl logs <pod-name> -n ozgen
```

---
