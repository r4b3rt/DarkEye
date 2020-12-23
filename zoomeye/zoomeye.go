package zoomeye

import (
	"encoding/json"
	"fmt"
	"github.com/zsdevX/DarkEye/common"
	"strconv"
)

//Run add comment
func (z *ZoomEye) Run() {
	tot, _ := strconv.Atoi(z.Pages)
	idx := 1
	if tot == -1 {
		tot = 65535
	}
	targets := make([]TargetInfo, 0)
	for idx <= tot {
		z.runAPI(idx)
		if z.Results.Total <= 0 {
			z.ErrChannel <- common.LogBuild("zoomEye",
				fmt.Sprintf("获取异常%s:%s", z.Query, "无信息"), common.INFO)
			break
		}
		if z.Results.Available == 0 {
			z.ErrChannel <- common.LogBuild("zoomEye",
				fmt.Sprintf("获取异常%s:%s", z.Query, "超过'ip'语法频率限制，规则不清楚，哭死！"), common.FAULT)
			break
		}
		//如果不一致是什么情况？
		if z.Results.Total != z.Results.Available {
			z.ErrChannel <- common.LogBuild("zoomEye",
				fmt.Sprintf("获取异常%s:%s", z.Query, "z.Results.Total != z.Results.Available"), common.INFO)
			break
		}
		targets = append(targets, z.Results.Matches...)
		if idx == 1 && z.Pages == "-1" {
			tot = z.Results.Total / 20
			if z.Results.Total%20 != 0 {
				tot++
			}
		}
		z.ErrChannel <- common.LogBuild("zoomEye",
			fmt.Sprintf("%s:共%d页获取第%d页信息共计%d条", z.Query, tot, idx, len(z.Results.Matches)), common.INFO)
		idx++
		if common.ShouldStop(&z.Stop) {
			break
		}
	}
	match, err := json.Marshal(targets)
	if err != nil {
		z.ErrChannel <- common.LogBuild("zoomEye",
			fmt.Sprintf("获取信息%s:%s", z.Query, err.Error()), common.FAULT)
		return
	}

	filename, err := common.Write2CSV(z.Query+"_zoomEye_", match)
	if err != nil {
		z.ErrChannel <- common.LogBuild("zoomEye",
			fmt.Sprintf("获取信息%s:%s", z.Query, err.Error()), common.FAULT)
		return
	}
	z.ErrChannel <- common.LogBuild("zoomEye",
		fmt.Sprintf("收集信息任务完成，有效数量%d, 已保存结果:%s", len(targets), filename), common.INFO)
}
