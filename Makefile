.PHONY: gen-go
gen-go:
	@for proto in ./proto/*.proto; do \
		name=$$(basename $$proto .proto); \
		mkdir -p ./go/gen_pb/$${name}pb; \
		protoc --proto_path=./proto \
			--go_out=paths=source_relative:./go/gen_pb/$${name}pb \
			--go-grpc_out=paths=source_relative:./go/gen_pb/$${name}pb \
			$$proto; \
	done