package deployment

import (
	"github.com/bpineau/kube-deployments-notifier/config"
	"github.com/bpineau/kube-deployments-notifier/pkg/controllers"
	"github.com/bpineau/kube-deployments-notifier/pkg/notifiers"

	appsv1beta1 "k8s.io/api/apps/v1beta1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
)

// Controller monitors Kubernetes' deployments objects in the cluster
type Controller struct {
	// https://golang.org/doc/effective_go.html#embedding
	controllers.CommonController
}

// Init initialize controller
func (c *Controller) Init(conf *config.KdnConfig, n notifiers.Notifier) controllers.Controller {
	c.CommonController = controllers.CommonController{
		Conf:      conf,
		Name:      "deployment",
		Notifiers: n,
	}

	client := c.Conf.ClientSet
	c.ObjType = &appsv1beta1.Deployment{}
	selector := meta_v1.ListOptions{LabelSelector: conf.Filter}
	c.ListWatch = &cache.ListWatch{
		ListFunc: func(options meta_v1.ListOptions) (runtime.Object, error) {
			return client.AppsV1beta1().Deployments(meta_v1.NamespaceAll).List(selector)
		},
		WatchFunc: func(options meta_v1.ListOptions) (watch.Interface, error) {
			return client.AppsV1beta1().Deployments(meta_v1.NamespaceAll).Watch(selector)
		},
	}

	return c
}
