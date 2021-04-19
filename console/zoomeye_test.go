package main

import (
	"fmt"
	"testing"
)

func Test_zoomEyeQuery(t *testing.T) {
	cmd := []string{
		"country: india",
		"+",
		"service: redis",
		"+",
		"port:6379",
	}
	z := zoomEyeRuntime{}
	_, s := z.buildQuery(cmd)
	if s != "country:india +service:redis +port:6379" {
		fmt.Println(s)
		t.Fail()
	}
}
