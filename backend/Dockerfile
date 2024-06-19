# Stage 1: Build the Go application
FROM golang:1.22 as builder

WORKDIR /app

# Cache and install dependencies
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# Copy the source code
COPY backend/ .

# Install Air for hot reload
RUN go install github.com/air-verse/air@latest

# Build the application
RUN go build -o main .

# Stage 2: Run the application
FROM golang:1.22

WORKDIR /app

# Copy the binary and Air config from the builder stage
COPY --from=builder /go/bin/air /usr/bin/air
COPY --from=builder /app /app

# Expose the application port
EXPOSE 8080

# Command to run the application with Air
CMD ["air", "-c", ".air.toml"]
