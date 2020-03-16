
all: build

tools:
	which gometalinter || ( go get -u github.com/alecthomas/gometalinter && gometalinter --install )
	which goveralls || go get github.com/mattn/goveralls

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

build:
	env CGO_ENABLED=0 go build -i

install:
	env CGO_ENABLED=0 go install

clean:
	rm -rf dist/
	go clean -i

coverall:
	goveralls -service=travis-ci -package github.com/bpineau/kube-deployments-notifier/pkg/...

test:
	go test -i github.com/bpineau/kube-deployments-notifier/...
	go test -race -cover github.com/bpineau/kube-deployments-notifier/...

.PHONY: tools lint fmt install clean coverall test all
