OUTPUT?=build/sng


all: build install test
.PHONY: all

build:
	go build -o $(OUTPUT) ./main.go

install:
	go install -o $(OUTPUT) ./main.go


test:
	@echo "--> Running test..."


#clean

