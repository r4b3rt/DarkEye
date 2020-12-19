package fofa

import (
	"fmt"
	"github.com/zsdevX/DarkEye/common"
	"strings"
	"time"
)

func (f *Fofa) Run() {
	//初始化用来存储结果
	f.ipNodes = make([]ipNode, 0)
	f.runIP()

	w, file, fileName, err := common.CreateCSV("fofa",
		[]string{"端口有效", "IP", "PORT", "网址", "标题", "中间件", "指纹"})
	if err != nil {
		f.ErrChannel <- common.LogBuild("fofa",
			fmt.Sprintf("创建记录文件失败:"+err.Error()), common.INFO)
		return
	}
	defer file.Close()

	logNumber := 0
	for _, n := range f.ipNodes {
		alive := "失效"
		if n.Alive == common.Alive {
			alive = "有效"
		} else if n.Alive == common.TimeOut {
			alive = "超时"
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
			fmt.Sprintf("收集信息任务完成，有效数量%d, 已保存结果:%s", logNumber, fileName), common.INFO)
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
