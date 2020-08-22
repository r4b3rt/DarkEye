package fofa

import (
	"github.com/zsdevX/DarkEye/common"
	"testing"
	"fmt"
)

func Test_get(t *testing.T) {
	fofa := NewConfig()
	fofa.Ip = "1.1.1.1"
	common.Console = true

	fofa.ErrChannel = make(chan string, 10)

	go fofa.Run()

	for {
		msg := <-fofa.ErrChannel
		fmt.Println(msg)
	}
}








