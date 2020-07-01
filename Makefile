DIST   ?= development
DEBUG  ?= --debug
VERSION = v0.6.x
LDFLAGS = -ldflags "-X uhppote.VERSION=$(VERSION)" 

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

debug: build
	go test ./... -run TestGrantALL


