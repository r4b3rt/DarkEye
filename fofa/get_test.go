package fofa

import (
	"fmt"
	"testing"
)

func Test_get(t *testing.T) {
	fofa := NewConfig()
	//fofa.Interval = 10
	fofa.Ip = "121.199.9.246"

	fofa.ErrChannel = make(chan string, 10)

	go fofa.Run()

	for {
		msg := <-fofa.ErrChannel
		fmt.Println(msg)
	}
}
