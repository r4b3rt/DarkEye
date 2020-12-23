package spider

import (
	"fmt"
	"testing"
	"time"
)

func Test_Spider(t *testing.T) {
	sp := NewConfig()
	sp.ErrChannel = make(chan string, 10)
	sp.Url = "http://www.baidu.com"
	sp.MaxDeps = 2

	go func() {
		for {
			msg := <-sp.ErrChannel
			fmt.Println(msg)
		}
	}()

	sp.Run()
	time.Sleep(time.Second * 10)
}

func Test_Search(t *testing.T) {
}
