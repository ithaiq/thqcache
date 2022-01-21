.PHONY: wire
wire:
	wire ./...

.PHONY: test
test:
	go test -cover -race ./...

.PHONY: proto
proto:
	protoc --go_out=plugins=grpc:./proto ./proto/*.proto