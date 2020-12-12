package main

import (
	"encoding/json"
	"github.com/zsdevX/DarkEye/common"
	"github.com/zsdevX/DarkEye/fofa"
	"github.com/zsdevX/DarkEye/spider"
	"github.com/zsdevX/DarkEye/subdomain"
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
	SubDomain subdomain.SubDomain `json:"sub_domain"`
	Fofa           fofa.Fofa                     `json:"fofa"`
	Spider         spider.Spider                 `json:"spider"`
	Valid          bool                          `json:"valid"`
}

func loadCfg() error {
	defer func() {
		//首次使用配置不存在, 初始化默认值
		if !mConfig.Valid {
			mConfig.Valid = true
			mConfig.Fofa = fofa.NewConfig()
			mConfig.SubDomain = subdomain.NewConfig()
			mConfig.Spider = spider.NewConfig()
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
