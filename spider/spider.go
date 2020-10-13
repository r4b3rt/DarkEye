package spider

import (
	"encoding/base64"
	"github.com/gocolly/colly"
	"github.com/zsdevX/DarkEye/common"
	"net/url"
	"regexp"
	"strings"
)

func (sp *Spider) Run() {
	sp.findInterface()
	if sp.SearchEnable {
		sp.Search()
	}
}

func (sp *Spider) findInterface() {
	allowDomain := ""
	if sp.LocalLink {
		url, err := url.Parse(sp.Url)
		if err != nil {
			sp.ErrChannel <- common.LogBuild("spider", err.Error(), common.FAULT)
			return
		}
		allowDomain = url.Host
	}

	disAllowedRequest, err := base64.StdEncoding.DecodeString(sp.DisAllowedRequest)
	if err != nil {
		sp.ErrChannel <- common.LogBuild("spider", err.Error(), common.FAULT)
		return
	}
	c := colly.NewCollector(
		colly.DisallowedURLFilters(
			regexp.MustCompile(string(disAllowedRequest)),
		),
		colly.AllowedDomains(allowDomain),
	)
	c.MaxDepth = sp.MaxDeps

	reqRules := strings.Split(sp.RequestMatchRule, ",")
	for _, rule := range reqRules {
		r := strings.Split(rule, ":")
		if len(r) != 2 {
			sp.ErrChannel <- common.LogBuild("spider", rule+"规则格式错误", common.FAULT)
			return
		}
		c.OnHTML(r[0]+"["+r[1]+"]", func(e *colly.HTMLElement) {
			_ = e.Request.Visit(e.Attr(r[1]))
		})
	}

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Cookie", sp.Cookie)
		if common.ShouldStop(&sp.Stop) {
			r.Abort()
		}
		//sp.ErrChannel <- common.LogBuild("", r.URL.String(), common.INFO)
	})

	responseMatchRule, err := base64.StdEncoding.DecodeString(sp.ResponseMatchRule)
	if err != nil {
		sp.ErrChannel <- common.LogBuild("spider", err.Error(), common.FAULT)
		return
	}
	responseFilter, err := base64.StdEncoding.DecodeString(sp.ResponseFilter)
	if err != nil {
		sp.ErrChannel <- common.LogBuild("spider", err.Error(), common.FAULT)
		return
	}
	reg := regexp.MustCompile(string(responseMatchRule))
	regFilter := regexp.MustCompile(string(responseFilter))

	c.OnResponse(func(resp *colly.Response) {
		matches := reg.FindAllString(string(resp.Body), -1)
		for _, url := range matches {
			if strings.Contains(url, "/") {
				if regFilter.MatchString(url) {
					continue
				}
				sp.ErrChannel <- common.LogBuild("spider", url, common.INFO)
			}
		}
	})
	_ = c.Visit(sp.Url)
}
