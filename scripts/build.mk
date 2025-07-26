lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.59.1
	go fmt ./...
	golangci-lint run --fix