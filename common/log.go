package common

import (
	"fmt"
	"github.com/fatih/color"
)

const (
	//INFO add comment
	INFO = 1
	//FAULT add comment
	FAULT = 2
	//ALERT add comment
	ALERT = 3
)

//LogBuild add comment
func Log(module, logCt string, level int) {
	switch level {
	case INFO:
		color.Green(fmt.Sprintf("[√] %s %s",
			module,
			logCt))
	case ALERT:
		color.Yellow(fmt.Sprintf("[!] %s %s",
			module,
			logCt))
	case FAULT:
		color.Yellow(fmt.Sprintf("[x] %s %s",
			module,
			logCt))
	}
}
