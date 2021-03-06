VERSION = v0.7.0
DIST   ?= development
DEBUG  ?= --debug
LDFLAGS = -ldflags "-X uhppote.VERSION=$(VERSION)" 

.PHONY: bump

all: test      \
	 benchmark \
     coverage

clean:
	go clean
	rm -rf bin

format: 
	go fmt ./...

build: format
	go build ./...

test: build
	go test ./...

vet: build
	go vet ./...

lint: build
	golint ./...

benchmark: build
	go test -bench ./...

coverage: build
	go test -cover ./...

build-all: test vet
	env GOOS=linux   GOARCH=amd64       go build ./...
	env GOOS=linux   GOARCH=arm GOARM=7 go build ./...
	env GOOS=darwin  GOARCH=amd64       go build ./...
	env GOOS=windows GOARCH=amd64       go build ./...

release: build-all

bump:
	go get -u github.com/uhppoted/uhppote-core
	go get -u golang.org/x/sys

debug: build
	go test ./... -run TestPutACLWithInvalidTimeProfile

godoc:
	godoc -http=:80	-index_interval=60s
