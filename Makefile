.PHONY: build clean test all gen-mocks check-mockery

OUTPUT = ./headache
GO_SOURCES = $(shell find . -type f -name '*.go')
GOBIN ?= $(shell go env GOPATH)/bin

all: build test docs

build: $(OUTPUT)

test:
	GO111MODULE=on go test -v ./...

check-mockery:
	@which mockery > /dev/null || (echo mockery not found: issue \"GO111MODULE=off go get -u  github.com/vektra/mockery/.../\" && false)

gen-mocks: check-mockery
	GO111MODULE=on mockery -output helper_mocks -outpkg helper_mocks -dir helper -name Clock \
	GO111MODULE=on mockery -output vcs_mocks -outpkg vcs_mocks -dir vcs -name Vcs \
	GO111MODULE=on mockery -output vcs_mocks -outpkg vcs_mocks -dir vcs -name VersioningClient \
	GO111MODULE=on mockery -output fs_mocks -outpkg fs_mocks -dir fs -name FileWriter \
	GO111MODULE=on mockery -output fs_mocks -outpkg fs_mocks -dir fs -name FileReader \
	GO111MODULE=on mockery -output fs_mocks -outpkg fs_mocks -dir fs -name File \
	GO111MODULE=on mockery -output fs_mocks -outpkg fs_mocks -dir fs -name PathMatcher \
	GO111MODULE=on mockery -output core_mocks -outpkg core_mocks -dir core -name ExecutionTracker

install: build
	cp $(OUTPUT) $(GOBIN)

$(OUTPUT): $(GO_SOURCES)
	GO111MODULE=on go build ./...

clean:
	rm -f $(OUTPUT)
