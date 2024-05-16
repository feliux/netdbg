all: help

BINARY_NAME = netdbg
BIN_FOLDER := bin/
GO_LDFLAGS = -ldflags "-s -w"

build:
	GOARCH=arm64 GOOS=darwin go build $(GO_LDFLAGS) -o $(BIN_FOLDER)${BINARY_NAME}-darwin main.go
	GOARCH=amd64 GOOS=linux go build $(GO_LDFLAGS) -o $(BIN_FOLDER)${BINARY_NAME} main.go
	GOARCH=amd64 GOOS=windows go build $(GO_LDFLAGS) -o $(BIN_FOLDER)${BINARY_NAME}-windows main.go

clean:
	rm -f $(BIN_FOLDER)*

help:
	@echo '	Usage:'
	@echo '	  author	Display author.'
	@echo '	  clean		Clean binay file.'
	@echo '	  build		Compile to binary.'

author:
	@echo "	Project by feliux"
	@echo "	https://github.com/feliux"
