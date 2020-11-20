package securitytrails

import (
	"fmt"
	"github.com/zsdevX/DarkEye/common"
	"testing"
	"time"
)

func Test_get(t *testing.T) {
	s := SecurityTrails{
		ApiKey:    "v94C1s0xgSR21tbSJsOV9G5rk6vpMuf3",
		Queries:   "baidu.com",
		DnsServer: "192.168.1.1:53",
		IpCheck:   true,
	}
	common.Console = true

	s.ErrChannel = make(chan string, 10)

	go s.Run()

	for {
		msg := <-s.ErrChannel
		fmt.Println(msg)
	}
}

//go test
func Test_parseTag(t *testing.T) {
	d := dnsInfo{
		domain: "ooxx.com",
	}
	s := SecurityTrails{
	}
	s.parseTag(&d)
	fmt.Println(d)
}

func Test_ipapiRate(t *testing.T) {
	i := 0
	for {
		if ipApiLimit.Allow() {
			fmt.Println("fuck it", i)
			i++
		} else {
			fmt.Println("limit 1 second")
			time.Sleep(time.Second * 10)
		}
	}
}
