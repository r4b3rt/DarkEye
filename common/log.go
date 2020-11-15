package common

import (
	"fmt"
)

const (
	INFO  = 1
	FAULT = 2
	ALERT = 3
)

var (
	Console = false
	logDesc = []string{
		0:     "None",
		INFO:  `[!]`,
		FAULT: `<font color="red">[x]</font>`,
		ALERT: `<font color="green">[√]</font>`,
	}
	logFile = "dark_eye.log"
)

func LogBuild(module string, logCt string, level int) string {
	return fmt.Sprintf("%s /%s/ %s",
		logDesc[level],
		module,
		logCt)
}
