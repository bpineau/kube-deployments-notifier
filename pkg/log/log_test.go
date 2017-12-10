package log

import (
	"fmt"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
)

func TestLog(t *testing.T) {
	logger := New("warning", "", "test")

	logger.Info("Changed: foo")
	logger.Warn("Changed: bar")
	logger.Error("Deleted: foo")

	hook := logger.Hooks[logrus.InfoLevel][0].(*test.Hook)
	if len(hook.Entries) != 2 {
		t.Errorf("Not the correct count of log entries")
	}

	logger.Warn("Changed: baz")
	if hook.LastEntry().Message != "Changed: baz" {
		t.Errorf("Unexpected log entry: %s", hook.LastEntry().Message)
	}

	logger = New("info", "192.0.2.0:514", "syslog")
	if fmt.Sprintf("%T", logger) != "*logrus.Logger" {
		t.Error("Failed to instantiate a syslog logger")
	}
}
