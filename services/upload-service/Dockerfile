FROM golang:1.22.2-alpine AS builder

WORKDIR /app
COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/upload-service ./cmd/server

FROM gcr.io/distroless/base-debian12
COPY --from=builder /bin/upload-service /bin/upload-service
COPY --from=builder /app/config /config

EXPOSE 8080
ENTRYPOINT ["/bin/upload-service", "-config", "/config/config.yaml"]
