package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)


func Test_portSplit(t *testing.T) {
	x := portSplit("81")
	assert.Equal(t, "81", x[0])
	assert.Equal(t, 1, len(x))
	x = portSplit("81,82")
	assert.Equal(t, 2, len(x))
	x = portSplit("81,82-84,85")
	assert.Equal(t, 5, len(x))
	assert.Equal(t, "83", x[2])
	x = portSplit("4430-4445")
	assert.Equal(t, 4445-4430+1, len(x))
}

func Test_splitIp2C(t *testing.T) {
	x, _ := splitIp2C("1.1.1.1")
	assert.Equal(t, 1, len(x))

	x, _ = splitIp2C("1.1.1.1-24")
	assert.Equal(t, 1, len(x))

	x, _ = splitIp2C("1.1.1-24")
	assert.Equal(t, 24, len(x))

	x, _ = splitIp2C("1.1-24")
	assert.Equal(t, 24*256, len(x))

	x, _ = splitIp2C("1.1.1.1/16")
	assert.Equal(t, 256, len(x))
}

func Test_splitIpC2Ip(t *testing.T) {
	x, _ := splitIpC2Ip("192.168.1.1-2")
	assert.Equal(t, 2, len(x))
	assert.Equal(t, "192.168.1.2", x[1].String())

	x, _ = splitIpC2Ip("192.168.1.1/24")
	assert.Equal(t, 254, len(x))
	assert.Equal(t, "192.168.1.1", x[0].String())
	assert.Equal(t, "192.168.1.254", x[253].String())

	x, _ = splitIpC2Ip("192.168.1.1/30")
	assert.Equal(t, 2, len(x))
	assert.Equal(t, "192.168.1.2", x[1].String())

	x, _ = splitIpC2Ip("192.168.1.60/30")
	assert.Equal(t, 2, len(x))
	assert.Equal(t, "192.168.1.62", x[1].String())

	x, _ = splitIpC2Ip("192.168.1.1")
	assert.Equal(t, 1, len(x))
}
