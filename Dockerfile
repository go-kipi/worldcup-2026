# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy dependency files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o worldcup-backend .

# Final stage
FROM alpine:latest

WORKDIR /app

# Install certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Copy binary from builder
COPY --from=builder /app/worldcup-backend .

# Copy front-end dist folder
# Note: Ensure the 'dist' folder exists in the build context
COPY dist ./dist

# Expose port
EXPOSE 8080

# Run the backend
CMD ["./worldcup-backend"]
