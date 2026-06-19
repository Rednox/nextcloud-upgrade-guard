BINARY ?= nc-guard

.PHONY: test vet build

test:
go test ./...

vet:
go vet ./...

build:
go build -o ./bin/$(BINARY) ./cmd/nc-guard
