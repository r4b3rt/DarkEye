package zoomeye

import (
	"fmt"
	"github.com/tidwall/gjson"
	"github.com/zsdevX/DarkEye/common"
	url2 "net/url"
	"time"
)

func (z *ZoomEye) run(page int) *gjson.Result {
	url := fmt.Sprintf("https://api.zoomeye.org/host/search?query=%s&page=%d", url2.QueryEscape(z.Query), page)
	req := common.HttpRequest{
		Url:     url,
		TimeOut: time.Duration(10),
		Method:  "GET",
		Headers: map[string]string{
			"User-Agent": common.UserAgents[1],
			"API-KEY":    z.ApiKey,
		},
	}
	response, err := req.Go()
	if err != nil && response == nil {
		z.ErrChannel <- common.LogBuild("zoomEye",
			fmt.Sprintf("获取信息失败%s:%s", z.Query, err.Error()), common.FAULT)
		return nil
	}
	switch response.Status {
	case 201:
		z.ErrChannel <- common.LogBuild("zoomEye", "不支持相关展示", common.FAULT)
	case 200:
		ret := gjson.GetBytes(response.Body, "matches")
		return &ret
	default:
		//4xx错误
		z.ErrChannel <- common.LogBuild("zoomEye",
			fmt.Sprintf("获取信息失败%s:%s", z.Query, string(response.Body)), common.FAULT)
	}
	return nil
}
