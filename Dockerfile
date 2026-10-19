# ---------- Build stage ----------

FROM golang:1.25.6-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/gophertalk ./cmd


# ---------- Runtime stage ----------

FROM alpine:latest

WORKDIR /app

# Сертификаты необходимы для HTTPS-соединений.
RUN apk --no-cache add ca-certificates \
    && addgroup -S app && adduser -S app -G app

COPY --from=builder --chown=app:app /app/gophertalk ./gophertalk

USER app

EXPOSE 8080

CMD ["./gophertalk"]