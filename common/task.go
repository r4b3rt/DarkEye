package common

import (
	"context"
	"github.com/elastic/beats/libbeat/common/atomic"
	"sync"
)

//Task add comment
type Task struct {
	wait   sync.WaitGroup
	jobs   chan int
	Ctx    context.Context
	cancel context.CancelFunc
	close atomic.Bool
	quit   chan int
	sync.RWMutex
}

//NewTask add comment
func NewTask(jobs int, parent context.Context) *Task {
	ctx, cancel := context.WithCancel(parent)
	t := &Task{
		jobs:   make(chan int, jobs),
		cancel: cancel,
		Ctx:    ctx,
		quit:   make(chan int, 1),
	}
	return t
}

//GetJob return true if success, else stop
func (t *Task) Job() bool {
	//task已经关闭申请失败
	if t.close.Load() {
		return false
	}
	//Parent要求关闭
	select {
	case <-t.Ctx.Done():
		t.Die()
		return false
	default:
	}
	//申请任务
	t.jobs <- 1
	t.wait.Add(1)
	return true
}

func (t *Task) UnJob() {
	<-t.jobs
	t.wait.Done()
	if len(t.jobs) == 0 {
		t.Die()
	}
}

func (t *Task) Wait(f string) {
	//fmt.Println(f)
	select {
	case <-t.quit:
		t.cancel()
	}
	t.wait.Wait()
}

func (t *Task) Die() {
	t.Lock()
	defer t.Unlock()
	if !t.close.Load() {
		t.close.Store(true)
		t.quit<-1
	}
}
