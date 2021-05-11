package zoomeye

import (
	"context"
	"fmt"
	"github.com/tidwall/gjson"
	"github.com/zsdevX/DarkEye/common"
	url2 "net/url"
	"time"
)

func (z *ZoomEye) run(ctx context.Context, page int) *gjson.Result {
	u, _ := url2.Parse("https://api.zoomeye.org/host/search")
	query := u.Query()
	query.Add("page", fmt.Sprintf("%v", page))
	query.Add("query", z.Query)
	u.RawQuery = query.Encode()

	req := common.HttpRequest{
		Url:     u.String(),
		TimeOut: time.Duration(10),
		Method:  "GET",
		Ctx:     ctx,
		Headers: map[string]string{
			"User-Agent": common.UserAgents[1],
			"API-KEY":    z.ApiKey,
		},
	}
	response, err := req.Go()
	if err != nil && response == nil {
		z.ErrChannel <-
			fmt.Sprintf("获取信息失败%s:%s", z.Query, err.Error())
		return nil
	}
	switch response.Status {
	case 201:
		z.ErrChannel <- "不支持相关展示"
	case 200:
		ret := gjson.GetBytes(response.Body, "matches")
		return &ret
	default:
		//4xx错误
		z.ErrChannel <-
			fmt.Sprintf("获取信息失败%s:%s", req.Url, string(response.Body))
	}
	return nil
}
