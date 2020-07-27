OUTPUT?=build/gsn


all: build install test
.PHONY: all

build:
	CGO_CFLAGS="-I/usr/lib" CGO_LDFLAGS="-L/usr/lib/rocksdb -lrocksdb -lstdc++ -lm -lz -lbz2 -lsnappy -llz4 -lzstd"  go build -o $(OUTPUT) .

install:
	go install ./main.go


test:
	@echo "--> Running test..."


clean:
	rm -rf build/

