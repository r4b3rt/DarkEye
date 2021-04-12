package main

import (
	"github.com/c-bata/go-prompt"
)

var (
	superScanSuggestions = []prompt.Suggest{
		// Command
		{"-ip", "Scan ip target"},
		{"-pps", "Crack rate(packet/second)"},
		{"-port-list", "Port list(80-88,8080"},
		{"-timeout", "网络超时(单位ms)"},
		{"-thread", "每个IP爆破端口的线程数量"},
		{"-plugin", "指定协议插件爆破"},
		{"-user-list", "字符串(u1,u2,u3)或文件(一个账号一行)"},
		{"-pass-list", "字符串(p1,p2,p3)或文件（一个密码一行"},
		{"-alive-host-check", "只检查活跃主机的网段(ping)"},
		{"-only-alive-network", "检查所有活跃主机(ping)"},
	}
	superScanValueCheck = map[string]bool{
		"-ip":                 false,
		"-pps":                false,
		"-port-list":          false,
		"-timeout":            false,
		"-thread":             false,
		"-plugin":             false,
		"-user-list":          false,
		"-pass-list":          false,
		"-alive-host-check":   true,
		"-only-alive-network": true,
	}
)

func init() {
	mSuggestions = append(mSuggestions, []prompt.Suggest{
		{superScan, "Live port/service/weak password scan"},
	}...)
}

func (s *superScanRuntime) Completer(args []string) []prompt.Suggest {
	if len(args) <= 1 {
		return filterSuggestions(runCompleteCheck(superScanSuggestions, s.parent.CmdArgs,
			[]string{
				"-ip",
			}), s.parent.CmdArgs)
	}
	//过滤重复的命令
	if len(filterSuggestions([]prompt.Suggest{
		{args[0], "any"},
	}, s.parent.CmdArgs)) == 0 {
		return []prompt.Suggest{}
	}
	switch args[0] {
	case "-ip":
		if len(args) == 2 {
			return []prompt.Suggest{{"192.168.1.1-192.168.1.255", "Scan ip target"},}
		}
	case "-pps":
		if len(args) == 2 {
			return []prompt.Suggest{{"0", "'0' means unlimited rate"},}
		}
	case "-port-list":
		if len(args) == 2 {
			return []prompt.Suggest{{"80,80-88", "Default 1000+ port"},}
		}
	case "-timeout":
		if len(args) == 2 {
			return []prompt.Suggest{{"3000", "default 3000ms"},}
		}
	case "-thread":
		if len(args) == 2 {
			return []prompt.Suggest{{"128", "default 128"},}
		}
	case "-plugin":
		if len(args) == 2 {
			return []prompt.Suggest{{"ssh", "default all"},}
		}
	case "-user-list":
		if len(args) == 2 {
			return []prompt.Suggest{{"admin", "default all"},}
		}
	case "-pass-list":
		if len(args) == 2 {
			return []prompt.Suggest{{"admin", "default all"},}
		}
	case "-alive-host-check":
		if len(args) == 2 {
			return []prompt.Suggest{{"", "Bool"},}
		}
	case "-only-alive-network":
		if len(args) == 2 {
			return []prompt.Suggest{{"", "Bool"},}
		}
	}
	return []prompt.Suggest{}
}
