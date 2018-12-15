VERSION=1.0.0-alpha2
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean

CC=gcc
CC_WIN=i686-w64-mingw32-gcc
CC_ARM=arm-linux-gnueabihf-gcc

GOWIN=CC=$(CC_WIN) CGO_ENABLED=1 GOOS=windows GOARCH=386
GOARM=CC=$(CC_ARM) CGO_ENABLED=1 GOOS=linux GOARCH=arm

GOARM_SHARED=CC=$(CC_ARM) CGO_ENABLED=1 GOOS=linux GOARCH=arm

BASE=bootstrap
BINARY_NAME=MiSTer-bootstrap-linux-amd64-$(VERSION)
BINARY_NAME_WIN=MiSTer-bootstrap-windows-amd64-$(VERSION).exe

all:	clean build build-windows build-arm build-shared build-example
build: clean
		$(GOBUILD) -o bin/$(BINARY_NAME) -v src/main.go
build-windows:
		$(GOWIN) $(GOBUILD) -o bin/$(BINARY_NAME_WIN) -v src/main.go
build-arm:
		$(GOARM) $(GOBUILD) -o bin/$(BINARY_NAME) -v src/main.go
build-shared:
		$(GOBUILD) -o bin/$(BASE).so -buildmode=c-archive -v src/main.go
build-shared-arm: clean
		$(GOARM_SHARED) $(GOBUILD) -o bin/$(BASE).a -buildmode=c-archive -v src/main.go
build-example: clean build-shared
		$(CC) -o bin/$(BASE)_example -I./bin -L./bin example/example.cpp bin/bootstrap.so
clean:
		$(GOCLEAN)
		rm -f bin/$(BASE)*
		rm -f bin/$(BINARY_NAME)*
		rm -f bin/$(BINARY_NAME_WIN)*
