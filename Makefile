GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
BINARY_NAME=bootstrap
CGO_ENABLED=0

all:	clean build
build: 
		$(GOBUILD) -o bin/$(BINARY_NAME) -v src/main.go
clean:
		$(GOCLEAN)
		rm -f bin/$(BINARY_NAME)
