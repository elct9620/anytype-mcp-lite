.PHONY: all build clean

all: build

build:
	go build -o ./dist/anytype-mcp ./cmd

clean:
	rm -rf ./dist