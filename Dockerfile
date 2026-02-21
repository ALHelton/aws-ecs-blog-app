# Step 1: Build the binary
FROM golang:1.26.0-alpine AS builder

# Install git or other necessary dependencies
RUN apk add --no-cache git

# Set the working directory
WORKDIR /aws-ecs-blog-app

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the binary for arm64
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o main cmd/api/main.go

# Step 2: Final lightweight image
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from the builder
COPY --from=builder /aws-ecs-blog-app/main .

# Expose port
EXPOSE 8080

CMD ["./main"]