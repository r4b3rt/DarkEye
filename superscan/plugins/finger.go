package plugins

import (
	"bytes"
	"context"
	"fmt"
	"github.com/zsdevX/DarkEye/common"
	"net"
	"time"
)

func (plg *Plugins) finger() {
	conn, err := common.DialCtx(context.Background(), "tcp",
		net.JoinHostPort(plg.TargetIp, plg.TargetPort), time.Duration(Config.TimeOut)*time.Millisecond)
	if err != nil {
		return
	}
	defer conn.Close()
	_ = conn.SetDeadline(time.Now().Add(time.Duration(Config.TimeOut) * time.Millisecond))
	_, _ = conn.Write([]byte("asd\n"))
	buff := make([]byte, 1024)
	n, err := conn.Read(buff)
	if err != nil {
		return
	}
	if n > 64 {
		n = 32
	}
	out := bytes.Buffer{}
	for _, v := range buff[:n] {
		if v >= 0x20 && v <= 0x7e {
			out.WriteString(string(v))
		} else {
			out.Write([]byte(fmt.Sprintf("\\x%X", v)))
		}
	}
	plg.Result.Output.Set("finger", out.String())
	return
}
