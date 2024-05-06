# Use an official Golang runtime as a parent image
FROM golang:1.22.2 as builder

# Set the working directory inside the container
WORKDIR /app

# Copy the local package files to the container's workspace.
ADD . /app

# Build the Go app
RUN make build

# Use a Docker multi-stage build to create a lean production image.
FROM golang:1.22.2

WORKDIR /app

COPY --from=builder /app/feedex /app

# Run the main binary.
CMD ["./feedex"]
