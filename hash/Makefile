PWD := $(shell pwd)
VERSION := $(shell cat VERSION.txt)
APP_NAME := hash
HASH_TMP_FILE=tests/hash.txt

default: build test

build:
	CGO_ENABLED=0 go build -ldflags="-X main.version=$(VERSION)" -o bin/$(APP_NAME) .

test:
	bin/$(APP_NAME) -v
	bin/$(APP_NAME) md5 -R -H --exclude "***_test.go" --exclude "vendor" --exclude ccc.goo ./

integration-test: build
	bin/$(APP_NAME) sha1 -R -U -o $(HASH_TMP_FILE) \
	 --exclude "***_test.go" \
	 --exclude "vendor" \
	 --exclude "tests" \
	 -f "{{.Hash}}  {{.Name}}" ./
	tests/check.sh $(HASH_TMP_FILE)

clear:
	rm -f $(HASH_TMP_FILE)

dev:
	echo "$(HASH_TMP_FILE)"