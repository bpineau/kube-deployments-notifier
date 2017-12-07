package deployment

import (
	"github.com/bpineau/kube-deployments-notifier/config"
	"github.com/bpineau/kube-deployments-notifier/pkg/controllers"

	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	appsv1beta1 "k8s.io/client-go/pkg/apis/apps/v1beta1"
	"k8s.io/client-go/tools/cache"
)

// Controller monitors deployments
type Controller struct {
	// https://golang.org/doc/effective_go.html#embedding
	controllers.CommonController
}

// Init initialize deployment controller
func (c *Controller) Init(conf *config.KdnConfig) controllers.Controller {
	c.CommonController = controllers.CommonController{
		Conf: conf,
		Name: "deployment",
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