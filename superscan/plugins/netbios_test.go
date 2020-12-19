package plugins

import (
	"github.com/fatih/color"
	"testing"
	"time"
)

func Test_nbCheck(t *testing.T) {
	plg := Plugins{
		TargetIp: "127.0.0.2",
	}
	nbCheck(&plg)
	color.Red("%v", plg)
	time.Sleep(time.Second * 5)
}
