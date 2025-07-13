PROTO_DIR = proto
PROTO_FILES = $(wildcard $(PROTO_DIR)/*.proto)

proto:
	protoc \
	  --go_out=. \
	  --go_opt=paths=source_relative \
	  --go-grpc_out=. \
	  --go-grpc_opt=paths=source_relative \
	  $(PROTO_FILES)


test-grpc:
	go run ./test/test_grpc_client.go