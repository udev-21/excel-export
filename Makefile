.PHONY: build build-linux build-mac build-win build-all

build-linux:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bin/linux_amd64 -v
build-mac:
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -o bin/darwin_amd64 -v
build-win:
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o bin/windows_amd64 -v

build-all:
	make build-linux
	make build-mac
	make build-win
