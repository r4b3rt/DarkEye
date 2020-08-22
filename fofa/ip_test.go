package fofa

import (
	"github.com/zsdevX/DarkEye/common"
	"testing"
	"fmt"
)

func Test_ip(t *testing.T) {
	fofa := NewConfig()
	fofa.Ip = "101.231.113.13"
	common.Console = true

	fofa.ErrChannel = make(chan string, 10)

	go fofa.Run()

	for {
		msg := <-fofa.ErrChannel
		fmt.Println(msg)
	}
}
