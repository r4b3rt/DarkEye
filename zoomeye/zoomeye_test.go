package zoomeye

import (
	"fmt"
	"testing"
)

func Test_zoomEye(t *testing.T) {
	z := New()
	z.Query = "title:aaa"
	z.ApiKey = "api"
	z.Pages = 1
	z.ErrChannel = make(chan string, 10)
	go func() {
		for {
			msg := <-z.ErrChannel
			fmt.Println(msg)
		}
	}()
	z.Run()

}
