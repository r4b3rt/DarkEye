package common

import (
	"fmt"
	"strings"
	"testing"
)

func Test_getIPRange(t *testing.T) {
	base, start, end, err := GetIPRange("39.98.122.200-39.98.122.230")
	fmt.Println(start, end, fmt.Sprintf("%v", err))
	for {
		nip := GenIP(base, start)
		if strings.Compare(nip, end) > 0 {
			break
		}
		fmt.Println(nip)
		start++
	}
}
