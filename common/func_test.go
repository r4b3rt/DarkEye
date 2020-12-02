package common

import (
	"fmt"
	"testing"
)

func Test_getIPRange(t *testing.T) {
	base, start, end, _ := GetIPRange("192.168.10.1-255")
	for start < end {
		fmt.Println(GenIP(base, start))
		start++
	}
}
