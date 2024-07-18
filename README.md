
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

