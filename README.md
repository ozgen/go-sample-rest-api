
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
