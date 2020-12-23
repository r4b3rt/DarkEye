package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	"testing"
)

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

func Test_pingNet(t *testing.T) {
	s := newScan("")
	s.PingNet("192.168.1.1-192.168.255.255")
}
