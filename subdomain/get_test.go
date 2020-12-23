package subdomain

import (
	"fmt"
	"github.com/zsdevX/DarkEye/common"
	"testing"
	"time"
)

func Test_get(t *testing.T) {
	s := SubDomain{
		ApiKey:    "**************",
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

func Test_brute(t *testing.T) {
	s := SubDomain{
		ApiKey:      "**************",
		Queries:     "baidu.com",
		DnsServer:   "192.168.1.1:53",
		IpCheck:     true,
		Brute:       true,
		BruteLength: "3",
		BruteRate:   "100",
	}
	common.Console = true

	s.ErrChannel = make(chan string, 10)

	go s.Run()

	for {
		msg := <-s.ErrChannel
		fmt.Println(msg)
	}
}

func Test_IPAPIRate(t *testing.T) {
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

func Test_GenBruteSource(t *testing.T) {
	target := make([]string, 0)
	length := 3
	i := 0
	for length > 0 {
		target = genSource(target, common.LowCaseAlpha+"0123456789", i+1)
		length--
		i++
	}
	if len(target) != 36+36*36+36*36*36 {
		t.Fail()
	}
}
