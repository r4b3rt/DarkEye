package zoomeye

import (
	"encoding/json"
	"fmt"
	"github.com/zsdevX/DarkEye/common"
)

func (z *ZoomEye) Run() {
	z.runAPI()

	if z.Results.Total <= 0 {
		z.ErrChannel <- common.LogBuild("zoomEye",
			fmt.Sprintf("获取信息%s:%s", z.Query, "无信息"), common.INFO)
		return
	}

	if z.Results.Available == 0 {
		z.ErrChannel <- common.LogBuild("zoomEye",
			fmt.Sprintf("获取信息%s:%s", z.Query, "超过'ip'语法频率限制，规则不清楚，哭死！"), common.FAULT)
		return
	}
	match, err := json.Marshal(z.Results.Matches)
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
		fmt.Sprintf("收集信息任务完成，有效数量%d, 已保存结果:%s", len(z.Results.Matches), filename), common.INFO)

}
