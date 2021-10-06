gen:
	protoc \
                -I test \
                -I third_party \
                --go_out=plugins=grpc,paths=source_relative:./test \
                --grpc-gateway_out=paths=source_relative:./test \
                test/example.proto

test:
	go test


