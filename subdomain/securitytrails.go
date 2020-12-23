package subdomain

import (
	"encoding/json"
	"fmt"
	"github.com/zsdevX/DarkEye/common"
	"math/rand"
	"time"
)

/*
curl "https://api.securitytrails.com/v1/domain/apple.com" \
 -H 'apikey: xxxxxxxxx'
*/

type subMeta struct {
	LimitReached bool `json:"limit_reached"`
}

type subResult struct {
	Subdomains []string `json:"subdomains"`
	Meta       subMeta  `json:"meta"`
}

func (s *SubDomain) getSecurityTrails(query string) {
	url := fmt.Sprintf("https://api.securitytrails.com/v1/domain/%s/subdomains?children_only=false", query)
	req := common.HttpRequest{
		Url:     url,
		TimeOut: time.Duration(10),
		Method:  "GET",
		Headers: map[string]string{
			"User-Agent":   common.UserAgents[rand.Int()%len(common.UserAgents)],
			"apikey":       s.ApiKey,
			"Content-Type": "application/json",
		},
	}
	response := s.fetchSubDomainResults(&req, query)
	if response == nil {
		return
	}
	res := subResult{}
	if err := json.Unmarshal(response.Body, &res); err != nil {
		s.ErrChannel <- common.LogBuild("subDomain.get",
			fmt.Sprintf("收集子域%s处理返回数据失败:%s", query, err.Error()), common.FAULT)
		return
	}
	if res.Meta.LimitReached {
		s.ErrChannel <- common.LogBuild("subDomain.get",
			fmt.Sprintf("子域%s返回数量%d超过达到SecurityTrails服务器允许上限", query, len(res.Subdomains)), common.FAULT)
	}

	for _, r := range res.Subdomains {
		s.try(r, query)
		if common.ShouldStop(&s.Stop) {
			break
		}
	}
}

func (s *SubDomain) fetchSubDomainResults(req *common.HttpRequest, query string) *common.HttpResponse {
	retry := 0
	for {
		if common.ShouldStop(&s.Stop) {
			break
		}
		response, err := req.Go()
		if err != nil {
			retry++
			s.ErrChannel <- common.LogBuild("subDomain.get",
				fmt.Sprintf("收集子域%s请求失败。网络错误，本次请求尝试（%d）次,错误:%s",
					query, retry, err.Error()), common.FAULT)
			time.Sleep(time.Second * 3)
			continue
		}
		return response
	}
	return nil
}
