FROM golang:1.24-alpine

WORKDIR /app

# Copy go.mod and go.sum first for caching
COPY go.mod ./
COPY go.sum ./

# Then copy source code (needed before tidy!)
COPY . .

# Tidy and download dependencies
RUN go mod tidy
RUN go mod download

EXPOSE 3000

# Build the app
RUN go build -o server .


