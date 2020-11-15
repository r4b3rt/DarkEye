package common

import (
	"fmt"
	"github.com/fatih/color"
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
		INFO:  "[!]",
		FAULT: color.HiRedString("[x]"),
		ALERT: color.HiGreenString("[√]"),
	}
	logFile = "dark_eye.log"
)

func LogBuild(module string, logCt string, level int) string {
	return fmt.Sprintf("%s /%s/ %s",
		logDesc[level],
		module,
		logCt)
}
