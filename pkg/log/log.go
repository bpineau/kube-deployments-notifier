package log

import (
	"io"
	"io/ioutil"
	"os"

	"log/syslog"

	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	ls "github.com/sirupsen/logrus/hooks/syslog"
)

// New initialize logrus and return a new logger
func New(logLevel string, logServer string, logOutput string) *logrus.Logger {
	var level logrus.Level
	var output io.Writer
	var hook logrus.Hook

	switch logOutput {
	case "stdout":
		output = os.Stdout
	case "stderr":
		output = os.Stderr
	case "test":
		output = ioutil.Discard
		_, hook = test.NewNullLogger()
	case "syslog":
		output = os.Stderr // does not matter ?
		if logServer == "" {
			panic("syslog output needs a log server (ie. 127.0.0.1:514)")
		}
		hook, _ = ls.NewSyslogHook("udp", logServer, syslog.LOG_INFO, "kube-deployments-notifier")
	default:
		output = os.Stderr
	}

	switch logLevel {
	case "debug":
		level = logrus.DebugLevel
	case "info":
		level = logrus.InfoLevel
	case "warning":
		level = logrus.WarnLevel
	case "error":
		level = logrus.ErrorLevel
	case "fatal":
		level = logrus.FatalLevel
	case "panic":
		level = logrus.PanicLevel
	default:
		level = logrus.InfoLevel
	}

	formatter := &logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	}

	log := &logrus.Logger{
		Out:       output,
		Formatter: formatter,
		Hooks:     make(logrus.LevelHooks),
		Level:     level,
	}

	if logOutput == "syslog" || logOutput == "test" {
		log.Hooks.Add(hook)
	}

	return log
}
