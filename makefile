proto_generate:
	mkdir -p pb && protoc --proto_path=proto --go_out=plugins=grpc:. proto/*.proto

proto_clean:
	rm pb/*.pb.go
