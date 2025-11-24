# --- Builder Stage ---
# Use an official Go image as the builder
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./
# Download dependencies
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the Go application
# CGO_ENABLED=0 disables Cgo, creating a static binary
# -o /go-htmx-todo specifies the output file
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /go-htmx-todo ./src/main

# --- Final Stage ---
# Use a minimal alpine image
FROM alpine:latest

WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /go-htmx-todo .

# Copy templates
COPY ./src/resources/ ./src/resources/

# Create the uploads directory
# Note: Data in here will be lost unless a volume is mounted
RUN mkdir -p /app/uploads

# Expose the port the app runs on
EXPOSE 8080

# The command to run the application
CMD ["./go-htmx-todo"]