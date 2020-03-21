// Package clientset initialize a Kubernete's client-go "clientset" (an initialized
// connection to the Kubernete's api-server) according the configuration.
package clientset

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	// Ensure we have various auth method linked in
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
)

// NewClientSet create a clientset (a client connection to a Kubernetes cluster).
// It will connect using the optional apiserver or kubeconfig options, or will
// default to the automatic, in cluster settings.
func NewClientSet(apiserver, context, kubeconfig string) (*kubernetes.Clientset, error) {
	overrides := clientcmd.ConfigOverrides{}
	loader := clientcmd.NewDefaultClientConfigLoadingRules()
	if kubeconfig != "" {
		loader.ExplicitPath = kubeconfig
	}

	if context != "" {
		overrides.CurrentContext = context
	}

	config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loader, &overrides).ClientConfig()

	if err != nil {
		return nil, err
	}

	if apiserver != "" {
		config.Host = apiserver
	}

	return kubernetes.NewForConfig(config)
}
