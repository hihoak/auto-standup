.Phony: build
build:
	go build -o ${GOPATH}/bin/auto-standup main.go

.Phony: build_amd
build:
	GOOS="linux" GOARCH="amd64" go build -o ${GOPATH}/bin/auto-standup main.go

.Phony: build_win
build:
	GOOS="windows" GOARCH="amd64" go build -o ${GOPATH}\bin\auto-standup main.go

.Phony: build_arm
build:
	GOOS="darwin" GOARCH="arm64" go build -o ${GOPATH}/bin/auto-standup main.go
