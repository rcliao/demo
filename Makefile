# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_FOLDER=bin
BINARY_NAME=demo
BINARY_UNIX=$(BINARY_NAME)_unix
BINARY_WINDOWS=$(BINARY_NAME)_windows
MAIN_FILE=main.go

all: clean install test build
install:
	$(GOGET) -t ./...
build: 
	$(GOBUILD) -o $(BINARY_FOLDER)/$(BINARY_NAME) -v $(MAIN_FILE)
test: install
	$(GOTEST) -v ./...
clean: 
	$(GOCLEAN)
	rm -f $(BINARY_FOLDER)/$(BINARY_NAME)
	rm -f $(BINARY_FOLDER)/$(BINARY_UNIX)
run: build
	./$(BINARY_FOLDER)/$(BINARY_NAME)

build-linux:    ## Cross compilation
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_FOLDER)/$(BINARY_UNIX) -v $(MAIN_FILE)
build-windows:    ## Cross compilation
	GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BINARY_FOLDER)/$(BINARY_WINDOWS).exe -v $(MAIN_FILE)



# ----------------------------------------------------------------------------
# Self-Documented Makefile
# ----------------------------------------------------------------------------
help:						## (DEFAULT) This help information
	@echo ====================================================================
	@grep -E '^## .*$$'  \
		$(MAKEFILE_LIST)  \
		| awk 'BEGIN { FS="## " }; {printf "\033[33m%-20s\033[0m \n", $$2}'
	@echo
	@grep -E '^[0-9a-zA-Z_-]+:.*?## .*$$'  \
		$(MAKEFILE_LIST)  \
		| awk 'BEGIN { FS=":.*?## " }; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'  \
#		 | sort
.PHONY: help
.DEFAULT_GOAL := all

