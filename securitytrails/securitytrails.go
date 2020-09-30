package securitytrails

import (
	"encoding/csv"
	"fmt"
	"github.com/zsdevX/DarkEye/common"
	"os"
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

	//先创建文件，创建失败结束
	f, err := os.Create(saveFile)
	if err != nil {
		_, _ = fmt.Fprint(os.Stderr, "初始化失败", err.Error())
		return
	}
	defer f.Close()
	_, _ = f.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM
	w := csv.NewWriter(f)
	_ = w.Write([]string{"域名", "CNAME", "中间件", "标题", "IP"})

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
			fmt.Sprintf("收集信息任务完成，有效数量%d, 已保存结果:%s", logNumber, saveFile), common.INFO)
	}
}
