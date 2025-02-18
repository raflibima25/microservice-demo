.PHONY: proto

# Install required tools
install-tools:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

# Generate proto files
proto:
	protoc --proto_path=grpc/proto \
		--go_out=grpc/pb --go_opt=paths=source_relative \
		--go-grpc_out=grpc/pb --go-grpc_opt=paths=source_relative \
		grpc/proto/auth/*.proto \
		grpc/proto/product/*.proto

# Generate protoset file for Postman
protoset:
	protoc --proto_path=grpc/proto \
		--descriptor_set_out=grpc/proto/api.protoset \
		--include_imports \
		grpc/proto/auth/*.proto \
		grpc/proto/product/*.proto

# Run the server
run-auth:
	cd auth-service && go run cmd/main.go

run-product:
	cd product-service && go run cmd/main.go

run-frontend:
	cd frontend && npm run start

run-api-gateway:
	cd api-gateway && go run main.go