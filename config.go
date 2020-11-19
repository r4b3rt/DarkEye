package main

import (
	"encoding/json"
	"github.com/zsdevX/DarkEye/common"
	"github.com/zsdevX/DarkEye/fofa"
	"github.com/zsdevX/DarkEye/poc"
	"github.com/zsdevX/DarkEye/securitytrails"
	"github.com/zsdevX/DarkEye/spider"
	"io/ioutil"
	"path/filepath"
)

var (
	mConfigFile = "dark_eye.cfg"
	mConfig     = config{
		Valid: false,
	}
)

func init() {
	mConfigFile = filepath.Join(common.BaseDir, mConfigFile)
}

type config struct {
	SecurityTrails securitytrails.SecurityTrails `json:"securitytrails"`
	Fofa           fofa.Fofa                     `json:"fofa"`
	Spider         spider.Spider                 `json:"spider"`
	Poc            poc.Poc                       `json:"poc"`
	Valid          bool                          `json:"valid"`
}

func loadCfg() error {
	defer func() {
		//首次使用配置不存在, 初始化默认值
		if !mConfig.Valid {
			mConfig.Valid = true
			mConfig.Fofa = fofa.NewConfig()
			mConfig.SecurityTrails = securitytrails.NewConfig()
			mConfig.Spider = spider.NewConfig()
			mConfig.Poc = poc.NewConfig()
		}
	}()
	data, err := ioutil.ReadFile(mConfigFile)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, &mConfig); err != nil {
		return err
	}
	return nil
}

func saveCfg() error {
	data, err := json.Marshal(mConfig)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(mConfigFile, data, 0700)
}
