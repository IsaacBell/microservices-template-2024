# Define the first stage for building the application
FROM golang:1.22.2 AS builder

# Install dependencies
RUN apt-get update && apt-get install -y make protobuf-compiler

# Set the working directory
WORKDIR /src

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the entire project
COPY . .

# Build the application
RUN make compile

# Define the second stage to reduce the final image size
FROM debian:stable-slim

# Install essential packages
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates netbase

# Reduce the image size
RUN rm -rf /var/lib/apt/lists/*

# Set the working directory
WORKDIR /app
RUN touch /app/.env

# Copy the built binary from the builder stage to the root of the app directory
COPY --from=builder /src/bin/core .

# Set environment variables
ENV PORT=8000
ENV METRICS_PORT=9000

# Expose application ports
EXPOSE $PORT
EXPOSE $METRICS_PORT

# Run the application
RUN chmod +x core
CMD ["./core"]
