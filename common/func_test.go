package common

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

//toTest_getIPRange do add comment
func Test_getIPRange(t *testing.T) {
	base, start, end, err := GetIPRange("39.98.122.200-39.98.122.230")
	assert.Equal(t, "39.98.122.200", base)
	assert.Equal(t, 0, start)
	assert.Equal(t, "39.98.122.230", end)
	assert.Equal(t, err, nil)

	nip := GenIP(base, start)
	assert.Equal(t, base, nip)
}
