package plugins

import (
	"bufio"
	"fmt"
	"github.com/fatih/color"
	"github.com/zsdevX/DarkEye/common"
	"os"
	"path/filepath"
	"strings"
)

var (
	checkFuncs = map[int]func(*Plugins) interface{}{}
)

func (plg *Plugins) Check() {
	plg.RateWait(plg.RateLimiter) //活跃端口发包限制
	if common.IsAlive(plg.TargetIp, plg.TargetPort, plg.TimeOut) != common.Alive {
		return
	}
	plg.PortOpened = true
	i := 0
	//预处理
	for i < PluginNR {
		if p := checkFuncs[i](plg); p != nil {
			if plg.highLight {
				color.Green("\n%s %s:%s %v\n", plg.TargetProtocol, plg.TargetIp, plg.TargetPort, p)
			} else {
				color.Yellow("\n%s %s:%s %v\n", plg.TargetProtocol, plg.TargetIp, plg.TargetPort, p)
			}
			break
		}
		i++
	}
	if i >= PluginNR {
		color.Yellow("\n%s %s:%s %v\n", "[-]", plg.TargetIp, plg.TargetPort, "Opened")
	}
	return
}

func (plg *Plugins) PreCheck() {
	//预处理
	//137端口机器检查
	nbCheck(plg)
	if plg.PortOpened {
		color.Yellow("\n%s %s:%s %v\n", plg.TargetProtocol, plg.TargetIp, 137, plg.NetBios)
	}
}

func loadDic(name string) []string {
	filename := filepath.Join("dic", fmt.Sprintf("dic_%s", name))
	file, err := os.Open(filename)
	if err != nil {
		color.Red("未发现字典文件:" + filename)
		return nil
	}
	defer file.Close()
	result := make([]string, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		one := scanner.Text()
		if strings.HasPrefix(one, "#") {
			continue
		}
		one = strings.TrimSpace(one)
		one = strings.Trim(one, "\r\n")
		if one == "空" { //超级口令的""特殊表示
			result = append(result, "")
		} else {
			result = append(result, one)
		}
	}
	return result
}
