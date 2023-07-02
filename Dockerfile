# Use the official GoLang base image
FROM golang:1.20-alpine

# Set the working directory inside the container
WORKDIR /app

# Install air for hot-reloading
RUN apk update && apk add --no-cache git

# Copy the GoLang application source code to the container
COPY . .

# Generate go.mod and go.sum files
RUN go mod init nice_stream

# Expose a port (if your application listens on a specific port)
EXPOSE 8080

# Specify the command to run the application when the container starts
CMD ["go", "run", "main.go"]
