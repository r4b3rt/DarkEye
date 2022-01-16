package scan

import (
	"context"
	"net"
	"time"
)

func dail(parent context.Context, protocol, addr string, timeout int) (net.Conn, error) {
	timeOut := time.Millisecond * time.Duration(timeout)
	d := net.Dialer{Timeout: timeOut}
	ctx, _ := context.WithTimeout(parent, timeOut)
	return d.DialContext(ctx, protocol, addr)
}

func hello(parent context.Context, protocol, addr string, hi []byte, timeout int) ([]byte, error) {
	c, err := dail(parent, protocol, addr, timeout)
	if err != nil {
		return nil, err
	}
	defer c.Close()
	if hi != nil {
		_, _ = c.Write(hi)
	}
	_ = c.SetReadDeadline(time.Now().Add(time.Millisecond * time.Duration(timeout)))
	buf := make([]byte, 1024)
	n, err := c.Read(buf)
	if err != nil {
		return nil, err
	}
	return buf[:n], nil
}
