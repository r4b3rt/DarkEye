package common

import (
	"fmt"
	"time"
)

const (
	INFO  = 1
	FAULT = 2
	ALERT = 3
)

var (
	Console = false
)

func LogBuild(module string, logCt string, level int) string {
	return fmt.Sprintf("%s: %s [%s] %s",
		time.Now().Format("2006/1/2 15:04:05"),
		logDesc[level],
		module,
		logCt)
}
