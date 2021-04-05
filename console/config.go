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
	//分析模块
	ModuleFuncs[moduleId(zoomEye)] = ModuleFunc{
		name:        zoomEye,
		start:       zoomEyeRuntimeRuntimeOptions.start,
		compileArgs: zoomEyeRuntimeRuntimeOptions.compileArgs,
		usage:       zoomEyeRuntimeRuntimeOptions.usage,
		init:        zoomEyeInitRunTime,
		valueCheck:  zoomEyeValueCheck,
		completer:   mContext.zoomEyeArgumentsCompleter,
	}
}

func moduleId(m string) string {
	return strings.ToLower(m)
}
