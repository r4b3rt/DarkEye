package main

import (
	"github.com/zsdevX/DarkEye/common"
	"os"
	"path/filepath"
)

type PortInfo struct {
	Port   int
	Server string
	Title  string
}

type Scan struct {
	//需配置参数
	Ip             string `json:"ip"`
	PortRange      string `json:"port_range"`
	ActivePort     string `json:"active_port"`
	DefaultTimeOut int    `json:"default_timeout"`
	MinTimeOut     int    `json:"min_timeout"`
	Test           bool   `json:"rate_test"`
	//任务执行结果
	TimeOut              int          `json:"timeout"`
	PortsHaveBeenScanned map[int]bool `json:"port_scanned"`
	PortsScannedOpened   []PortInfo   `json:"ports_opened"`
	Title                bool         `json:"title"`
	//用于回显示
	Callback    func(a ...interface{}) `json:"-"`
	BarCallback func()                 `json:"-"`
}

var (
	mBasedir = filepath.Join(common.BaseDir, "tmp")
)

func init() {
	_ = os.Mkdir(mBasedir, 0700)
}
