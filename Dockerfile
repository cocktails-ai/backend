# Use the official Go image as the base image
FROM golang:1.17-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files to the working directory
COPY go.mod go.sum ./

# Download the dependencies
RUN go mod download

# Copy the rest of the source code to the working directory
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Start a new stage with a minimal image to reduce the final image size
FROM alpine:3.14

# Set the working directory inside the container
WORKDIR /app

# Copy the compiled Go binary from the previous stage
COPY --from=0 /app/main /app/

# Expose the port that the server will run on
EXPOSE 8080

# Run the server
CMD ["/app/main"]
