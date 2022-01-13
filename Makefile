.PHONY: wire
wire:
	wire ./...

.PHONY: test
test:
	go test -cover -race ./...