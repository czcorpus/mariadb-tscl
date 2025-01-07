VERSION=`git describe --tags --always`
BUILD=`date +%FT%T%z`
HASH=`git rev-parse --short HEAD`


LDFLAGS=-ldflags "-w -s -X main.version=${VERSION} -X main.buildDate=${BUILD} -X main.gitCommit=${HASH}"

all: test build

build:
	go build ${LDFLAGS}

clean:
	rm ./mariadb-tscl

test:
	go test ./...

.PHONY: clean install test build

