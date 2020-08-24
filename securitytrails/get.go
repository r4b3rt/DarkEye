package securitytrails

import (
	"encoding/json"
	"fmt"
	"github.com/miekg/dns"
	"github.com/zsdevX/DarkEye/common"
	"math/rand"
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
	s.ErrChannel <- common.LogBuild("securitytrails",
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
		s.ErrChannel <- common.LogBuild("securitytrails.get",
			fmt.Sprintf("收集子域%s发起请求失败（如果不是网络问题请检查api是否使用到期，如果到期大佬多申请几个账号吧）:%s", query, err.Error()), common.ALERT)
		return
	}
	res := subResult{}
	if err = json.Unmarshal(body, &res); err != nil {
		s.ErrChannel <- common.LogBuild("securitytrails.get",
			fmt.Sprintf("收集子域%s处理返回数据失败:%s", query, err.Error()), common.ALERT)
		return
	}
	if res.Meta.Limit_reached {
		s.ErrChannel <- common.LogBuild("securitytrails.get",
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
		s.parseHistory(&d)
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
	r, _, err := c.Exchange(m1, s.DnsServer)
	if err != nil {
		s.ErrChannel <- common.LogBuild("securitytrails.get.parseIP",
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
	s.ErrChannel <- common.LogBuild("securitytrails",
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
		s.ErrChannel <- common.LogBuild("securitytrails.get.parseIP.parseIpLoc",
			fmt.Sprintf("获取IP%s所在地失败%s", ipi.ip, err.Error()), common.ALERT)
		return
	}
	_ = json.Unmarshal(body, ipi)
	time.Sleep(time.Second * 1)
}

type DnsHistoryRecords struct {
	Records []DnsHistory `json:"records"`
}

func (s *SecurityTrails) parseHistory(d *dnsInfo) {
	//https://api.securitytrails.com/v1/history/${host}/dns/a
	url := fmt.Sprintf("https://api.securitytrails.com/v1/history/%s/dns/a", d.domain)
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
		s.ErrChannel <- common.LogBuild("securitytrails.get.parseHistory",
			fmt.Sprintf("请求收集子域%s历史IP失败:%s", d.domain, err.Error()), common.ALERT)
		return
	}

	rr := DnsHistoryRecords{}
	if err = json.Unmarshal(body, &rr); err != nil {
		s.ErrChannel <- common.LogBuild("securitytrails.get.parseHistory",
			fmt.Sprintf("解析子域%s历史IP失败:%s", d.domain, err.Error()), common.ALERT)
		return
	}

	history := ""
	for k, v := range rr.Records {
		if k != 0 {
			history += "|"
		}
		history += "ip="
		for j, ip := range v.Values {
			if j != 0 {
				history += "+"
			}
			history += ip.Ip
		}
		history += "&organizations="
		for i, org := range v.Organizations {
			if i != 0 {
				history += "+"
			}
			history += org
		}
		history += "&first_seen=" + v.First_seen + "&last_seen=" + v.Last_seen
	}
	d.history = history
}
