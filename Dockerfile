# Use the official Golang image to build the application
FROM golang:1.22-alpine as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Install git and other dependencies
RUN apk update && apk add --no-cache git

# Copy go mod and go sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go build -o main .

# Use a minimal base image
FROM alpine:latest

# Set the Current Working Directory inside the container
WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .

# Install ca-certificates for HTTPS requests
RUN apk add --no-cache ca-certificates

# Set the entry point to the executable
ENTRYPOINT ["./main"]
