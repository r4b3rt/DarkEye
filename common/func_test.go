package common

import (
	"fmt"
	"testing"
)

func Test_getIPRange(t *testing.T) {
	base, start, end, err := GetIPRange("192.168.10.1-255/")
	fmt.Println(fmt.Sprintf("%v", err))
	for start < end {
		fmt.Println(GenIP(base, start))
		start++
	}
}
