package zoomeye

import (
	"fmt"
	"testing"
)

func Test_zoomEye(t *testing.T) {
	z := New()
	z.Query = "title:aaa"
	z.ApiKey = "563e7683-2B49-03589-4655-377840a766a"
	z.Pages = "3"
	z.ErrChannel = make(chan string, 10)
	go z.Run()

	for {
		msg := <-z.ErrChannel
		fmt.Println(msg)
	}
}
