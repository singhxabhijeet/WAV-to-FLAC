# Start with a base Go image
FROM golang:alpine AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod .
COPY go.sum .

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app
RUN GOOS=linux GOARCH=amd64 go build -o main .

# Start a new stage from scratch
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .

# Command to run the executable
CMD ["./main"]
