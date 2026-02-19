# Build stage
FROM golang:1.24-alpine AS builder

# Set the working directory
WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy dependency files
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project source
COPY . .

# Build the Go application
# CGO_ENABLED=0 ensures a statically linked binary
# GOOS=linux targets the Linux operating system
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Final stage (Runtime)
FROM alpine:latest

# Set the working directory
WORKDIR /app

# Install runtime dependencies
# ca-certificates: for making HTTPS requests
# tzdata: for timezone support
RUN apk add --no-cache ca-certificates tzdata

# Set timezone to Asia/Kuala_Lumpur (matches database connection logic)
ENV TZ=Asia/Kuala_Lumpur

# Copy the pre-built binary from the build stage
COPY --from=builder /app/main .

# Note: The application uses viper to load .env.
# In a production Docker environment, it's recommended to pass variables 
# via environment variables rather than a .env file.
# If you still need a .env file, uncomment the line below:
# COPY .env .

# Expose the application port (matching the default in main.go)
EXPOSE 6060

# Command to run the application
ENTRYPOINT ["./main"]
