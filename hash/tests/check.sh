#!/bin/bash

docker run --rm -v $(PWD):/opt/ -w /opt  alpine:edge sha1sum -c 