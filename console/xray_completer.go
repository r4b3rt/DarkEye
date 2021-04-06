package main

import (
	"github.com/c-bata/go-prompt"
)

var (
	xRayProgram     = "xRay"
	xRaySuggestions = []prompt.Suggest{
		// Command
		{"-download", "https://binary_url"},
		{"-url", "url or url-file-list"},
		{"-proxy-port", "被动监听端口"},
		{"-chrome", "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"},
	}
	xRayValueCheck = map[string]bool{
		"-download":   false,
		"-url":        false,
		"-proxy-port": false,
		"-chrome":     false,
	}
)

func init() {
	mSuggestions = append(mSuggestions, []prompt.Suggest{
		{xRayProgram, "Vulnerable Scanner"},
	}...)
}

func (ctx *RequestContext) xRayArgumentsCompleter(args []string) []prompt.Suggest {
	if len(args) <= 1 {
		return filterSuggestions(runCompleteCheck(analysisSuggestions, ctx.CmdArgs,
			[]string{
				"-url",
			}), ctx.CmdArgs)
	}
	switch args[0] {
	case "-url":
		if len(args) == 2 {
			return []prompt.Suggest{{"https://vuln.com.cn", "vulnerable url"},
				{"url-file-list", "url.txt"},}
		}

	case "-proxy-port":
		if len(args) == 2 {
			return []prompt.Suggest{{"7777", "被动监听端口"},}
		}
	case "-chrome":
		if len(args) == 2 {
			return []prompt.Suggest{{"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome", "Chrome浏览器程序路径"},}
		}
	case "-download":
		if len(args) == 2 {
			return []prompt.Suggest{
				{
					"https://ghproxy.com/https://github.com/zsdevX/helper/releases/download/1/xray_darwin_amd64",
					"xRay binary address"},
				{
					"https://ghproxy.com/https://github.com/zsdevX/helper/releases/download/1/xray_linux_amd64",
					"xRay binary address"},
				{
					"https://ghproxy.com/https://github.com/zsdevX/helper/releases/download/1/xray_windows_386",
					"xRay binary address"},
			}
		}
	}
	return []prompt.Suggest{}
}
