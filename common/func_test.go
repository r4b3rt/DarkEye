package common

import (
	"fmt"
	"testing"
)

func Test_getIPRange(t *testing.T) {
	base, start, end, err := GetIPRange("192.168.10.3-9")
	fmt.Println(fmt.Sprintf("%v", err))
	for start < end {
		fmt.Println(GenIP(base, start))
		start++
	}
	fmt.Println(fmt.Sprintf("%20s", "fuck"))
}
