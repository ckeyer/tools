PWD := $(shell pwd)
VERSION := $(shell cat VERSION.txt)
APP_NAME := hash

default: build test

build:
	CGO_ENABLED=0 go build -ldflags="-X main.version=$(VERSION)" -o bin/$(APP_NAME) .

godep:

test:
	bin/$(APP_NAME) -v
	bin/$(APP_NAME) md5 -R --exclude "*.go" --exclude ccc.goo hello
	bin/$(APP_NAME) md5 -R -U

integration-test:
	bin/$(APP_NAME) sha1 -R -f "{{.Hash}} {{.Name}}" hello
	./check.sh