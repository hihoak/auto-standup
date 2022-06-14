.Phony: build
build:
	go build -o ${GOPATH}/bin/auto-standup main.go

.Phony: generate
generate:
	go generate ./...

.Phony: tests
tests:
	go test -v -race -covermode=atomic -coverpkg `go list ./... | grep -v mocks | tr '\n' ','` -coverprofile=coverage.out ./...
