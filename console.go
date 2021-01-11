package gcplog

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

type ConsoleFormatter struct {
	ProjectID string
}

func (f *ConsoleFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var latency time.Duration

	if l, ok := entry.Data["latency"]; ok {
		latency, _ = l.(time.Duration)
	}

	s := fmt.Sprintf("%-10s: %10s %s\n", getGCPLevel(entry.Level), latency, entry.Message)
	return []byte(s), nil
}
