# kube-deployments-notifier

[![Build Status](https://github.com/bpineau/kube-deployments-notifier/workflows/CI/badge.svg)](https://github.com/bpineau/kube-deployments-notifier/actions)
[![Coverage Status](https://coveralls.io/repos/github/bpineau/kube-deployments-notifier/badge.svg?branch=master)](https://coveralls.io/github/bpineau/kube-deployments-notifier?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/bpineau/kube-deployments-notifier)](https://goreportcard.com/report/github.com/bpineau/kube-deployments-notifier)

An example Kubernetes controller that list and watch deployments, and send
them as json payload to a remote API endpoint.

## Build

Assuming you have go 1.14:

```shell
make build
```

## Usage

The daemon may run either as a pod, or outside of the Kubernetes cluster.
He should find the Kubernetes api-server automatically (or you can use the
"-s", "-x" or "-k" flags). You can pass parameters from cli args, env, config
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
  -x, --context string         kube context
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
  -i, --resync-interval int    resync interval in seconds (0 to disable) (default 900)
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
resync-interval: 900

log:
  output: "stdout"
  level: "debug"
```

The environment variable consumed by kube-deployments-notifier are option names prefixed
by ```KDN_``` and using underscore instead of dash. Except KUBECONFIG,
used without a prefix (to match kubernetes conventions).

## Docker image

A ready to use, public docker image is available at [Docker Hub](https://hub.docker.com/r/bpineau/kube-deployments-notifier/), published at each release.
You can use it directly from your Kubernetes deployments, ie.

```yaml
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: kube-deployments-notifier
  namespace: kube-system
  labels:
    k8s-app: kube-deployments-notifier
spec:
  selector:
    matchLabels:
      k8s-app: kube-deployments-notifier
  replicas: 1
  template:
    metadata:
      labels:
        k8s-app: kube-deployments-notifier
    spec:
      serviceAccountName: kube-deployments-notifier
      containers:
        - name: kube-deployments-notifier
          image: bpineau/kube-deployments-notifier:v0.6.0
          args:
            - --log-output=stdout
            - --filter=vendor=mycompany,app!=mmp-database
            - --endpoint=http://mycollector
            - --healthcheck-port=8080
          resources:
            requests:
              cpu: 0.1
              memory: 50Mi
            limits:
              cpu: 0.2
              memory: 100Mi
          livenessProbe:
            httpGet:
              path: /health
              port: 8080
            timeoutSeconds: 5
            initialDelaySeconds: 10

---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: kube-deployments-notifier
  namespace: kube-system
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: kube-deployments-notifier
  namespace: kube-system
rules:
  - apiGroups:
      - ""
      - apps
    resources:
      - namespaces
      - deployments
    verbs:
      - list
      - get
      - watch
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: kube-deployments-notifier
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: kube-deployments-notifier
subjects:
  - kind: ServiceAccount
    name: kube-deployments-notifier
    namespace: kube-system
```
