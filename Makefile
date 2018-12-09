GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
BINARY_NAME=bootstrap

all:	clean build build-shared build-example
build: 
		$(GOBUILD) -o bin/$(BINARY_NAME) -v src/main.go
build-shared:
		$(GOBUILD) -o bin/$(BINARY_NAME).so -buildmode=c-shared -v src/main.go
build-example: build-shared
		g++ -o bin/$(BINARY_NAME)_example -I./bin -L./bin example/example.cpp bin/bootstrap.so
clean:
		$(GOCLEAN)
		rm -f bin/$(BINARY_NAME)*