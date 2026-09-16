# ---- Build stage ----
FROM golang:1.22-alpine AS build
WORKDIR /app

# Download dependencies first so Docker can cache this layer
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /ticket-system .

# ---- Run stage ----
FROM alpine:3.20
RUN adduser -D -H appuser
USER appuser

COPY --from=build /ticket-system /ticket-system

ENV PORT=8080
EXPOSE 8080
CMD ["/ticket-system"]
