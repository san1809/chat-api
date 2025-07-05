# syntax=docker/dockerfile:1

FROM golang:1.21-alpine

# Install Git (needed for go mod)
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Copy the rest of the code
COPY . .

# Build the app
RUN go build -o server .

# Expose port 8080 (your Go app port)
EXPOSE 8080

# Run the app
CMD ["./server"]
