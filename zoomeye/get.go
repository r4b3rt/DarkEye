package zoomeye

import (
	"encoding/json"
	"fmt"
	"github.com/zsdevX/DarkEye/common"
	url2 "net/url"
	"time"
)

func (z *ZoomEye) runAPI() {
	url := fmt.Sprintf("https://api.zoomeye.org/host/search?query=%s", url2.QueryEscape(z.Query))
	req := common.HttpRequest{
		Url:     url,
		TimeOut: time.Duration(5),
		Method:  "GET",
		Headers: map[string]string{
			"User-Agent": common.UserAgents[0],
			"API-KEY":    z.ApiKey,
		},
	}
	response, err := req.Go()
	if err != nil {
		z.ErrChannel <- common.LogBuild("zoomEye",
			fmt.Sprintf("获取信息失败%s:%s", z.Query, err.Error()), common.FAULT)
		return
	}

	switch response.Status {
	case 201:
		z.ErrChannel <- common.LogBuild("zoomEye", "不支持相关展示", common.FAULT)
	case 200:
		_ = json.Unmarshal(response.Body, &z.Results)
	default:
		//4xx错误
		z.ErrChannel <- common.LogBuild("zoomEye",
			fmt.Sprintf("获取信息失败%s:%s", z.Query, string(response.Body)), common.FAULT)
		return
	}
}
