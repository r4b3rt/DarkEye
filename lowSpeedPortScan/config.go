package main

import (
	"encoding/json"
	"github.com/zsdevX/DarkEye/common"
	"io/ioutil"
	"os"
	"path/filepath"
)

type FromTo struct {
	from int
	to   int
}

type Scan struct {
	Ip          string `json:"ip"`
	PortRange   string `json:"port_range"`
	ActivePort  string `json:"active_port"`
	DefaultRate int    `json:"default_rate"`
	MinRate     int    `json:"min_rate"`
	Test        bool   `json:"rate_test"`

	Rate                 int          `json:"rate"`
	PortsHaveBeenScanned map[int]bool `json:"port_scanned"`
	PortsScannedOpened   []int        `json:"ports_opened"`

	valid bool
}

var (
	mConfigFile = "scan.cfg"
	mBasedir    = filepath.Join(common.BaseDir, "tmp")
)

func init() {
	_ = os.Mkdir(mBasedir, 0700)
	mConfigFile = filepath.Join(mBasedir, mConfigFile)
}

func (s *Scan) loadCfg() error {
	data, err := ioutil.ReadFile(mConfigFile + "." + s.Ip)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, s); err != nil {
		return err
	}
	s.valid = true
	return nil
}

func (s *Scan) saveCfg() error {
	data, err := json.Marshal(s)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(mConfigFile + "." + s.Ip, data, 0700)
}
