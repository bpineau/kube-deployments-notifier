package run

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/bpineau/kube-deployments-notifier/config"
	"github.com/bpineau/kube-deployments-notifier/pkg/controllers"
	"github.com/bpineau/kube-deployments-notifier/pkg/controllers/deployment"
	"github.com/bpineau/kube-deployments-notifier/pkg/health"
	"github.com/bpineau/kube-deployments-notifier/pkg/notifiers"
)

var conts = []controllers.Controller{
	&deployment.Controller{},
}

// Run launchs the effective controllers goroutines
func Run(config *config.KdnConfig) {
	wg := sync.WaitGroup{}
	wg.Add(len(conts))
	defer wg.Wait()

	notifiers := notifiers.Init(notifiers.Backends)

	for _, c := range conts {
		go c.Init(config, notifiers).Start(&wg)
		defer func(c controllers.Controller) {
			go c.Stop()
		}(c)
	}

	go log.Println(health.HeartBeatService(config))

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGTERM)
	signal.Notify(sigterm, syscall.SIGINT)
	<-sigterm

	config.Logger.Infof("Stopping all controllers")
}
