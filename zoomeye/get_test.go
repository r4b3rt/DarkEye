package zoomeye

import (
	"context"
	"fmt"
	"testing"
)

func Test_run(t *testing.T) {
	z := ZoomEye{
		Query:"telnet",
		ApiKey:"",
	}
	fmt.Println(z.run(context.TODO(), 1))
}
