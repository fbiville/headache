.PHONY: all build clean test gen-mocks check-mockery help

OUTPUT = ./headache
GO_SOURCES = $(shell find . -type f -name '*.go')
GOBIN ?= $(shell go env GOPATH)/bin

.DEFAULT_GOAL := help

all: clean gen-mocks build test install ## run everything

build: $(OUTPUT) ## build the project binary

test: ## run the project tests
	GO111MODULE=on go test -v ./...

check-mockery: ## check whether mockery is installed
	@which mockery > /dev/null || ((echo 'mockery not found, please run: "cd `mktemp -d` && go get -u github.com/vektra/mockery/.../ && cd -"') && false)

gen-mocks: check-mockery ## generate the project mocks
	GO111MODULE=on mockery -output internal/pkg/helper_mocks -outpkg helper_mocks -dir internal/pkg/helper -name Clock \
	GO111MODULE=on mockery -output internal/pkg/vcs_mocks -outpkg vcs_mocks -dir internal/pkg/vcs -name Vcs \
	GO111MODULE=on mockery -output internal/pkg/vcs_mocks -outpkg vcs_mocks -dir internal/pkg/vcs -name VersioningClient \
	GO111MODULE=on mockery -output internal/pkg/fs_mocks -outpkg fs_mocks -dir internal/pkg/fs -name FileWriter \
	GO111MODULE=on mockery -output internal/pkg/fs_mocks -outpkg fs_mocks -dir internal/pkg/fs -name FileReader \
	GO111MODULE=on mockery -output internal/pkg/fs_mocks -outpkg fs_mocks -dir internal/pkg/fs -name File \
	GO111MODULE=on mockery -output internal/pkg/fs_mocks -outpkg fs_mocks -dir internal/pkg/fs -name PathMatcher \
	GO111MODULE=on mockery -output internal/pkg/core_mocks -outpkg core_mocks -dir internal/pkg/core -name ExecutionTracker

install: build ## copy the binary to GOBIN
	cp $(OUTPUT) $(GOBIN)

$(OUTPUT): $(GO_SOURCES)
	GO111MODULE=on go build -gcflags="all=-N -l" ./cmd/headache

clean: ## remove the binary
	rm -f $(OUTPUT)

# source: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help: ## Print help for each make target
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'