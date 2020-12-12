package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
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
	*mIp = "192.168.1.1-254"
	*mThread = 30
	*mTimeOut = 2000
	//*mPortList = "53"
	//*mPortList ="1-65535"
	go func() {
		log.Println(http.ListenAndServe("localhost:10000", nil))
	}()
	Start()
}
