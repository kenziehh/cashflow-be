FROM golang:1.25-alpine

# Set working directory
WORKDIR /app

# Install dependencies
RUN apk add --no-cache git curl bash

# Install Air for hot reload
RUN go install github.com/air-verse/air@latest

# Copy go mod files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire source
COPY . .

# Expose the app port
EXPOSE 8080

# Command for Air (auto reload)
CMD ["air", "-c", ".air.toml"]
