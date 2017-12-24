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

var (
	conts = []controllers.Controller{
		&deployment.Controller{},
	}
)

// Run launchs the effective controllers goroutines
func Run(config *config.KdnConfig, notif notifiers.Notifier) {
	wg := sync.WaitGroup{}
	wg.Add(len(conts))
	defer wg.Wait()

	for _, c := range conts {
		go c.Init(config, notif).Start(&wg)
		defer func(c controllers.Controller) {
			go c.Stop()
		}(c)
	}

	go func() {
		if err := health.HeartBeatService(config); err != nil {
			log.Fatal("Healtcheck service failed: ", err)
		}
	}()

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGTERM)
	signal.Notify(sigterm, syscall.SIGINT)
	<-sigterm

	config.Logger.Infof("Stopping all controllers")
}
