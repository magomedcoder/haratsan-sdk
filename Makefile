.PHONY: gen
gen: gen-go gen-python

gen-go:
	@for proto in ./proto/*.proto; do \
		name=$$(basename $$proto .proto); \
		mkdir -p ./go/gen_pb/$${name}pb; \
		protoc --proto_path=./proto \
			--go_out=paths=source_relative:./go/gen_pb/$${name}pb \
			--go-grpc_out=paths=source_relative:./go/gen_pb/$${name}pb \
			$$proto; \
	done

gen-python:
	@mkdir -p ./python/haratsan_sdk/gen_pb && \
	 python3 -m grpc_tools.protoc -I./proto \
		--python_out=./python/haratsan_sdk/gen_pb \
		--grpc_python_out=./python/haratsan_sdk/gen_pb \
		./proto/bot_api.proto && \
	 sed -i 's/^import bot_api_pb2/from . import bot_api_pb2/' ./python/haratsan_sdk/gen_pb/bot_api_pb2_grpc.py