# Start with a lightweight Go image
FROM golang:1.23.4 as builder

# Set the working directory inside the container
WORKDIR /app

# Copy Go module files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the entire backend source code
COPY . .

# Build the Go binary
RUN go build -o main .

# Create a lightweight runtime image
FROM gcr.io/distroless/base-debian11

# Copy the built Go binary from the builder stage
COPY --from=builder /app/main /main

# Expose the application port
EXPOSE 8080

# Run the Go application
ENTRYPOINT ["/main"]
