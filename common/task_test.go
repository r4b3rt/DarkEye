package common

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_task(t *testing.T) {
	ctx, cancel := context.WithCancel(context.TODO())
	time.AfterFunc(time.Second*10, func() {
		t.Log("called cancel")
		cancel()
	})
	jobs := 100
	task := NewTask(jobs, ctx)
	defer task.Wait("test")
	for {
		if !task.Job() {
			break
		}
		assert.LessOrEqual(t, len(task.jobs), 100)
		go func() {
			defer task.UnJob()
			time.Sleep(time.Millisecond * 5)
			return
		}()
	}
}
