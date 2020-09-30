package fofa

import (
	"encoding/csv"
	"fmt"
	"github.com/zsdevX/DarkEye/common"
	"os"
	"strings"
)

func (f *Fofa) Run() {
	//初始化用来存储结果
	f.ipNodes = make([]ipNode, 0)

	if !f.FofaComma {
		f.runIP()
	} else {
		f.runRaw()
	}

	saveFile := common.GenFileName("fofa")
	file, err := os.Create(saveFile)
	if err != nil {
		_, _ = fmt.Fprint(os.Stderr, "初始化失败", err.Error())
		return
	}
	defer file.Close()
	_, _ = file.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM
	w := csv.NewWriter(file)
	_ = w.Write([]string{"端口有效", "IP", "PORT", "网址", "标题", "中间件", "指纹"})

	logNumber := 0
	for _, n := range f.ipNodes {
		alive := "失效"
		if n.Alive {
			alive = "有效"
		}
		_ = w.Write([]string{alive, n.Ip, n.Port, n.Domain, n.Title, n.Server, n.Finger})
		logNumber++
	}
	w.Flush()
	if logNumber == 0 {
		f.ErrChannel <- common.LogBuild("fofa",
			fmt.Sprintf("收集信息任务完成，未有结果"), common.INFO)
	} else {
		f.ErrChannel <- common.LogBuild("fofa",
			fmt.Sprintf("收集信息任务完成，有效数量%d, 已保存结果:%s", logNumber, saveFile), common.INFO)
	}
	return
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
		base, start, end, err := common.ParseNmapIP(ip)
		if err != nil {
			f.ErrChannel <- err.Error()
			return
		}
		for {
			if start > end {
				break
			}
			if common.ShouldStop(&f.Stop) {
				break
			}
			f.get(fmt.Sprintf("%s.%d", base, start))
			start++
		}
	}
}
