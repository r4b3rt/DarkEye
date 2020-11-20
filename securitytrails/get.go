package securitytrails

import (
	"encoding/json"
	"fmt"
	"github.com/miekg/dns"
	"github.com/zsdevX/DarkEye/common"
	"golang.org/x/time/rate"
	"math/rand"
	"strings"
	"time"
)

/*
curl "https://api.securitytrails.com/v1/domain/apple.com" \
 -H 'apikey: xxxxxxxxx'
*/

type subMeta struct {
	Limit_reached bool `json:"limit_reached"`
}

type subResult struct {
	Subdomains []string `json:"subdomains"`
	Meta       subMeta  `json:"meta"`
}

var (
	ipApiLimit = rate.NewLimiter(rate.Every(1500*time.Millisecond), 40) //burst 40，以后1.5秒分配资源
)

func (s *SecurityTrails) get(query string) {
	//查询ip历史信息太浪费api，大佬有需要在查吧。
	s.ErrChannel <- common.LogBuild("SecurityTrails",
		fmt.Sprintf("若需要查到域名历史信息请直接用如下命令:\n %s",
			`curl --request GET --url "https://api.securitytrails.com/v1/history/${host}/dns/a" -H "apikey: ${you-api-key}" --header 'accept: application/json' `),
		common.INFO)

	s.ErrChannel <- common.LogBuild("SecurityTrails",
		fmt.Sprintf("开始收集子域%s", query), common.INFO)

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
		s.ErrChannel <- common.LogBuild("SecurityTrails.get",
			fmt.Sprintf("收集子域%s处理返回数据失败:%s", query, err.Error()), common.FAULT)
		return
	}
	if res.Meta.Limit_reached {
		s.ErrChannel <- common.LogBuild("SecurityTrails.get",
			fmt.Sprintf("子域%s返回数量%d超过达到SecurityTrails服务器允许上限", query, len(res.Subdomains)), common.FAULT)
	}

	for _, r := range res.Subdomains {
		d := dnsInfo{
			domain: r + "." + query,
			ip:     make([]ipInfo, 0),
		}
		if s.IpCheck {
			s.parseIP(&d)
		}
		s.parseTag(&d)

		s.dns = append(s.dns, d)
		if common.ShouldStop(&s.Stop) {
			break
		}
	}
}

func (s *SecurityTrails) fetchSubDomainResults(req *common.HttpRequest, query string) (*common.HttpResponse) {
	retry := 0
	for {
		if common.ShouldStop(&s.Stop) {
			break
		}
		response, err := req.Go()
		if err != nil {
			retry++
			s.ErrChannel <- common.LogBuild("SecurityTrails.get",
				fmt.Sprintf("收集子域%s请求失败。网络错误，本次请求尝试（%d）次,错误:%s",
					query, retry, err.Error()), common.FAULT)
			time.Sleep(time.Second * 3)
			continue
		}
		return response
	}
	return nil
}

func (s *SecurityTrails) parseTag(d *dnsInfo) {
	d.server, d.title = common.GetHttpTitle("http", d.domain)
	if d.server == "" && d.title == "" {
		d.server, d.title = common.GetHttpTitle("https", d.domain)
	}
	return
}

func (s *SecurityTrails) parseIP(d *dnsInfo) {
	//做dns解析
	c := dns.Client{
		Timeout: 5 * time.Second,
	}
	m1 := &dns.Msg{}
	m1.SetQuestion(d.domain+".", dns.TypeA)
	server := s.DnsServer
	if !strings.Contains(server, ":") {
		server = server + ":53"
	}
	r, _, err := c.Exchange(m1, server)
	if err != nil {
		s.ErrChannel <- common.LogBuild("SecurityTrails.get.parseIP",
			fmt.Sprintf("解析域名失败%s:%s", d.domain, err.Error()), common.FAULT)
		return
	}
	for _, a := range r.Answer {
		if cn, ok := a.(*dns.CNAME); ok {
			if d.cname != "" {
				d.cname += "|"
			}
			d.cname += cn.Target
		}

		if ip, ok := a.(*dns.A); ok {
			ipi := ipInfo{
				ip: ip.A.String(),
			}
			s.parseIpLoc(&ipi)
			d.ip = append(d.ip, ipi)
		}
	}
	s.ErrChannel <- common.LogBuild("SecurityTrails",
		fmt.Sprintf("解析域名%s完成CNAME=%s,ip数量%d", d.domain, d.cname, len(d.ip)), common.ALERT)
}

func (s *SecurityTrails) parseIpLoc(ipi *ipInfo) {
	for {
		if ipApiLimit.Allow() {
			break
		} else {
			time.Sleep(time.Millisecond * 1500)
		}
	}

	url := fmt.Sprintf("http://ip-api.com/json/%s?lang=zh-CN", ipi.ip)
	req := common.HttpRequest{
		Url:     url,
		TimeOut: time.Duration(5),
		Method:  "GET",
		Headers: map[string]string{
			"User-Agent":   common.UserAgents[0],
			"Content-Type": "application/json",
		},
	}
	response, err := req.Go()
	if err != nil {
		s.ErrChannel <- common.LogBuild("SecurityTrails.get.parseIP.parseIpLoc",
			fmt.Sprintf("获取IP%s所在地失败%s", ipi.ip, err.Error()), common.FAULT)
		return
	}
	_ = json.Unmarshal(response.Body, ipi)
	time.Sleep(time.Second * 1)
}
