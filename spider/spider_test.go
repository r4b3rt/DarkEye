package spider

import (
	"fmt"
	"testing"
)

func Test_Spider(t *testing.T) {
	sp := NewConfig()
	sp.ErrChannel = make(chan string, 10)
	sp.Url = "http://www.varbing.com"
	sp.MaxDeps = 2

	go func() {
		for {
			msg := <-sp.ErrChannel
			fmt.Println(msg)
		}
	}()

	sp.Run()

}
