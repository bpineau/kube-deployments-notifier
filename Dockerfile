FROM golang:1.9.2 as builder
WORKDIR /go/src/github.com/bpineau/kube-deployments-notifier
COPY . .
RUN go get -u github.com/Masterminds/glide
RUN make deps
RUN make build

FROM alpine:3.7
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/src/github.com/bpineau/kube-deployments-notifier/kube-deployments-notifier /usr/bin/
ENTRYPOINT ["/usr/bin/kube-deployments-notifier"]
