package zoomeye

import (
	"fmt"
	"testing"
)

func Test_zoomEye(t *testing.T) {
	z := New()
	z.Query = "title:aaa"
	z.ApiKey = "DF885877-5e59-e5de4-f4df-1ad9cea4e5b"
	z.Pages = "1"
	z.ErrChannel = make(chan string, 10)
	go z.Run()

	for {
		msg := <-z.ErrChannel
		fmt.Println(msg)
	}
}
