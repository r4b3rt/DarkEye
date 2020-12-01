package main

import (
	"github.com/zsdevX/DarkEye/hacktools/portscan/plugins"
	"sync"
)

type Scan struct {
	//需配置参数
	Ip                  string `json:"ip"`
	PortRange           string `json:"port_range"`
	ActivePort          string `json:"active_port"`
	PortRangeThresholds int    `json:"port_range_thresholds"`
	ThreadNumber        int    `json:"thread_number"`
	//任务执行结果
	TimeOut              int               `json:"timeout"`
	PortsHaveBeenScanned map[int]bool      `json:"port_scanned"`
	PortsScannedOpened   []plugins.Plugins `json:"ports_opened"`
	//用于回显示
	Callback    func([]byte) `json:"-"`
	BarCallback func(i int)  `json:"-"`
	lock        sync.RWMutex
}
