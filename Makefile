GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
BINARY_NAME=bootstrap

all:	clean build
build: 
		$(GOBUILD) -o $(BINARY_NAME) -v src/main.go
clean:
		$(GOCLEAN)
		rm -f $(BINARY_NAME)
