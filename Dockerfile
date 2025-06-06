# Stage 1: Build environment
FROM golang:1-alpine AS build-env

# Set application name and root directory for the build stage
ENV APP_ROOT=/app
WORKDIR ${APP_ROOT}

# Copy Go module files first to leverage Docker layer caching for dependencies
COPY go.mod go.sum ./
RUN go mod download
RUN go mod tidy # Ensure dependencies are clean and vendor if necessary

# Copy the entire application source code
COPY . .

# Build the Go application
# - CGO_ENABLED=0: Create a static binary, important for alpine runtime and distroless
# - ldflags="-s -w": Strip debug symbols and DWARF information to reduce binary size
# Output the binary to /${APP_NAME} in this build stage (root of the build stage)
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /${APP_NAME} ./cmd/parallax

# Stage 2: Runtime environment
# Use a distroless image for a minimal and secure runtime.
# static-debian11 is suitable for CGO_ENABLED=0 Go binaries.
# FROM debian AS runtime
FROM gcr.io/distroless/static-debian12 AS runtime
#
# RUN apt update && apt install -y ca-certificates && apt install curl -y && rm -rf /var/lib/apt/lists/* && apt clean

# These ENV vars might be overridden by .env or could be removed if .env is the sole source of truth
ENV APP_NAME=parallax
ENV ENV=production
ENV GIN_MODE=release
WORKDIR /app

# Copy .env file for runtime configuration
# IMPORTANT: Ensure .env does not contain sensitive data if this image is pushed to a public registry.
# Consider using Docker secrets or environment variables passed at runtime for sensitive data.
# COPY .env.example /app/.env

# Copy the compiled application binary from the build stage
COPY --from=build-env /${APP_NAME} /app/${APP_NAME}

# Expose the port the application listens on
# This should match the port configured in the application (e.g., via .env or default)
EXPOSE 80

# Define the command to run the application
# Using the application binary copied into /app/
ENTRYPOINT ["/app/parallax"]
