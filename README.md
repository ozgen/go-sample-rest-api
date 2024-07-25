
# go-sample-rest-api

## Project Description

This repository demonstrates how to implement a RESTful API in Go with integration to Azure Blob Storage for handling image uploads and downloads. The API follows best practices and utilizes popular Go libraries to manage routing, environment variables, database interactions, and more.

## Prerequisites

- **Go Version**: This repository requires [Go 1.22.0](https://golang.org/dl/) or higher. Ensure you have the correct version installed to run the provided examples.

## Dependencies

This project uses several third-party libraries to handle various functionalities:

### Web Routing
- **[Gorilla Mux](https://github.com/gorilla/mux)**  
  Version: 1.8.1  
  Gorilla Mux is a powerful HTTP router and URL matcher for building Go web servers.

### Environment Variables
- **[godotenv](https://github.com/joho/godotenv)**  
  Version: 1.5.1  
  godotenv loads environment variables from a `.env` file into the system environment.

### Database Connectivity
- **[lib/pq](https://github.com/lib/pq)**  
  Version: 1.10.9  
  lib/pq is a pure Go Postgres driver.

### Data Validation
- **[Go Playground Validator](https://github.com/go-playground/validator)**  
  Version: v10  
  Validates structs and individual fields based on tags.

### Authentication
- **[Go JWT](https://github.com/golang-jwt/jwt)**  
  Version: v5  
  Provides methods for creating and verifying JSON Web Tokens.

### Database Migration Tool
- **[golang-migrate](https://github.com/golang-migrate/migrate)**  
  Handles database migrations.

### SQL Mocking for Testing
- **[go-sqlmock](https://github.com/DATA-DOG/go-sqlmock)**
  A library for mocking SQL operations in tests.

### Testify
- **[Testify](https://github.com/stretchr/testify)**  
  A toolkit with common assertions and mocks for testing.

### UUID
- **[UUID](https://github.com/google/uuid)**  
  Used to generate universally unique identifiers.

### Logging
- **[Logrus](https://github.com/sirupsen/logrus)**  
  Version: Latest  
  A structured logger for Go, API compatible with the standard library logger.

### Azure Blob Storage
- **[Azure Storage Blob Go](https://github.com/Azure/azure-storage-blob-go)**  
  Azure SDK for Go that allows managing blob storage operations.

## Installation of Dependencies

To install these dependencies, run:

```bash
go get -u github.com/gorilla/mux github.com/joho/godotenv github.com/lib/pq github.com/go-playground/validator/v10 github.com/golang-jwt/jwt/v5 github.com/DATA-DOG/go-sqlmock github.com/stretchr/testify github.com/google/uuid github.com/sirupsen/logrus github.com/Azure/azure-storage-blob-go/azblob
```

### Installation of golang-migrate

golang-migrate can be installed using Homebrew on macOS:

```bash
brew install golang-migrate
```

For other platforms, follow the installation instructions on the [golang-migrate GitHub page](https://github.com/golang-migrate/migrate).

## Running Commands Using Makefile

The project includes a Makefile for simplifying common tasks:

- **Build the Application**
  ```bash
  make build
  ```

- **Run Tests**
  ```bash
  make test
  ```

- **Run the Application**
  ```bash
  make run
  ```

- **Create a Database Migration**
  ```bash
  make migration name=create_user_table
  ```

- **Apply Migrations**
  ```bash
  make migrate-up
  ```

- **Revert Migrations**
  ```bash
  make migrate-down
  ```

- **Generate Swagger Documentation**
  ```bash
  make swagger
  ```

## Debugging the Application

### Setting Up for Debugging
To debug the application, you need to have Delve Debugger installed. Delve provides a more sophisticated debugging experience for Go applications than traditional print statement debugging.

- **Install Delve:**
  If Delve is not already installed on your machine, you can install it using the following command:

  ```bash
  go install github.com/go-delve/delve/cmd/dlv@latest
  ```

### Start Debugging
Once Delve is installed, you can start the application in debug mode using the following Makefile command. This command sets up the application to run in headless mode, allowing you to connect with an IDE or a remote debugger interface.

- **Debug the Application**
  To start the application in debug mode, run:

  ```bash
  make debug
  ```

This command starts the application with Delve in headless mode, listening for debugger connections on port `2345`. You can attach to this session using your preferred IDE configured for remote debugging.

### Connecting with IntelliJ IDEA
If you are using IntelliJ IDEA or similar IDEs for Go development, follow these steps to connect to the debug session:

1. **Open Run/Debug Configurations:**
  - Navigate to 'Run' -> 'Edit Configurations...'
  - Add a new 'Go Remote' configuration.
2. **Configure the Debugger:**
  - Set the host to `localhost` and the port to `2345`.
3. **Start Debugging:**
  - Select the created debug configuration and click on the debug icon to start debugging.

Ensure that your firewall settings allow traffic on the port if debugging remotely.

### Tips for Effective Debugging
- **Breakpoints:** Set breakpoints in your code where you want the execution to pause.
- **Evaluate Expressions:** Use the debugger to evaluate expressions and inspect the state of your application.
- **Step Execution:** Utilize step over, step into, and step out functionalities to navigate through your code.

This setup provides a comprehensive debugging environment that aids in developing robust applications by allowing developers to find and fix issues efficiently.

---

## Deployment Instructions

This section provides detailed instructions on how to deploy the Go application and PostgreSQL database using Kubernetes and Kustomize.

### Prerequisites

- **Kubernetes Cluster**: Ensure you have access to a Kubernetes cluster.
- **kubectl**: You need to have `kubectl` installed and configured to interact with your Kubernetes cluster.
- **Kustomize**: Ensure that `kustomize` is installed on your machine to customize Kubernetes configurations.

### Setup Kubernetes Namespace

Before deploying any resources, create the `ozgen` namespace to isolate resources in your Kubernetes cluster:

```bash
kubectl create namespace ozgen
```

This command creates a new namespace named `ozgen` where all your Kubernetes resources will be deployed.

### Deployment Process

#### Deploy PostgreSQL

1. **Deploy PostgreSQL Resources**:
   Navigate to the project directory and deploy the PostgreSQL resources within the `ozgen` namespace. This step should be done first, as the application depends on the database being available.

   ```bash
   kustomize build k8s/postgres | kubectl apply -f - -n ozgen
   ```

2. **Set Up Port Forwarding**:
   After deploying PostgreSQL, set up port forwarding to access PostgreSQL from your local machine for initial database setup and migrations.

   ```bash
   kubectl port-forward svc/my-postgres 5432:5432 -n ozgen
   ```

   Replace `svc/my-postgres` with the actual service name of your PostgreSQL instance.

3. **Initialize the Database**:
   With port forwarding in place, proceed to create the database schema and run migrations.

   ```bash
   make migration up
   ```

   Ensure that your `Makefile` includes the necessary commands for handling database migrations.

#### Deploy Go Application

1. **Deploy Application**:
   Once the database is ready, deploy the Go application into the `ozgen` namespace.

   ```bash
   kustomize build k8s/api | kubectl apply -f - -n ozgen
   ```

2. **Verify Deployment**:
   Confirm that the application is running smoothly by checking the status of the deployed pods.

   ```bash
   kubectl get pods -n ozgen
   ```

### Additional Notes

- **Configuration**: Check that the application's configuration, especially the database connection details, are correctly set in `api/secret.yaml`.
- **Security**: Always ensure your Kubernetes configurations and secrets are secure, particularly when managing sensitive data like database credentials.

### Troubleshooting

- If you run into any deployment issues, review the logs of your Kubernetes pods for clues:

  ```bash
  kubectl logs <pod-name> -n ozgen
  ```

- For database connectivity problems, verify that port forwarding is correctly established and that the database credentials are properly configured.


---

## API Documentation

### Access Swagger UI

Our API documentation is available via Swagger UI, which allows you to interact with the API's endpoints directly through your browser. To view the interactive documentation, visit:

[Swagger API Documentation](http://localhost:8080/api/v1/documentation/index.html)

### Integration Details

The documentation is automatically generated and served using the `swaggo/http-swagger` middleware integrated into our Go application. This setup ensures that the API documentation is always up-to-date with the latest codebase.

```go
import (
    "github.com/gorilla/mux"
    "github.com/swaggo/http-swagger"       // HTTP Swagger middleware
    "github.com/swaggo/http-swagger/swaggerFiles" // Swagger embed files
)

func setupRouter() *mux.Router {
    r := mux.NewRouter()
	r.PathPrefix("/documentation/").Handler(httpSwagger.WrapHandler)
    return r
}
```

Please make sure to update the URL appropriately if your application is deployed to a different environment.

---
