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
	checkFuncs          = map[int]func(*Plugins) interface{}{}
)

func (plg *Plugins) Check() {
	if common.IsAlive(plg.TargetIp, plg.TargetPort, plg.TimeOut) != common.Alive {
		return
	}
	plg.PortOpened = true
	i := 0

	for i < PluginNR {
		if p := checkFuncs[i](plg); p != nil {
			color.Yellow("\n%s %s:%s %v\n", plg.TargetProtocol, plg.TargetIp, plg.TargetPort, p)
			break
		}
		i++
	}
	if i >= PluginNR {
		color.Yellow("%s %s:%s %v", "[-]", plg.TargetIp, plg.TargetPort, "Opened")
	}
	return
}

func loadDic(name string) []string {
	filename := filepath.Join("dic", fmt.Sprintf("dic_%s", name))
	file, err := os.Open(filename)
	if err != nil {
		panic("No dictionary:" + filename)
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
		if one == "空" {
			result = append(result, "")
		} else {
			result = append(result, one)
		}
	}
	return result
}
