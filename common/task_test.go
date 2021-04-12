package common

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func Test_task(t *testing.T) {
	task := NewTask(2000*2, context.TODO())
	defer task.Wait()
	i := 0
	for {
		task.Job()
		go func(idx int) {
			defer task.UnJob()
			fmt.Println("hi", idx)
			time.Sleep(time.Second)
			return
		}(i)
		i++
	}

}
