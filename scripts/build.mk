lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go fmt ./...
	golangci-lint run
