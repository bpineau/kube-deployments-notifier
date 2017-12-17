# kube-deployments-notifier

[![Build Status](https://travis-ci.org/bpineau/kube-deployments-notifier.svg?branch=master)](https://travis-ci.org/bpineau/kube-deployments-notifier)
[![Go Report Card](https://goreportcard.com/badge/github.com/bpineau/kube-deployments-notifier)](https://goreportcard.com/report/github.com/bpineau/kube-deployments-notifier)

An example Kubernetes controller that list and watch deployments, and send
them as json payload to a remote API endpoint.

## Build

Assuming you have go 1.9 and glide in the path, and GOPATH configured:

```shell
make deps
make build
```

## Usage

The daemon may run either as a pod, or outside of the Kubernetes cluster.
He should find the Kubernetes api-server automatically (or you can use the
"-s" or "-k" flags). You can pass parameters from cli args, env, config
files, or both.

Example:
```
kube-deployments-notifier \
  -l 'vendor=mycompany,app!=mmp-database' \
  -t Authorization -a "Bearer vbnf3hjklp5iuytre" \
  -e http://myapiserver:8042 
```

The command line flags are (all optionals):
```
Usage:
  kube-deployments-notifier [flags]

Flags:
  -s, --api-server string      kube api server url
  -c, --config string          configuration file (default "/etc/kdn/kube-deployments-notifier.yaml")
  -d, --dry-run                dry-run mode
  -e, --endpoint string        API endpoint
  -l, --filter string          Label filter
  -p, --healthcheck-port int   port for answering healthchecks
  -h, --help                   help for kube-deployments-notifier
  -k, --kube-config string     kube config path
  -v, --log-level string       log level (default "debug")
  -o, --log-output string      log output (default "stderr")
  -r, --log-server string      log server (if using syslog)
  -t, --token-header string    token header name
  -a, --token-value string     token header value
```

Using an (optional) configuration file:
```yaml
dry-run: false
healthcheck-port: 8080
api-server: http://example.com:8080
endpoint: https://my-api-endpoint/foo/bar
token-header: Authorization
token-value: "Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9"
filter: "vendor=mycompany,app!=mmp-database"

log:
  output: "stdout"
  level: "debug"
```

The environment variable consumed by kube-deployments-notifier are option names prefixed
by ```KDN_``` and using underscore instead of dash. Except KUBECONFIG,
used without a prefix (to match kubernetes conventions).
