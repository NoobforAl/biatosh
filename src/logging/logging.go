package logging

import (
	"biatosh/contract"
	"os"
	"sync"

	"github.com/sirupsen/logrus"
)

var once sync.Once
var logger *logrus.Logger

func New() contract.Logger {
	once.Do(
		func() {
			level := logrus.InfoLevel

			/* if config.Debug {
				level = logrus.DebugLevel
			} */

			logger = &logrus.Logger{
				Out:          os.Stderr,
				Formatter:    new(logrus.TextFormatter),
				Hooks:        make(logrus.LevelHooks),
				Level:        level,
				ExitFunc:     os.Exit,
				ReportCaller: false,
			}
		})

	return logger
}
