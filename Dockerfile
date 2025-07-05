# Build stage
FROM golang:1.24-alpine

WORKDIR /app

# Copy go.mod and go.sum first for better caching
COPY go.mod ./
COPY go.sum ./

# Tidy modules before downloading
RUN go mod tidy
RUN go mod download

# Now copy the rest of the app
COPY . .

# Build the binary
RUN go build -o server .

