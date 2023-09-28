package misc

import (
	"os"

	"github.com/sirupsen/logrus"
)

func init() {
	level, err := logrus.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		level = logrus.WarnLevel
	}
	logrus.SetLevel(level)
	logrus.SetReportCaller(true)
}
