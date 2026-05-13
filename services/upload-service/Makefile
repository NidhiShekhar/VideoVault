APP_NAME=upload-service

.PHONY: run test tidy

run:
	go run ./cmd/server -config ./config/config.yaml

test:
	go test ./...

tidy:
	go mod tidy
