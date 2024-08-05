# Use the official Golang image to create a build artifact.
FROM golang:1.20 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the Go app
RUN go build -o auth-service .

# Start a new stage from scratch
FROM golang:1.20

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/auth-service .

# Expose port 8000 to the outside world
EXPOSE 8000

# Command to run the executable
CMD ["./auth-service"]
