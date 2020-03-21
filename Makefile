GOLANGCI_VERSION = 1.24.0
export GO111MODULE := on

ifndef $(GOPATH)
    GOPATH=$(shell go env GOPATH)
    export GOPATH
endif

all: build

bin/golangci-lint:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v${GOLANGCI_VERSION}

${GOPATH}/bin/goveralls:
	which goveralls || go get github.com/mattn/goveralls

tools: bin/golangci-lint ${GOPATH}/bin/goveralls

lint: tools
	./bin/golangci-lint run --concurrency=1 --timeout=300s --disable-all \
		--enable=golint \
		--enable=vet \
		--enable=vetshadow \
		--enable=varcheck \
		--enable=errcheck \
		--enable=structcheck \
		--enable=deadcode \
		--enable=ineffassign \
		--enable=dupl \
		--enable=varcheck \
		--enable=interfacer \
		--enable=goconst \
		--enable=megacheck \
		--enable=unparam \
		--enable=misspell \
		--enable=gas \
		--enable=goimports \
		--enable=gocyclo

fmt:
	go fmt ./...

build:
	env CGO_ENABLED=0 go build

install:
	env CGO_ENABLED=0 go install

clean:
	rm -rf dist/ bin/
	go clean -i

coverall: ${GOPATH}/bin/goveralls
	goveralls -service=travis-ci -package github.com/bpineau/kube-deployments-notifier/pkg/...

test:
	go test github.com/bpineau/kube-deployments-notifier/...
	go test -race -cover github.com/bpineau/kube-deployments-notifier/...

.PHONY: tools lint fmt install clean coverall test all ${GOPATH}/bin/goveralls bin/golangci-lint
