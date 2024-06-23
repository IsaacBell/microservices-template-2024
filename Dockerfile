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

# Ensure Google Wire is installed and generate necessary files
RUN go install github.com/google/wire/cmd/wire@latest
RUN make init
RUN make config
RUN make api
RUN make wire

# Build the application
RUN make build

# Define the second stage to reduce the final image size
FROM debian:stable-slim

# Install essential packages
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates netbase

# Set the working directory
WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /src/bin /app/bin

# Expose application ports
EXPOSE 8000
EXPOSE 9000

# Run the application
CMD ["./bin/core"]