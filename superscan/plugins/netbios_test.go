package plugins

import (
	"github.com/fatih/color"
	"golang.org/x/time/rate"
	"testing"
	"time"
)

func Test_nbCheck(t *testing.T) {
	plg := Plugins{
		TargetIp: "127.0.0.2",
		RateWait: func(r *rate.Limiter) {
			return
		},
	}
	nbCheck(&plg)
	color.Red("%v", plg)
	time.Sleep(time.Second * 5)
}
