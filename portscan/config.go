package main

import (
	"github.com/zsdevX/DarkEye/common"
	"os"
	"path/filepath"
)

type FromTo struct {
	from int
	to   int
}

type PortInfo struct {
	Port   int
	Server string
	Title  string
}

type Scan struct {
	Ip             string `json:"ip"`
	PortRange      string `json:"port_range"`
	ActivePort     string `json:"active_port"`
	DefaultTimeOut int    `json:"default_timeout"`
	MinTimeOut     int    `json:"min_timeout"`
	Test           bool   `json:"rate_test"`

	TimeOut              int          `json:"timeout"`
	PortsHaveBeenScanned map[int]bool `json:"port_scanned"`
	PortsScannedOpened   []PortInfo   `json:"ports_opened"`
	Title                bool         `json:"title"`

	valid bool
}

var (
	mBasedir = filepath.Join(common.BaseDir, "tmp")
)

func init() {
	_ = os.Mkdir(mBasedir, 0700)
}
