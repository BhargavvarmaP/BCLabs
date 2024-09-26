# Use the official Golang image as the build stage
FROM golang:1.23 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files for dependency management
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the entire project to the working directory
COPY . .

# Build the Go application
RUN go build -o main .

# Start a new stage from scratch
FROM alpine:latest

# Set the Current Working Directory in the new stage
WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/main .

# Expose the port that the app runs on
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
