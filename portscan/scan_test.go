package main

import (
	"fmt"
	"testing"
)

func Test_SpeedTest(t *testing.T) {
	s := New("192.168.1.1")
	if err := s.InitConfig(); err != nil {
		t.Fatal(err.Error())
	}
	fmt.Println(s)
	s.ActivePort = "80"
	s.TimeOut = 300
	s.Run()
}

func Test_Run(t *testing.T) {
	*mIp = "192.167.1.1-254"
	*mThread = 10
	Start()
}
