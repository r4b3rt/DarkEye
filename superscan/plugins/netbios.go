package plugins

import (
	"bytes"
	"encoding/hex"
	"net"
	"strings"
	"time"
)

//抄的"github.com/shadow1ng/fscan/common"，很棒的项目！
var (
	bufferV1, _ = hex.DecodeString("05000b03100000004800000001000000b810b810000000000100000000000100c4fefc9960521b10bbcb00aa0021347a00000000045d888aeb1cc9119fe808002b10486002000000")
	bufferV2, _ = hex.DecodeString("050000031000000018000000010000000000000000000500")
	bufferV3, _ = hex.DecodeString("0900ffff0000")
)

func nbCheck(plg *Plugins, f *funcDesc) {
	plg.TargetPort = f.port
	nbConn(plg, f)
}

func nbConn(plg *Plugins, f *funcDesc) {
	conn, err := net.DialTimeout("tcp",
		plg.TargetIp+":"+plg.TargetPort, time.Duration(plg.TimeOut)*time.Millisecond)
	if err != nil {
		return
	}
	_ = conn.SetDeadline(time.Now().Add(time.Duration(plg.TimeOut) * time.Millisecond))
	defer conn.Close()
	_, _ = conn.Write(bufferV1)
	reply := make([]byte, 4096)
	_, err = conn.Read(reply)
	if err != nil {
		return
	}
	_, _ = conn.Write(bufferV2)
	if n, err := conn.Read(reply); err != nil || n < 42 {
		return
	}
	text := reply[42:]

	for i := 0; i < len(text)-5; i++ {
		if bytes.Equal(text[i:i+6], bufferV3) {
			text = text[:i-4]
			collectNbi(text, &plg.NetBios)
			plg.TargetProtocol = f.name
			plg.PortOpened = true
			return
		}
	}
}
func collectNbi(text []byte, nb *NetBios) {
	encodedStr := hex.EncodeToString(text)
	hosts := strings.Replace(encodedStr, "0700", "", -1)
	hostname := strings.Split(hosts, "000000")
	for i := 0; i < len(hostname); i++ {
		hostname[i] = strings.Replace(hostname[i], "00", "", -1)
		host, err := hex.DecodeString(hostname[i])
		if err != nil {
			return
		}
		nb.Ip = nb.Ip + "," + string(host)
	}
	nb.Ip = strings.TrimPrefix(nb.Ip, ",")
	return
}
