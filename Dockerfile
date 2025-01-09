# Stage 1: Build the Go binary
FROM golang:1.22.2-alpine AS builder

# Set working directory in the container
WORKDIR /app

# Install necessary dependencies
RUN apk add --no-cache git

# Copy the source code
COPY . .
WORKDIR /app
# RUN go mod download

# Build the Go binary
# RUN go get
RUN go build -o /app/main ./src

# Stage 2: Create the minimal runtime container
FROM alpine:3.18

# Set up directories for static assets
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/main /app/main

# Copy Assets
COPY --from=builder /app/src/static /app/static
COPY --from=builder /app/src/pages /app/pages
COPY --from=builder /app/src/vacations.json /app/pages/vacations.json

# Ensure the binary is executable
RUN chmod +x /app/main

# Set the default command to run the binary
CMD ["/app/main"]


