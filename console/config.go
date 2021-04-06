package main

import (
	"context"
	"github.com/c-bata/go-prompt"
	"github.com/elastic/beats/libbeat/common/atomic"
	"strings"
)

//RequestContext add comment
type RequestContext struct {
	CmdArgs     []string
	ctx         context.Context
	cancel      context.CancelFunc
	running     atomic.Bool
	Interactive bool
}

type analysisEntity struct {
	ID      int64  `json:"id" gorm:"primaryKey"`
	Task    string `json:"task" gorm:"unique_index:UNIQ_hi;column:task"`
	Ip      string `json:"ip" gorm:"unique_index:UNIQ_hi;column:ip"`
	Port    string `json:"port" gorm:"unique_index:UNIQ_hi;column:port"`
	Service string `json:"service" gorm:"unique_index:UNIQ_hi;column:service"`

	Url             string `json:"url" gorm:"column:url"`
	Title           string `json:"title" gorm:"column:title"`
	WebServer       string `json:"web_server" gorm:"column:web_server"`
	WebResponseCode int32  `json:"http_code" gorm:"column:http_code"`

	Hostname  string
	Os        string
	Device    string
	Banner    string
	Version   string
	ExtraInfo string
	RDns      string
	Country   string

	NetBios     string `json:"netbios" gorm:"column:netbios"`
	WeakAccount string `json:"weak_account" gorm:"column:weak_account"`
	Vulnerable  string `json:"vulnerable" gorm:"column:vulnerable"`
}

//ModuleFunc add comment
type ModuleFunc struct {
	name        string
	start       func(ctx context.Context)
	init        func()
	valueCheck  map[string]bool
	completer   func(args []string) []prompt.Suggest
	compileArgs func(args []string) error
	usage       func()
}

var (
	//ModuleFuncs add comment
	ModuleFuncs = make(map[string]ModuleFunc, 0)
)

func init() {
	//端口、弱口令扫描模块
	ModuleFuncs[moduleId(superScan)] = ModuleFunc{
		name:        superScan,
		start:       superScanRuntimeOptions.start,
		init:        superScanInitRunTime,
		compileArgs: superScanRuntimeOptions.compileArgs,
		usage:       superScanRuntimeOptions.usage,
		valueCheck:  superScanValueCheck,
		completer:   mContext.superScanArgumentsCompleter,
	}
	//分析模块
	ModuleFuncs[moduleId(analysisProgram)] = ModuleFunc{
		name:        analysisProgram,
		start:       analysisRuntimeOptions.start,
		init:        analysisInitRunTime,
		compileArgs: analysisRuntimeOptions.compileArgs,
		usage:       analysisRuntimeOptions.usage,
		valueCheck:  analysisValueCheck,
		completer:   mContext.analysisArgumentsCompleter,
	}
	//资产采集
	ModuleFuncs[moduleId(zoomEye)] = ModuleFunc{
		name:        zoomEye,
		start:       zoomEyeRuntimeOptions.start,
		compileArgs: zoomEyeRuntimeOptions.compileArgs,
		usage:       zoomEyeRuntimeOptions.usage,
		init:        zoomEyeInitRunTime,
		valueCheck:  zoomEyeValueCheck,
		completer:   mContext.zoomEyeArgumentsCompleter,
	}
	//脆弱性检查
	ModuleFuncs[moduleId(xRayProgram)] = ModuleFunc{
		name:        xRayProgram,
		start:       xRayRuntimeOptions.start,
		compileArgs: xRayRuntimeOptions.compileArgs,
		usage:       xRayRuntimeOptions.usage,
		init:        xRayInitRunTime,
		valueCheck:  xRayValueCheck,
		completer:   mContext.xRayArgumentsCompleter,
	}
}

func moduleId(m string) string {
	return strings.ToLower(m)
}
