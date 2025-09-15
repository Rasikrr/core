lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go fmt ./...
	golangci-lint run


build:
	go build -o bin/app ./cmd/app


deps:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Генерация Go кода из proto файлов
generate:
	mkdir -p ./pkg/api/proto
	@echo "Генерация Go кода из proto файлов..."
	protoc \
		-I=. \
		-I=api/proto \
		--go_out=./pkg \
		--go_opt=paths=source_relative \
		--go-grpc_out=./pkg \
		--go-grpc_opt=paths=source_relative \
		api/proto/orders/orders.proto
	@echo "Генерация завершена!"