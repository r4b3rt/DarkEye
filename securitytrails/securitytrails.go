package securitytrails

import (
	"fmt"
	"github.com/zsdevX/DarkEye/common"
	"strings"
)

func (s *SecurityTrails) Run() {
	queries := strings.Split(s.Queries, ",")
	s.dns = make([]dnsInfo, 0)
	for _, q := range queries {
		s.get(q)
	}

	saveFile := common.GenFileName("dns")
	logNumber := 0
	for k, n := range s.dns {
		if k == 0 {
			_ = common.SaveFile("域名,CNAME,IP", saveFile)
		}
		ipi := ""
		for i, ip := range n.ip {
			if i != 0 {
				ipi += "|"
			}
			ipi += fmt.Sprintf("%s_%s_%s", ip.ip, ip.RegionName, ip.Isp)
		}
		line := fmt.Sprintf("%s,%s,%s",
			n.domain,
			n.cname,
			ipi)

		if err := common.SaveFile(line, saveFile); err != nil {
			s.ErrChannel <- common.LogBuild("SecurityTrails",
				fmt.Sprintf("收集信息任务失败，无法保存结果:%s", saveFile), common.FAULT)
			return
		}
		logNumber++
	}
	if logNumber == 0 {
		s.ErrChannel <- common.LogBuild("SecurityTrails",
			fmt.Sprintf("收集信息任务完成，未有结果"), common.INFO)
	} else {
		s.ErrChannel <- common.LogBuild("SecurityTrails",
			fmt.Sprintf("收集信息任务完成，有效数量%d, 已保存结果:%s", logNumber, saveFile), common.INFO)
	}
}
