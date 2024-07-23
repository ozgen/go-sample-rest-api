# Use an official Go runtime as a parent image as the builder stage
FROM golang:1.22.0 as builder

# Set the working directory inside the container for all subsequent operations
WORKDIR /app

# Copy go.mod and go.sum for dependency installations
COPY go.mod go.sum ./

# Download and cache the dependencies
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the application binary.
# Disable CGO to create a statically linked binary to run in a minimal container.
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o myapp .

# Use a minimal Alpine image to run the application
FROM alpine:3.14

# Set the working directory in this new stage
WORKDIR /root/

# Copy the compiled binary from the previous stage into this lighter image
COPY --from=builder /app/myapp .

# Expose the port the app runs on
EXPOSE 8080

# Command to execute the binary
CMD ["./myapp"]
