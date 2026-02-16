FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /backend-1 .

FROM alpine:3.20
RUN apk add --no-cache curl \
    && curl -fsSL https://github.com/golang-migrate/migrate/releases/download/v4.18.1/migrate.linux-amd64.tar.gz \
       | tar xz -C /usr/local/bin \
    && apk del curl
COPY --from=builder /backend-1 /backend-1
COPY migrations/ /migrations/
EXPOSE 8080
ENTRYPOINT ["/backend-1"]
