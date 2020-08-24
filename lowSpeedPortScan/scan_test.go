package main

import (
	"fmt"
	"testing"
)

func Test_SpeedTest(t *testing.T) {
	if err := initConfig(); err != nil {
		t.Fatal(err.Error())
	}
	s := &scanCfg
	fmt.Println(s)
	s.Ip = "101.231.113.1310"
	s.ActivePort = "8443"
	s.Speed = 300
	//if err := s.SpeedTest(); err != nil {
	//	t.Fatal(err.Error())
	//}
	s.Run()
}
