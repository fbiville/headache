#!/bin/bash
set -e

dep ensure
rm mocks/*.go || true && go get github.com/vektra/mockery/.../ && mockery -output mocks -dir versioning -name Vcs
go build ./...
go clean -testcache && go test ./...
rm header 2> /dev/null || true && go build
