GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
BINARY_NAME=bootstrap
GOARM_ENV=CC=$(CC_ARM) CGO_ENABLED=1 GOOS=linux GOARCH=arm

CC=gcc
CC_ARM=arm-linux-gnueabihf-gcc

all:	clean build build-shared build-example
build: clean
		$(GOBUILD) -o bin/$(BINARY_NAME) -v src/main.go
build-shared:
		$(GOBUILD) -o bin/$(BINARY_NAME).so -buildmode=c-shared -v src/main.go
build-shared-arm: clean
		$(GOARM_ENV) $(GOBUILD) -o bin/$(BINARY_NAME).so -buildmode=c-shared -v src/main.go
build-example: clean build-shared
		$(CC) -o bin/$(BINARY_NAME)_example -I./bin -L./bin example/example.cpp bin/bootstrap.so
clean:
		$(GOCLEAN)
		rm -f bin/$(BINARY_NAME)*
