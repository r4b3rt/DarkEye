package common

import (
	"github.com/sirupsen/logrus"
)

const (
	//INFO add comment
	INFO = 1
	//FAULT add comment
	FAULT = 2
	//ALERT add comment
	ALERT = 3
)

//Log add comment
func Log(module, logCt string, level int) {
	switch level {
	case INFO:
		logrus.Info(module, ":", logCt)
	case ALERT:
		logrus.Warn(module, ":", logCt)
	case FAULT:
		logrus.Error(module, ":", logCt)
	}
}

func LogUi(logCt string, ui chan string, level int) {
	switch level {
	case INFO:
		ui <- logCt
	case ALERT:
		ui <- logCt
	case FAULT:
		ui <- logCt
	}
}
