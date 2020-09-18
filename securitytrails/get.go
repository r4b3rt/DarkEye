package securitytrails

import (
	"encoding/json"
	"fmt"
	"github.com/miekg/dns"
	"github.com/zsdevX/DarkEye/common"
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

func (s *SecurityTrails) get(query string) {
	//查询ip历史信息太浪费api，大佬有需要在查吧。
	s.ErrChannel <- common.LogBuild("SecurityTrails",
		fmt.Sprintf("若需要查到域名历史信息请直接用如下命令:\n %s",
			`curl --request GET --url "https://api.securitytrails.com/v1/history/${host}/dns/a" -H "apikey: ${you-api-key}" --header 'accept: application/json' `),
		common.INFO)

	s.ErrChannel <- common.LogBuild("SecurityTrails",
		fmt.Sprintf("开始收集子域%s", query), common.INFO)

	url := fmt.Sprintf("https://api.securitytrails.com/v1/domain/%s/subdomains?children_only=false", query)
	userAgent := common.UserAgents[rand.Int()%len(common.UserAgents)]
	req := common.Http{
		Url:         url,
		TimeOut:     time.Duration(5),
		Method:      "GET",
		Referer:     url,
		H:           "apikey=" + s.ApiKey,
		Agent:       userAgent,
		ContentType: "application/json",
	}
	body, err := req.Http()
	if err != nil {
		s.ErrChannel <- common.LogBuild("SecurityTrails.get",
			fmt.Sprintf("收集子域%s发起请求失败。\n 如果不是网络问题请检查api是否使用到期（返回429错误）,如果到期大佬多申请几个账号吧。\n错误码 :%s",
				query, err.Error()), common.ALERT)
		return
	}
	res := subResult{}
	if err = json.Unmarshal(body, &res); err != nil {
		s.ErrChannel <- common.LogBuild("SecurityTrails.get",
			fmt.Sprintf("收集子域%s处理返回数据失败:%s", query, err.Error()), common.ALERT)
		return
	}
	if res.Meta.Limit_reached {
		s.ErrChannel <- common.LogBuild("SecurityTrails.get",
			fmt.Sprintf("收集子域%s达到服务器允许上限", query), common.ALERT)
	}

	for _, r := range res.Subdomains {
		d := dnsInfo{
			domain: r + "." + query,
			ip:     make([]ipInfo, 0),
		}
		if s.IpCheck {
			s.parseIP(&d)
		}
		s.dns = append(s.dns, d)
		if common.ShouldStop(&s.Stop) {
			break
		}
	}
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
			fmt.Sprintf("解析域名失败%s:%s", d.domain, err.Error()), common.ALERT)
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
		fmt.Sprintf("解析域名%s完成CNAME=%s,ip数量%d", d.domain, d.cname, len(d.ip)), common.INFO)
}

func (s *SecurityTrails) parseIpLoc(ipi *ipInfo) {
	url := fmt.Sprintf("http://ip-api.com/json/%s?lang=zh-CN", ipi.ip)
	userAgent := common.UserAgents[0]
	req := common.Http{
		Url:         url,
		TimeOut:     time.Duration(5),
		Method:      "GET",
		Referer:     url,
		Agent:       userAgent,
		ContentType: "application/json",
	}
	body, err := req.Http()
	if err != nil {
		s.ErrChannel <- common.LogBuild("SecurityTrails.get.parseIP.parseIpLoc",
			fmt.Sprintf("获取IP%s所在地失败%s", ipi.ip, err.Error()), common.ALERT)
		return
	}
	_ = json.Unmarshal(body, ipi)
	time.Sleep(time.Second * 1)
}
