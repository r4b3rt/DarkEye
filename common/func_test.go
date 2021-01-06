package common

import (
	"fmt"
	"testing"
)

//toTest_getIPRange do add comment
func Test_getIPRange(t *testing.T) {
	base, start, end, err := GetIPRange("39.98.122.200-39.98.122.230")
	fmt.Println(start, end, fmt.Sprintf("%v", err))
	for {
		nip := GenIP(base, start)
		if CompareIP(nip, end) > 0 {
			break
		}
	}
}
