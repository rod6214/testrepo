# Start from the latest golang base image
FROM golang:alpine AS build-env
# Add Maintainer Info
LABEL project="Gnalog"
# Set the Current Working Directory inside the container
WORKDIR /app
# Copy go mod and sum files 
COPY go.mod go.sum ./
# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download
# Copy the source from the current directory to the Working Directory inside the container
COPY . .
# Build the Go app
RUN go build -o main .
# Expose port 8080 to the outside world
# EXPOSE 80

# # Final stage
# FROM alpine:3.11
# WORKDIR /app
# COPY --from=build-env /app/main .
# # Expose port 8080 to the outside world
EXPOSE 80
ENTRYPOINT ./main