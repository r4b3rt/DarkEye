package fofa

import (
	"encoding/json"
	"fmt"
	"github.com/zsdevX/DarkEye/common"
	"strings"
	"time"
)

func (f *Fofa) Run() {
	//初始化用来存储结果
	f.IpNodes = make([]IpNode, 0)
	f.runIP()

	if len(f.IpNodes) == 0 {
		f.ErrChannel <- common.LogBuild("foFa",
			fmt.Sprintf("收集信息任务完成，未有结果"), common.INFO)
		return
	}
	nodes, _ := json.Marshal(f.IpNodes)
	filename, err := common.Write2CSV(f.Ip+"foFa", nodes)
	if err != nil {
		f.ErrChannel <- common.LogBuild("foFa",
			fmt.Sprintf("获取信息%s:%s", f.Ip, err.Error()), common.FAULT)
		return
	}
	f.ErrChannel <- common.LogBuild("foFa",
		fmt.Sprintf("收集信息任务完成，有效数量%d, 已保存结果:%s", len(f.IpNodes), filename), common.INFO)
}

func (f *Fofa) runRaw() {
	f.ErrChannel <- common.LogBuild("fofa",
		fmt.Sprintf("开始收集信息:查询语法%s,间隔%d秒", f.Ip, f.Interval),
		common.INFO)
}

func (f *Fofa) runIP() {
	f.ErrChannel <- common.LogBuild("fofa",
		fmt.Sprintf("开始收集信息:IP范围%s,间隔%d秒", f.Ip, f.Interval),
		common.INFO)
	ips := strings.Split(f.Ip, ",")
	for _, ip := range ips {
		if common.ShouldStop(&f.Stop) {
			break
		}
		base, start, end, err := common.GetIPRange(ip)
		if err != nil {
			f.ErrChannel <- err.Error()
			return
		}
		for {
			nip := common.GenIP(base, start)
			if common.ShouldStop(&f.Stop) {
				break
			}
			if common.CompareIP(nip, end) > 0 {
				break
			}
			f.get(nip)
			start++
			//学做人，防止fofa封
			time.Sleep(time.Second * time.Duration(common.GenHumanSecond(f.Interval)))
		}
	}
}
