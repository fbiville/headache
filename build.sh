#!/bin/bash
set -euo pipefail

dep ensure

rm helper_mocks/*.go || true \
    && rm vcs_mocks/*.go || true \
    && rm fs_mocks/*.go || true \
    && rm core_mocks/*.go || true

go get github.com/vektra/mockery/.../ \
    && mockery -output helper_mocks -outpkg helper_mocks -dir helper -name Clock \
    && mockery -output vcs_mocks -outpkg vcs_mocks -dir vcs -name Vcs \
    && mockery -output vcs_mocks -outpkg vcs_mocks -dir vcs -name VersioningClient \
    && mockery -output fs_mocks -outpkg fs_mocks -dir fs -name FileWriter \
    && mockery -output fs_mocks -outpkg fs_mocks -dir fs -name FileReader \
    && mockery -output fs_mocks -outpkg fs_mocks -dir fs -name File \
    && mockery -output fs_mocks -outpkg fs_mocks -dir fs -name PathMatcher \
    && mockery -output core_mocks -outpkg core_mocks -dir core -name ExecutionTracker

go build ./...
go clean -testcache && go test -v ./...
rm headache 2> /dev/null || true && go build
