# This is a multi-stage Dockerfile and requires >= Docker 17.05
# https://docs.docker.com/engine/userguide/eng-image/multistage-build/
FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache gcc musl-dev

# Install buffalo
RUN go install github.com/gobuffalo/cli/cmd/buffalo@latest

WORKDIR /app

# Copy the entire project
COPY . .

# Download dependencies
RUN go mod download

# Build the Buffalo application
# -k flag skips asset compilation
RUN buffalo build -k -o bin/app

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy the binary and required files
COPY --from=builder /app/bin/app .
COPY --from=builder /app/database.yml .
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/public ./public
COPY --from=builder /app/locales ./locales

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata mysql-client

# Set environment variables
ENV GO_ENV=production
ENV PORT=3000

# Expose the application port
EXPOSE 3000

# Run migrations and start the app
CMD ["./app"]
