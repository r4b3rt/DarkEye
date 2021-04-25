package main

import (
	"github.com/c-bata/go-prompt"
)

var (
	zoomEyeSuggestions = []prompt.Suggest{
		// Command
		{"-api", "API-KEY"},
		{"-page", "返回查询页面范围(每页20条):开始页-结束页"},
		{"+", "AND运算"},
		{"-", "排除运算"},
		{"ip:", "搜索指定IPv4地址相关资产"},
		{"cidr:", "搜索IP的C段资产"},
		{"ssl:", `搜索ssl证书存在"google"字符串的资产。常常用来提过产品名及公司名搜索对应目标`},
		{"port:", "搜索相关端口资产"},
		{"service:", "搜索对应服务协议的资产"},
		{"title:", `搜索html内容里标题中搜索数据`},
		{"banner:", `服务器应答信息`},
		{"site:", "搜索域名相关的资产"},
		{"app:", `搜索中间件`},
		{"hostname:", `搜索相关IP"主机名"的资产`},
		{"device:", `搜索路由器相关的设备类型`},
		{"os:", `搜索相关操作系统`},
		{"organization:", `常常用来定位大学、结构、大型互联网公司对应IP资产`},
		{"isp:", `搜索相关网络服务提供商的资产`},
		{"asn:", `搜索对应ASN（Autonomous system number）自治系统编号相关IP资产`},
		{"city:", `搜索相关城市资产`},
		{"subdivisions:", `搜索相关指定行政区的资产,中国省会支持拼音及汉字书写如subdivisions:"北京"subdivisions:"beijing"`},
		{"country:", `搜索国家地区资产,可以使用国家缩写，也可以使用中/英文全称如country:"中国"country:"china"`},
	}
	zoomEyeValue = map[string]bool{
		"-api":          false,
		"-page":         false,
		"+":             true,
		"-":             true,
		"ip:":           false,
		"cidr:":         false,
		"ssl:":          false,
		"port:":         false,
		"service:":      false,
		"title:":        false,
		"banner:":       false,
		"site:":         false,
		"app:":          false,
		"hostname:":     false,
		"device:":       false,
		"os:":           false,
		"organization:": false,
		"isp:":          false,
		"asn:":          false,
		"city:":         false,
		"subdivisions:": false,
		"country:":      false,
	}
)

func init() {
	mSuggestions = append(mSuggestions, []prompt.Suggest{
		{zoomEye, "Analysis the result of intelligence collection"},
	}...)
}

func (z *zoomEyeRuntime) Completer(args []string) []prompt.Suggest {
	if len(args) <= 1 {
		filtered := make([]string, 0)
		for _, v := range z.parent.CmdArgs {
			if v != "-api" && v != "-page" {
				continue
			}
			filtered = append(filtered, v)
		}
		return filterSuggestions(runCompleteCheck(zoomEyeSuggestions, z.parent.CmdArgs,
			[]string{
				"-api",
			}), filtered)
	}
	//过滤重复的命令
	if args[0] == "-api" || args[0] == "-page" {
		if isDuplicateArg(args[0], z.parent.CmdArgs) {
			return []prompt.Suggest{}
		}
	}
	switch args[0] {
	case "-api":
		if len(args) == 2 {
			return []prompt.Suggest{{"you-api-key", "获取方式 https://www.zoomeye.org/profile"},}
		}
	case "-page":
		if len(args) == 2 {
			return []prompt.Suggest{{"1-5", "返回查询页面范围(每页20条):开始页-结束页"},}
		}
	case "+":
		if len(args) == 2 {
			return []prompt.Suggest{{`+`, "AND运算"},}
		}
	case "-":
		if len(args) == 2 {
			return []prompt.Suggest{{`-`, "排除运算"},}
		}
	case "ip:":
		if len(args) == 2 {
			return []prompt.Suggest{{`"8.8.8.8"`, "搜索指定IPv4地址相关资产"},}
		}
	case "cidr:":
		if len(args) == 2 {
			return []prompt.Suggest{{`"8.8.8.8/24"`, "搜索IP的C段资产"},}
		}
	case "ssl:":
		if len(args) == 2 {
			return []prompt.Suggest{{`"google""`,
				`搜索ssl证书存在"google"字符串的资产。常常用来提过产品名及公司名搜索对应目标`},}
		}
	case "port:":
		if len(args) == 2 {
			return []prompt.Suggest{{`80`, "搜索相关端口资产"},}
		}
	case "service:":
		if len(args) == 2 {
			return []prompt.Suggest{{`ssh`, "搜索对应服务协议的资产"},}
		}
	case "site:":
		if len(args) == 2 {
			return []prompt.Suggest{{`google.com`, "搜索域名相关的资产"},}
		}
	case "title:":
		if len(args) == 2 {
			return []prompt.Suggest{{`"Cisco"`, `搜索html内容里标题中存在"Cisco"的数据`},}
		}
	case "banner:":
		if len(args) == 2 {
			return []prompt.Suggest{{`"OpenSSh"`, `服务器应答包含的信息`},}
		}
	case "app:":
		if len(args) == 2 {
			return []prompt.Suggest{{`"Cisco ASA SSL VPN"`, `搜索思科ASA-SSL-VPN的设备`},}
		}
	case "hostname:":
		if len(args) == 2 {
			return []prompt.Suggest{{`"google.com"`, `搜索相关IP"主机名"的资产`},}
		}
	case "device:":
		if len(args) == 2 {
			return []prompt.Suggest{{`"router"`, `搜索路由器相关的设备类型`},}
		}
	case "os:":
		if len(args) == 2 {
			return []prompt.Suggest{{`"Linux"`, `搜索相关操作系统`},}
		}
	case "organization:":
		if len(args) == 2 {
			return []prompt.Suggest{{`"北京大学"`, `常常用来定位大学、结构、大型互联网公司对应IP资产`},}
		}
	case "isp:":
		if len(args) == 2 {
			return []prompt.Suggest{{`"China Mobile"`, `搜索相关网络服务提供商的资产`},}
		}
	case "asn:":
		if len(args) == 2 {
			return []prompt.Suggest{{`42893`,
				`搜索对应ASN（Autonomous system number）自治系统编号相关IP资产`},}
		}
	case "city:":
		if len(args) == 2 {
			return []prompt.Suggest{{`"changsha"`,
				`搜索相关城市资产`},}
		}
	case "subdivisions:":
		if len(args) == 2 {
			return []prompt.Suggest{{`"beijing"`,
				`搜索相关指定行政区的资产,中国省会支持拼音及汉字书写如subdivisions:"北京"subdivisions:"beijing"`},}
		}
	case "country:":
		if len(args) == 2 {
			return []prompt.Suggest{{`"CN"`,
				`搜索国家地区资产,可以使用国家缩写，也可以使用中/英文全称如country:"中国"country:"china"`},}
		}
	}
	return []prompt.Suggest{}
}
