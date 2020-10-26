package spider

import (
	"encoding/base64"
	"github.com/gocolly/colly"
	"github.com/zsdevX/DarkEye/common"
	"net/url"
	"regexp"
	"strings"
)

func (sp *Spider) ApiFinder() {
	c := sp.setup()
	if c == nil {
		return
	}
	//设置Cookie
	c.OnRequest(func(r *colly.Request) {
		if sp.Cookie != "" {
			r.Headers.Set("Cookie", sp.Cookie)
		}
		if common.ShouldStop(&sp.Stop) {
			r.Abort()
		}
	})

	sp.requestLinkExtrack(c)
	sp.responseResultExtract(c)

	_ = c.Visit(sp.Url)
}

func (sp *Spider) responseResultExtract(c *colly.Collector) {
	//设置贪婪匹配规则
	responseMatchRule, err := base64.StdEncoding.DecodeString(sp.ResponseMatchRule)
	if err != nil {
		sp.ErrChannel <- common.LogBuild("spider", err.Error(), common.FAULT)
		return
	}
	reg := regexp.MustCompile(string(responseMatchRule))
	//设置过滤数据规则
	responseFilter, err := base64.StdEncoding.DecodeString(sp.ResponseFilter)
	if err != nil {
		sp.ErrChannel <- common.LogBuild("spider", err.Error(), common.FAULT)
		return
	}
	regFilter := regexp.MustCompile(string(responseFilter))
	//Apply
	c.OnResponse(func(resp *colly.Response) {
		//匹配贪婪提取敏感路径和接口
		matches := reg.FindAllString(string(resp.Body), -1)
		for _, url := range matches {
			if strings.Contains(url, "/") {
				//匹配结果中去除一些垃圾数据
				if regFilter.MatchString(url) {
					continue
				}
				sp.ErrChannel <- common.LogBuild("spider", url, common.INFO)
			}
		}
	})
}

func (sp *Spider) requestLinkExtrack(c *colly.Collector) {
	reqRules := strings.Split(sp.RequestMatchRule, ",")
	for _, rule := range reqRules {
		r := strings.Split(rule, ":")
		if len(r) != 2 {
			sp.ErrChannel <- common.LogBuild("spider", rule+"规则格式错误", common.FAULT)
			return
		}
		//设置爬取的链接提取方式：例如a[href]表示<a href="ooxx.com"></a>
		c.OnHTML(r[0]+"["+r[1]+"]", func(e *colly.HTMLElement) {
			_ = e.Request.Visit(e.Attr(r[1]))
		})
	}
}

func (sp *Spider) setup() *colly.Collector {
	allowDomain := ""
	if sp.LocalLink {
		url, err := url.Parse(sp.Url)
		if err != nil {
			sp.ErrChannel <- common.LogBuild("spider", err.Error(), common.FAULT)
			return nil
		}
		allowDomain = url.Host
	}

	disAllowedRequest, err := base64.StdEncoding.DecodeString(sp.DisAllowedRequest)
	if err != nil {
		sp.ErrChannel <- common.LogBuild("spider", err.Error(), common.FAULT)
		return nil
	}
	c := colly.NewCollector(
		colly.DisallowedURLFilters(
			regexp.MustCompile(string(disAllowedRequest)),
		),
		colly.AllowedDomains(allowDomain),
	)
	c.MaxDepth = sp.MaxDeps
	return c
}
