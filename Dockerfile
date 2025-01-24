# First stage: Build stage
FROM golang:1.23 AS builder

WORKDIR /app

# Copy source code to the working directory
COPY . .

# Build main.go
RUN go mod tidy && \
 go mod vendor && \
 go build -o main cmd/main.go

# Second stage: Run stage
FROM ubuntu:24.04

WORKDIR /root/

# Install ffmpeg
RUN apt-get update && apt-get install -y ffmpeg

# Copy the build artifact from the build stage
COPY --from=builder /app/main /root/main
COPY --from=builder /app/configs/default.yml /app/configs/default.yml

# Ensure the binary has executable permissions
RUN chmod +x /root/main

# Set ENTRYPOINT
ENTRYPOINT ["/root/main", "restful", "-c", "/app/configs/default.yml"]
