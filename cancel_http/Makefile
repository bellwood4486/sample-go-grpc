.PHONY: proto
proto:
	protoc \
	-I=./helloworld \
	--go_out=plugins=grpc:./helloworld \
	./helloworld/*.proto
