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
	w, file, fileName, err := common.CreateCSV("dns",
		[]string{"域名", "CNAME", "中间件", "标题", "IP"})
	if err != nil {
		s.ErrChannel <- common.LogBuild("fofa",
			fmt.Sprintf("创建记录文件失败:"+err.Error()), common.INFO)
		return
	}
	defer file.Close()

	logNumber := 0
	for _, n := range s.dns {
		ipi := ""
		for i, ip := range n.ip {
			if i != 0 {
				ipi += "|"
			}
			ipi += fmt.Sprintf("%s_%s_%s", ip.ip, ip.RegionName, ip.Isp)
		}
		_ = w.Write([]string{n.domain, n.cname, n.server, n.title, ipi})
		logNumber++
	}
	w.Flush()
	if logNumber == 0 {
		s.ErrChannel <- common.LogBuild("SecurityTrails",
			fmt.Sprintf("收集信息任务完成，未有结果"), common.INFO)
	} else {
		s.ErrChannel <- common.LogBuild("SecurityTrails",
			fmt.Sprintf("收集信息任务完成，有效数量%d, 已保存结果:%s", logNumber, fileName), common.INFO)
	}
}
