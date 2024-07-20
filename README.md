
# go-sample-rest-api

## Project Description

coming soon...

---
## Prerequisites

- **Go Version**: This repository requires [Go 1.22.0](https://golang.org/dl/) or higher. Make sure you have the correct
  version installed to run the examples provided.

## Dependencies

This project uses several third-party libraries to handle various functionalities:

### Web Routing
- **[Gorilla Mux](https://github.com/gorilla/mux)**  
  Version: 1.8.1  
  Gorilla Mux is a powerful HTTP router and URL matcher for building Go web servers. It helps in routing incoming requests to their respective handlers.

### Environment Variables
- **[godotenv](https://github.com/joho/godotenv)**  
  Version: 1.5.1  
  godotenv is a Go library for loading environment variables from a `.env` file into the system environment, facilitating easier configuration management in development and production environments.

### Database Connectivity
- **[lib/pq](https://github.com/lib/pq)**  
  Version: 1.10.9  
  lib/pq is a pure Go Postgres driver for the database/sql package. It is designed to work with PostgreSQL and allows Go applications to interact with PostgreSQL databases using native Go APIs.

### Data Validation
- **[Go Playground Validator](https://github.com/go-playground/validator)**  
  Version: v10  
  Go Playground Validator is a library for performing validations on structs and individual fields based on tags. It is highly useful for validating input data in RESTful APIs.

### Authentication
- **[Go JWT](https://github.com/golang-jwt/jwt)**  
  Version: v5  
  Go JWT is a Go library that provides tools for creating and verifying JSON Web Tokens, which are an open standard used for securely transmitting information between parties as a JSON object.

### Database Migration Tool
- **[golang-migrate](https://github.com/golang-migrate/migrate)**  
  This tool handles database migrations. It is essential for managing database schema changes and is used to apply and revert database migrations.

### SQL Mocking for Testing
- **[go-sqlmock](https://github.com/DATA-DOG/go-sqlmock)**
 This is a mock library for SQL database operations that allows you to test your application's data access logic without needing to interact with a real database, thus speeding up tests and avoiding side effects.

### Testify
- **[Testify](https://github.com/stretchr/testify)** 

This is a toolkit with common assertions and mocks that are used in testing Go code. It provides a friendly and comprehensive set of tools that enhance the Go testing experience.

### UUID
- **[UUID](https://github.com/google/uuid)** 

This library is used in our project to generate universally unique identifiers (UUIDs). UUIDs are essential for creating unique keys in our database and ensuring that our data can be uniquely identified across distributed systems.
### Installation of Dependencies
To install these Go dependencies, run the following command:

```bash
go get -u github.com/gorilla/mux github.com/joho/godotenv github.com/lib/pq github.com/go-playground/validator/v10 github.com/golang-jwt/jwt/v5 github.com/DATA-DOG/go-sqlmock github.com/stretchr/testify github.com/google/uuid

```
### Installation of golang-migrate
golang-migrate can be installed using Homebrew on macOS:

```bash
brew install golang-migrate
```
For other platforms, follow the installation instructions on the golang-migrate GitHub page.

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
  To create a new migration file, specify the migration name:
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


### Adding Debug Instructions to the README

After the existing sections, you can add a new section that explains how to start the application in debug mode. This will guide users on how to use the newly added Makefile target `debug`.

```markdown
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
```
