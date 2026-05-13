# ---- Build stage ----
FROM golang:1.22 AS build

WORKDIR /app

# Copy go.mod (go.sum may not exist yet)
COPY go.mod ./

# Download dependencies (this will generate go.sum if missing)
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the ratelimiter binary
RUN go build -o ratelimiter ./cmd/ratelimiter

# ---- Runtime stage ----
FROM debian:stable-slim

WORKDIR /app

COPY --from=build /app/ratelimiter .

CMD ["./ratelimiter"]
