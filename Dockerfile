# Build stage
FROM golang:1.24-alpine AS base

# Set working directory
WORKDIR /app

FROM base AS builder
# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

RUN go build -o main cmd/app/main.go

RUN go build -o db cmd/db/main.go

# Final stage
FROM base AS final

COPY --from=builder /app/main /app/db .
COPY ./cmd/db/migration ./cmd/db/migration
COPY run-app.sh .

RUN chmod +x run-app.sh

CMD ["./run-app.sh"]

