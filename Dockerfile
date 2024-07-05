# Build stage
FROM golang:1.19-alpine AS build

# Set the Current Working Directory inside the container
WORKDIR /app

# Install necessary build tools
RUN apk add --no-cache gcc g++ make

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go build -o server .

# Final stage
FROM alpine:latest

# Set the Current Working Directory inside the container
WORKDIR /root/

# Install GStreamer and necessary plugins
RUN apk add --no-cache \
    gstreamer \
    gstreamer-dev \
    gst-plugins-base \
    gst-plugins-good \
    gst-plugins-bad \
    gst-plugins-ugly \
    gst-libav \
    gstreamer-tools

# Copy the Pre-built binary file from the previous stage
COPY --from=build /app/server .
COPY --from=build /app/static ./static

# 18080 - 백엔드
# 15000 - udp port
EXPOSE 18080
EXPOSE 15000

# Command to run the executable
CMD ["./server"]