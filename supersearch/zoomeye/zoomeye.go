package zoomeye

import (
	"context"
	"fmt"
	"github.com/b1gcat/DarkEye/common"
	"github.com/tidwall/gjson"
	"net/url"
	"strconv"
	"strings"
	"time"
)

/*Run add comment
@n:Number of results
@authCode: key or password
@query: type:query
@facet: stat
*/
func Run(ctx context.Context, n int, authCode, query, facet string, log chan string) []Match {
	//初始化key
	accessToken, err := login(ctx, authCode)
	if err != nil {
		common.LogUi(err.Error(), log, common.FAULT)
		return nil
	}
	apiKey := authCode
	if accessToken != "" {
		apiKey = ""
	}
	//获取用户信息
	if err := resource(ctx, accessToken, apiKey, log); err != nil {
		common.LogUi(err.Error(), log, common.FAULT)
		return nil
	}
	ret := make([]Match, 0)
	i := 1
	if n <= 0 {
		n = 65535
	}
	n = n/20 + 1
	for i <= n {
		data := gogo(ctx, i, accessToken, apiKey, query, facet, log)
		if data == nil {
			return ret
		}
		if i == 1 {
			firstShow(data, facet, log)
		}
		ret = getData(ret, data)
		log <- fmt.Sprintf("%s:获取第%d页信息共%d个", query, i, len(ret))
		time.Sleep(time.Second * 3)
		i++
	}
	return ret
}

func firstShow(data []byte, facet string, log chan string) {
	matches := gjson.GetBytes(data, "facets")
	f := strings.Split(facet, ",")

	total := gjson.GetBytes(data, "total").String()
	avail := gjson.GetBytes(data, "available").String()
	common.LogUi("Total/available:"+total+"/"+avail, log, common.INFO)
	for _, fv := range f {
		common.LogUi(fv+":", log, common.INFO)
		for _, v := range matches.Get(fv).Array() {
			name := v.Get("name").String()
			count := v.Get("count").String()
			if name == "" {
				name = "unknown"
			}
			common.LogUi("	"+name+":"+count, log, common.INFO)
		}
	}
}

func getData(ret []Match, data []byte) []Match {
	matches := gjson.GetBytes(data, "matches")
	for _, match := range matches.Array() {
		m := Match{}
		m.Ip = match.Get("ip").String()
		m.Country = match.Get("geoinfo.country.names.en").String()

		m.Port = int(match.Get("portinfo.port").Int())
		m.Os = match.Get("portinfo.os").String()
		m.Hostname = match.Get("portinfo.hostname").String()
		m.Service = match.Get("portinfo.service").String()
		m.Banner = match.Get("portinfo.banner").String()
		if len(m.Banner) > 128 {
			//Banner太大有点乱
			m.Banner = m.Banner[:128]
		}
		m.Title = match.Get("portinfo.title").String()
		m.Version = match.Get("portinfo.version").String()
		m.Device = match.Get("portinfo.device").String()
		m.ExtraInfo = match.Get("portinfo.extrainfo").String()
		m.RDns = match.Get("portinfo.rdns").String()
		m.App = match.Get("portinfo.app").String()

		if m.Service == "http" || m.Service == "https" {
			m.Url = m.Service + "://" + m.Ip + ":" + strconv.Itoa(m.Port)
			if m.Banner != "" {
				x := strings.Split(m.Banner, " ")
				if len(x) == 3 {
					m.HttpCode, _ = strconv.Atoi(x[2])
				}
			}
		}
		ret = append(ret, m)
	}
	return ret
}

//https://www.zoomeye.org/doc#resources-info
func resource(ctx context.Context, accessToken, apiKey string, log chan string) error {
	req := common.HttpRequest{
		Url:     "https://api.zoomeye.org/resources-info",
		TimeOut: time.Duration(10),
		Method:  "GET",
		Ctx:     ctx,
		Headers: map[string]string{
			"User-Agent":    common.UserAgents[1],
			"Authorization": "JWT " + accessToken,
			"API-KEY":       apiKey,
		},
	}
	resp, err := req.Go()
	if err != nil {
		return err
	}
	common.LogUi(string(resp.Body), log, common.INFO)
	return nil
}

//https://www.zoomeye.org/doc#login
func login(ctx context.Context, authCode string) (string, error) {
	if authCode == "" {
		return "", fmt.Errorf("未设置认证key")
	}
	if !strings.Contains(authCode, "/") {
		return "", nil
	}
	auth := strings.Split(authCode, "/")
	req := common.HttpRequest{
		Url:     "https://api.zoomeye.org/user/login",
		TimeOut: time.Duration(10),
		Method:  "POST",
		Ctx:     ctx,
		Headers: map[string]string{
			"Content_type": "application/json; charset=UTF-8",
			"User-Agent":   common.UserAgents[1],
		},
		Body: []byte(fmt.Sprintf(`{"username":"%s", "password":"%s"}`, auth[0], auth[1])),
	}
	resp, err := req.Go()
	if err != nil {
		return "", err
	}

	accessToken := gjson.GetBytes(resp.Body, "access_token").String()
	if accessToken == "" {
		return "", fmt.Errorf("认证错误 %v", string(resp.Body))
	}
	return accessToken, nil
}

func gogo(ctx context.Context, page int, accessToken, apiKey, query, facet string, log chan string) []byte {
	q := strings.SplitN(query, "@", 2)
	if strings.ToLower(q[0]) == "web" {
		common.LogUi("暂时不支持web类型查询", log, common.FAULT)
		return nil
	}
	url := fmt.Sprintf("https://api.zoomeye.org/%s/search?query=%s&page=%d&facets=%s",
		strings.ToLower(q[0]), url.QueryEscape(q[1]), page, url.QueryEscape(facet))
	req := common.HttpRequest{
		Url:     url,
		TimeOut: time.Duration(30),
		Method:  "GET",
		Ctx:     ctx,
		Headers: map[string]string{
			"User-Agent":    common.UserAgents[1],
			"Authorization": "JWT " + accessToken,
			"API-KEY":       apiKey,
		},
	}
	response, err := req.Go()
	if err != nil && response == nil {
		common.LogUi(fmt.Sprintf("获取信息失败%s:%s", query, err.Error()), log, common.FAULT)
		return nil
	}
	switch response.Status {
	case 201:
		common.LogUi("不支持相关展示", log, common.FAULT)
	case 200:
		return response.Body
	default:
		//4xx错误
		common.LogUi(fmt.Sprintf("获取信息失败%s:%s", query, string(response.Body)), log, common.FAULT)
	}
	return nil
}

type Match struct {
	Ip      string
	Country string
	//PortInfo
	Port      int
	Os        string
	Hostname  string
	Service   string
	Banner    string
	Title     string
	Version   string
	Device    string
	ExtraInfo string
	RDns      string
	App       string
	//
	Url      string
	HttpCode int
}
