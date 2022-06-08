.Phony: build
build:
	go build -o ${GOPATH}/bin/auto-standup main.go

.Phony: generate
generate:
	go generate ./...

.Phony: tests
tests:
	go test -p 1 ./...
