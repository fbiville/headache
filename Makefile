.PHONY: build clean test gen-mocks check-mockery help

OUTPUT = ./headache
GO_SOURCES = $(shell find . -type f -name '*.go')
GOBIN ?= $(shell go env GOPATH)/bin

.DEFAULT_GOAL := help

build: $(OUTPUT) ## build the project binary

test: ## run the project tests
	GO111MODULE=on go test -v ./...

check-mockery: ## check whether mockery is installed
	@which mockery > /dev/null || (echo mockery not found: issue \"GO111MODULE=off go get -u  github.com/vektra/mockery/.../\" && false)

gen-mocks: check-mockery ## generate the project mocks
	GO111MODULE=on mockery -output helper_mocks -outpkg helper_mocks -dir helper -name Clock \
	GO111MODULE=on mockery -output vcs_mocks -outpkg vcs_mocks -dir vcs -name Vcs \
	GO111MODULE=on mockery -output vcs_mocks -outpkg vcs_mocks -dir vcs -name VersioningClient \
	GO111MODULE=on mockery -output fs_mocks -outpkg fs_mocks -dir fs -name FileWriter \
	GO111MODULE=on mockery -output fs_mocks -outpkg fs_mocks -dir fs -name FileReader \
	GO111MODULE=on mockery -output fs_mocks -outpkg fs_mocks -dir fs -name File \
	GO111MODULE=on mockery -output fs_mocks -outpkg fs_mocks -dir fs -name PathMatcher \
	GO111MODULE=on mockery -output core_mocks -outpkg core_mocks -dir core -name ExecutionTracker

install: build ## copy the binary to GOBIN
	cp $(OUTPUT) $(GOBIN)

$(OUTPUT): $(GO_SOURCES)
	GO111MODULE=on go build

clean: ## remove the binary
	rm -f $(OUTPUT)

# source: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help: ## Print help for each make target
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'