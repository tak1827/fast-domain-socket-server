PROTO_SRC_FILES=$(shell find ./proto -type f -name "*.proto" | sed 's/\/proto//g')

.PHONY: proto
proto:
	rm ./data/*.pb.go; \
	cd ./proto; \
	protoc -I=. -I=${GOPATH}/src/github.com/protobuf \
		--gofast_out=paths=source_relative:../data  \
		$(PROTO_SRC_FILES); \

bench:
	go test ./... -bench=. -benchtime=10s

fmt:
	go fmt ./...

lint:
	go vet ./...

