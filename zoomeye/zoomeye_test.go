package zoomeye

import (
	"fmt"
	"testing"
)

func Test_zoomEye(t *testing.T) {
	z := New()
	z.Query = "title:aaa"
	z.ApiKey = "9CFD584F-9701-3EF03-3653-e6bda60993f"
	z.Pages = "1"
	z.ErrChannel = make(chan string, 10)
	go z.Run()

	for {
		msg := <-z.ErrChannel
		fmt.Println(msg)
	}
}
