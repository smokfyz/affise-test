run:
	@docker build -t server .
	@docker run --rm -p 8080:8080 -e LOG_LEVEL=${LOG_LEVEL} server ./server

test:
	@docker run --rm -v $(PWD):/app -w /app golang:1.22.1-alpine go test ./...

lint:
	@docker run --rm -v $(PWD):/app -w /app golangci/golangci-lint:v1.57.1 golangci-lint run
