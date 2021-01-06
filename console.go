package gcplog

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

type ConsoleFormatter struct {
	ProjectID string
}

func (f *ConsoleFormatter) Format(entry *logrus.Entry) ([]byte, error) {

	s := fmt.Sprintf("%-10s: %s\n", getGCPLevel(entry.Level), entry.Message)
	return []byte(s), nil
}
