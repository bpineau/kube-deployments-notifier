package config

import (
	"fmt"
	"time"

	"github.com/bpineau/kube-deployments-notifier/pkg/clientset"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// KdnConfig is the main program configuration, passed to controllers Init()
type KdnConfig struct {
	DryRun     bool
	Logger     *logrus.Logger
	ClientSet  kubernetes.Interface
	Endpoint   string
	TokenHdr   string
	TokenVal   string
	Filter     string
	HealthPort int
	ResyncIntv time.Duration
}

// Init initialize the configuration (creating the ClientSet for the cluster)
func (c *KdnConfig) Init(apiserver string, kubeconfig string) {
	var err error

	c.ClientSet, err = clientset.NewClientSet(apiserver, kubeconfig)
	if err != nil {
		panic(fmt.Errorf("Failed init Kubernetes clientset: %+v", err))
	}

	_, err = c.ClientSet.CoreV1().Namespaces().List(metav1.ListOptions{})
	if err != nil {
		panic(fmt.Errorf("Failed to query Kubernetes api-server: %+v", err))
	}

	c.Logger.Info("Kubernetes clientset initialized")
}
