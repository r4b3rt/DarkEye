package zoomeye

import (
	"fmt"
	"testing"
	"time"
)

func Test_zoomEye(t *testing.T) {
	z := New()
	z.Query = "ip:27.17.15.195"
	z.ApiKey = "21540Bf6-6EAE-7fD2E-369c-96B6Cd3109c"
	go z.Run()

	for {
		msg := <-z.ErrChannel
		fmt.Println(msg)
	}
	time.Sleep(10 * time.Second)
}
