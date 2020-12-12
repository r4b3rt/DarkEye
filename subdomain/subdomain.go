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

func (s *SubDomain) Run() {
	queries := strings.Split(s.Queries, ",")
	s.dns = make([]dnsInfo, 0)
	for _, q := range queries {
		s.get(q)
	}
	w, file, fileName, err := common.CreateCSV("dns",
		[]string{"域名", "CNAME", "中间件", "标题", "IP"})
	if err != nil {
		s.ErrChannel <- common.LogBuild("fofa",
			fmt.Sprintf("创建记录文件失败:"+err.Error()), common.INFO)
		return
	}
	defer file.Close()

	logNumber := 0
	for _, n := range s.dns {
		ipi := ""
		for i, ip := range n.ip {
			if i != 0 {
				ipi += "|"
			}
			ipi += fmt.Sprintf("%s_%s_%s", ip.ip, ip.RegionName, ip.Isp)
		}
		_ = w.Write([]string{n.domain, n.cname, n.server, n.title, ipi})
		logNumber++
	}
	w.Flush()
	if logNumber == 0 {
		s.ErrChannel <- common.LogBuild("subDomain",
			fmt.Sprintf("收集信息任务完成，未有结果"), common.INFO)
	} else {
		s.ErrChannel <- common.LogBuild("subDomain",
			fmt.Sprintf("收集信息任务完成，有效数量%d, 已保存结果:%s", logNumber, fileName), common.INFO)
	}
}

func (s *SubDomain) get(query string) {
	s.ErrChannel <- common.LogBuild("subDomain",
		fmt.Sprintf("开始收集子域%s", query), common.INFO)

	if ! s.Brute {
		s.getSecurityTrails(query)
		return
	} else {
		s.getBrute(query)
	}
}

func (s *SubDomain) try(prefix, query string) {
	d := dnsInfo{
		domain: prefix + "." + query,
		ip:     make([]ipInfo, 0),
	}
	if s.IpCheck {
		s.parseIP(&d)
	}
	if len(d.ip) > 0 {
		s.parseTag(&d)
		s.dns = append(s.dns, d)
	}
}

func (s *SubDomain) parseTag(d *dnsInfo) {
	d.server, d.title, _ = common.GetHttpTitle("http", d.domain, 5)
	if d.server == "" && d.title == "" {
		d.server, d.title, _ = common.GetHttpTitle("https", d.domain, 5)
	}
	return
}

func (s *SubDomain) parseIP(d *dnsInfo) {
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
		s.ErrChannel <- common.LogBuild("subDomain.get.parseIP",
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
	if len(d.ip) > 0 {
		s.ErrChannel <- common.LogBuild("subDomain",
			fmt.Sprintf("解析域名%s完成CNAME=%s,ip数量%d", d.domain, d.cname, len(d.ip)), common.ALERT)
	}
}

func (s *SubDomain) parseIpLoc(ipi *ipInfo) {
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
