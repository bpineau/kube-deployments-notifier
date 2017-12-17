
all: build

tools:
	which gometalinter || ( go get -u github.com/alecthomas/gometalinter && gometalinter --install )
	which glide || go get -u github.com/Masterminds/glide

lint:
	gometalinter --concurrency=1 --deadline=300s --vendor --disable-all \
		--enable=golint \
		--enable=vet \
		--enable=vetshadow \
		--enable=varcheck \
		--enable=errcheck \
		--enable=structcheck \
		--enable=deadcode \
		--enable=ineffassign \
		--enable=dupl \
		--enable=gotype \
		--enable=varcheck \
		--enable=interfacer \
		--enable=goconst \
		--enable=megacheck \
		--enable=unparam \
		--enable=misspell \
		--enable=gas \
		--enable=goimports \
		--enable=gocyclo \
		./...

fmt:
	go fmt ./...

deps:
	glide install

build:
	env CGO_ENABLED=0 go build -i

install:
	env CGO_ENABLED=0 go install

clean:
	rm -rf dist/
	go clean -i

test:
	go test -race -cover -v github.com/bpineau/kube-deployments-notifier/pkg/...

.PHONY: tools lint fmt install clean test all
