package http

import (
	"bytes"
	"time"
	api "net/http"

	"github.com/bpineau/kube-deployments-notifier/config"
)

// Notifier implements notifiers.Notifier
type Notifier struct {
}

// Changed sends notification to the configured logrus logger
func (l *Notifier) Changed(c *config.KdnConfig, msg string) error {
	return l.push(c, "POST", msg)
}

// Deleted sends notification to the configured logrus logger
func (l *Notifier) Deleted(c *config.KdnConfig, msg string) error {
	return l.push(c, "DELETE", msg)
}

func (l *Notifier) push(c *config.KdnConfig, method string, msg string) error {
	if len(c.Endpoint) == 0 {
		return nil
	}

	req, err := api.NewRequest(method, c.Endpoint, bytes.NewBuffer([]byte(msg)))
	req.Header.Set("Content-Type", "application/json")
	if c.TokenHdr != "" && c.TokenVal != "" {
		req.Header.Set(c.TokenHdr, c.TokenVal)
	}

	timeout := time.Duration(10 * time.Second)
	client := &api.Client{Timeout: timeout}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	return resp.Body.Close()
}
