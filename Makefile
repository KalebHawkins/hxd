.DEFAULT_GOAL := build

BIN_DIR := bin
BINARY_NAME := hxd

ifeq ($(OS),Windows_NT)
	BINARY := $(BIN_DIR)/$(BINARY_NAME).exe
	CLEAN_CMD := if exist $(BIN_DIR) rmdir /s /q $(BIN_DIR)
else
	BINARY := $(BIN_DIR)/$(BINARY_NAME)
	CLEAN_CMD := rm -rf $(BIN_DIR)
endif

.PHONY: fmt vet build clean

fmt:
	go fmt ./...

vet: fmt
	go vet ./...

build: vet
	go build -o $(BINARY) .

clean:
	go clean
	$(CLEAN_CMD)