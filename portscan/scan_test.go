package main

import (
	"fmt"
	"testing"
)

func Test_SpeedTest(t *testing.T) {
	s := New("192.168.1.1")
	fmt.Println(s)
	s.ActivePort = "80"
	s.TimeOut = 300
	s.Run()
}

func Test_Run(t *testing.T) {
	*mIp = "192.168.1.3"
	*mThread = 3
	*mTimeOut = 200
	Start()
}
