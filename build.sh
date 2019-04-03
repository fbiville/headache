#!/bin/bash
set -euo pipefail

go install

rm helper_mocks/*.go || true \
    && rm vcs_mocks/*.go || true \
    && rm fs_mocks/*.go || true \
    && rm core_mocks/*.go || true

GO111MODULE=off go get -u  github.com/vektra/mockery/.../ \
    && GO111MODULE=on mockery -output helper_mocks -outpkg helper_mocks -dir helper -name Clock \
    && GO111MODULE=on mockery -output vcs_mocks -outpkg vcs_mocks -dir vcs -name Vcs \
    && GO111MODULE=on mockery -output vcs_mocks -outpkg vcs_mocks -dir vcs -name VersioningClient \
    && GO111MODULE=on mockery -output fs_mocks -outpkg fs_mocks -dir fs -name FileWriter \
    && GO111MODULE=on mockery -output fs_mocks -outpkg fs_mocks -dir fs -name FileReader \
    && GO111MODULE=on mockery -output fs_mocks -outpkg fs_mocks -dir fs -name File \
    && GO111MODULE=on mockery -output fs_mocks -outpkg fs_mocks -dir fs -name PathMatcher \
    && GO111MODULE=on mockery -output core_mocks -outpkg core_mocks -dir core -name ExecutionTracker

go build ./...
go clean -testcache && go test -v ./...
rm headache 2> /dev/null || true && go build
