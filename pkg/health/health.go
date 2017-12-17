package health

import (
	"fmt"
	"io"
	"net/http"

	"github.com/bpineau/kube-deployments-notifier/config"
)

func healthCheckReply(w http.ResponseWriter, r *http.Request) {
	if _, err := io.WriteString(w, "ok\n"); err != nil {
		fmt.Printf("Failed to reply to http healtcheck: %s\n", err)
	}
}

// HeartBeatService exposes an http healthcheck handler
func HeartBeatService(c *config.KdnConfig) error {
	if c.HealthPort == 0 {
		return nil
	}
	http.HandleFunc("/health", healthCheckReply)
	return http.ListenAndServe(fmt.Sprintf(":%d", c.HealthPort), nil)
}
