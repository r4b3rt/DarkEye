package securitytrails

import (
	"fmt"
	"github.com/zsdevX/DarkEye/common"
	"testing"
)

func Test_get(t *testing.T) {
	s := SecurityTrails{
		ApiKey:    "******************",
		Queries:   "52pojie.cn",
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
