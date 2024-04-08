# Build stage
FROM golang:1.22.1 as builder
# Define build env
ENV GOOS linux
ENV CGO_ENABLED 0
# Add a work directory
WORKDIR /app
# Cache and install dependencies
COPY go.mod go.sum ./
RUN go mod download
# Copy app files
COPY . .
# Build app
RUN go build -o main main.go

# Run stage
FROM alpine:3.14 as production
WORKDIR /app
# Add certificates
RUN apk add --no-cache ca-certificates
# Copy built binary from builder
COPY --from=builder /app/main .

COPY config.json .
COPY db/migration ./db/migration

# Expose port
EXPOSE 4000
# Exec built binary
CMD [ "/app/main" ]