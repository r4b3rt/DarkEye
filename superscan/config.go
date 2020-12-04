package main

import (
	"github.com/zsdevX/DarkEye/superscan/plugins"
	"golang.org/x/time/rate"
	"sync"
)

type Scan struct {
	//需配置参数
	Ip           string `json:"ip"`
	PortRange    string `json:"port_range"`
	ActivePort   string `json:"active_port"`
	ThreadNumber int    `json:"thread_number"`
	NoTrust      bool
	PluginWorker int

	//任务执行结果
	Rate               *rate.Limiter
	TimeOut            int               `json:"timeout"`
	PortsScannedOpened []plugins.Plugins `json:"ports_opened"`
	//用于回显示
	Callback               func([]byte)   `json:"-"`
	BarCallback            func(i int)    `json:"-"`
	BarDescriptionCallback func(i string) `json:"-"`
	lock                   sync.RWMutex
}
