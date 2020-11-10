gen:
	protoc --proto_path=test --proto_path=third_party --go_out=plugins=grpc:. example.proto
	protoc --proto_path=test --proto_path=third_party --grpc-gateway_out=logtostderr=true:./test example.proto

test:
	go test


