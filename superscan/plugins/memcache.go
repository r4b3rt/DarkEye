package plugins

import (
	"fmt"
	"net"
	"strings"
	"time"
)

func init() {
	checkFuncs[MemoryCacheSrv] = memcachedCheck
	supportPlugin["memcached"] = "memcached"
}

func memcachedCheck(plg *Plugins) {
	if !plg.NoTrust && plg.TargetPort != "11211" {
		return
	}
	if memcachedUnAuth(plg) {
		plg.Cracked = append(plg.Cracked, Account{Username: "空", Password: "空"})
		plg.PortOpened = true
		plg.highLight = true
		plg.TargetProtocol = "memcached"
		return
	}
}

func memcachedUnAuth(plg *Plugins) (ok bool) {
	conn, err := net.DialTimeout("tcp",
		fmt.Sprintf("%s:%s", plg.TargetIp, plg.TargetPort), time.Duration(plg.TimeOut)*time.Millisecond)
	if err != nil {
		return
	}
	defer conn.Close()
	_ = conn.SetDeadline(time.Now().Add(time.Duration(plg.TimeOut) * time.Millisecond))
	_, _ = conn.Write([]byte("stats\n"))
	buff := make([]byte, 1024)
	n, err := conn.Read(buff)
	if err != nil {
		return
	}
	if strings.Contains(string(buff[:n]), "STAT") {
		/*EXP:
		apt install libmemcached-tools
		memcdump --servers=192.168.1.33 (列出key）
		memccat --servers=192.168.1.33 key1 （列出key1内容）
		*/
		return true
	}
	return
}
