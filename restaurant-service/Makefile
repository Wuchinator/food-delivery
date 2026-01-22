.PHONY: gen-proto docker-build

GOBIN := $(shell go env GOPATH)/bin

export PATH := $(GOBIN):$(PATH)

PROTO_DIR = api/proto/v1
PROTO_OUT_DIR = pkg/restaurant_v1
PLATFORM= linux/amd64

gen-proto:
	@mkdir -p $(PROTO_OUT_DIR)
	protoc \
		-I $(PROTO_DIR) \
		--go_out=$(PROTO_OUT_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(PROTO_OUT_DIR) --go-grpc_opt=paths=source_relative \
		$(PROTO_DIR)/restaurant.proto
