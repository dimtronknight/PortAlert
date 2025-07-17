# Build stage
FROM golang:1.22.5-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o notify cmd/notify.go
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o database cmd/database.go

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/

# Copy binaries from builder
COPY --from=builder /app/notify .
COPY --from=builder /app/database .

# Copy config files
COPY --from=builder /app/config ./config/

# Create script for daily run
RUN echo '#!/bin/sh' > daily-run.sh && \
    echo 'echo "🚀 Starting daily portfolio tracking..."' >> daily-run.sh && \
    echo './notify now' >> daily-run.sh && \
    echo 'echo "💾 Saving to database..."' >> daily-run.sh && \
    echo './database save' >> daily-run.sh && \
    echo 'echo "✅ Daily tracking completed!"' >> daily-run.sh && \
    chmod +x daily-run.sh

CMD ["./daily-run.sh"]