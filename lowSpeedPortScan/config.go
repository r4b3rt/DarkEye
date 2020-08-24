package main

import (
	"encoding/json"
	"github.com/zsdevX/DarkEye/common"
	"io/ioutil"
	"path/filepath"
)

type FromTo struct {
	from int
	to   int
}

type Scan struct {
	Ip           string `json:"ip"`
	PortRange    string `json:"port_range"`
	ActivePort   string `json:"active_port"`
	DefaultSpeed int    `json:"default_speed"`
	MinSpeed     int    `json:"min_speed"`
	Test         bool   `json:"speed_test"`

	Speed                int          `json:"speed"`
	PortsHaveBeenScanned map[int]bool `json:"port_scanned"`
	PortsScannedOpened   []int        `json:"ports_opened"`

	valid  bool
}

var (
	mConfigFile = "lowspeedportscan.cfg"
)

func init() {
	mConfigFile = filepath.Join(common.BaseDir, mConfigFile)
}

func loadCfg() error {
	data, err := ioutil.ReadFile(mConfigFile)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, &scanCfg); err != nil {
		return err
	}
	scanCfg.valid = true
	return nil
}

func saveCfg() error {
	data, err := json.Marshal(&scanCfg)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(mConfigFile, data, 0700)
}
