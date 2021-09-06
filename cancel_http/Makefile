.PHONY: proto
proto:
	protoc \
	-I=./helloworld \
	--go_out=plugins=grpc:./helloworld \
	./helloworld/*.proto

.PHONY: run
run:
	curl -v -m 2 http://localhost:18081/sleep
