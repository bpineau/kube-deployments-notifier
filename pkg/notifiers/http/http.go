// Package http notifier convey events notifications as json payloads to a
// configured HTTP API endpoint. Changes/creation notifications are sent
// with POST requests, while deletion are sent with DELETE requests.
package http

import (
	"bytes"
	"fmt"
	api "net/http"
	"time"

	"github.com/bpineau/kube-deployments-notifier/config"
)

// Notifier implements notifiers.Notifier.
type Notifier struct {
}

// Changed sends a notification to the configured HTTP API endpoint.
func (l *Notifier) Changed(c *config.KdnConfig, msg string) error {
	return l.push(c, "POST", msg)
}

// Deleted sends a notification to the configured HTTP API endpoint.
func (l *Notifier) Deleted(c *config.KdnConfig, msg string) error {
	return l.push(c, "DELETE", msg)
}

func (l *Notifier) push(c *config.KdnConfig, method string, msg string) error {
	if len(c.Endpoint) == 0 {
		return nil
	}

	req, err := api.NewRequest(method, c.Endpoint, bytes.NewBuffer([]byte(msg)))
	if err != nil {
		return err
	}

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

	if resp.StatusCode >= 400 {
		return fmt.Errorf("HTTP request failed (code=%d)", resp.StatusCode)
	}

	return resp.Body.Close()
}
