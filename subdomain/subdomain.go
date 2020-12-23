package subdomain

import (
	"encoding/json"
	"fmt"
	"github.com/miekg/dns"
	"github.com/zsdevX/DarkEye/common"
	"golang.org/x/time/rate"
	"strings"
	"time"
)

var (
	ipApiLimit = rate.NewLimiter(rate.Every(1500*time.Millisecond), 40) //burst 40，以后1.5秒分配资源
)

//Run add comment
func (s *SubDomain) Run() {
	queries := strings.Split(s.Queries, ",")
	s.Dns = make([]DnsInfo, 0)
	for _, q := range queries {
		s.get(q)
	}
	if len(s.Dns) == 0 {
		s.ErrChannel <- common.LogBuild("subDomain",
			fmt.Sprintf("收集信息任务完成，未有结果"), common.INFO)
		return
	}
	dns, _ := json.Marshal(s.Dns)
	filename, err := common.Write2CSV(s.Queries+"subDomain", dns)
	if err != nil {
		s.ErrChannel <- common.LogBuild("foFa",
			fmt.Sprintf("获取信息%s:%s", s.Queries, err.Error()), common.FAULT)
		return
	}

	s.ErrChannel <- common.LogBuild("subDomain",
		fmt.Sprintf("收集信息任务完成，有效数量%d, 已保存结果:%s", len(s.Dns), filename), common.INFO)
}

func (s *SubDomain) get(query string) {
	s.ErrChannel <- common.LogBuild("subDomain",
		fmt.Sprintf("开始收集子域%s", query), common.INFO)

	if !s.Brute {
		s.getSecurityTrails(query)
		return
	}
	s.getBrute(query)
}

func (s *SubDomain) try(prefix, query string) {
	d := DnsInfo{
		Domain: prefix + "." + query,
		Ip:     make([]IpInfo, 0),
	}
	if s.IpCheck {
		s.parseIP(&d)
	}
	if len(d.Ip) > 0 {
		s.parseTag(&d)
		s.Dns = append(s.Dns, d)
	}
}

func (s *SubDomain) parseTag(d *DnsInfo) {
	d.Server, d.Title, d.Code = common.GetHttpTitle("http", d.Domain, 5)
	if d.Server == "" && d.Title == "" {
		d.Server, d.Title, _ = common.GetHttpTitle("https", d.Domain, 5)
	}
	return
}

func (s *SubDomain) parseIP(d *DnsInfo) {
	//做dns解析
	c := dns.Client{
		Timeout: 5 * time.Second,
	}
	m1 := &dns.Msg{}
	m1.SetQuestion(d.Domain+".", dns.TypeA)
	server := s.DnsServer
	if !strings.Contains(server, ":") {
		server = server + ":53"
	}
	r, _, err := c.Exchange(m1, server)
	if err != nil {
		s.ErrChannel <- common.LogBuild("subDomain.get.parseIP",
			fmt.Sprintf("解析域名失败%s:%s", d.Domain, err.Error()), common.FAULT)
		return
	}
	for _, a := range r.Answer {
		if cn, ok := a.(*dns.CNAME); ok {
			if d.Cname != "" {
				d.Cname += "|"
			}
			d.Cname += cn.Target
		}

		if ip, ok := a.(*dns.A); ok {
			ipi := IpInfo{
				Ip: ip.A.String(),
			}
			s.parseIpLoc(&ipi)
			d.Ip = append(d.Ip, ipi)
		}
	}
	if len(d.Ip) > 0 {
		s.ErrChannel <- common.LogBuild("subDomain",
			fmt.Sprintf("解析域名%s完成CNAME=%s,ip数量%d", d.Domain, d.Cname, len(d.Ip)), common.ALERT)
	}
}

func (s *SubDomain) parseIpLoc(ipi *IpInfo) {
	for {
		if ipApiLimit.Allow() {
			break
		} else {
			time.Sleep(time.Millisecond * 1500)
		}
	}

	url := fmt.Sprintf("http://ip-api.com/json/%s?lang=zh-CN", ipi.Ip)
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
			fmt.Sprintf("获取IP%s所在地失败%s", ipi.Ip, err.Error()), common.FAULT)
		return
	}
	_ = json.Unmarshal(response.Body, ipi)
	time.Sleep(time.Second * 1)
}
