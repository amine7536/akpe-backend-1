FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /backend-1 .

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=builder /backend-1 /backend-1
EXPOSE 8080
ENTRYPOINT ["/backend-1"]
