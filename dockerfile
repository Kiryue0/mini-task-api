# ---- Builder stage ----
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o mini-task-api .

# ---- Final stage ----
FROM alpine:latest

RUN apk update && apk upgrade --no-cache

RUN adduser -D -u 10001 appuser

WORKDIR /app

COPY --from=builder /app/mini-task-api .

USER appuser

EXPOSE 8080

ENTRYPOINT ["./mini-task-api"]