package spider

import (
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"github.com/zsdevX/DarkEye/common"
	"net/url"
	"regexp"
	"strings"
)

//ApiFinder add comment
func (sp *Spider) ApiFinder() {
	myUrls := strings.Split(sp.Url, ",")
	sp.sensitiveInterface = make([]SensitiveInterface, 0)
	for _, myUrl := range myUrls {
		sp.apiFinderUrl(myUrl)
	}

	if len(sp.sensitiveInterface) == 0 {
		sp.ErrChannel <- common.LogBuild("spider",
			fmt.Sprintf("获取信息%s:%s", sp.Url, "无信息"), common.INFO)
		return
	}
	res, _ := json.Marshal(sp.sensitiveInterface)

	filename, err := common.Write2CSV(sp.Url+"spider", res)
	if err != nil {
		sp.ErrChannel <- common.LogBuild("spider",
			fmt.Sprintf("获取信息%s:%s", sp.Url, err.Error()), common.FAULT)
		return
	}
	sp.ErrChannel <- common.LogBuild("spider",
		fmt.Sprintf("收集信息任务完成，有效数量%d, 已保存结果:%s", len(sp.sensitiveInterface), filename), common.INFO)
}

func (sp *Spider) apiFinderUrl(myUrl string) {
	c := sp.setup(myUrl)
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

	sp.requestLinkExtract(c)
	sp.responseResultExtract(c)
	err := c.Visit(myUrl)
	if err != nil {
		sp.ErrChannel <- common.LogBuild("spider", sp.Url+":"+err.Error(), common.FAULT)
	}
}

func (sp *Spider) responseResultExtract(c *colly.Collector) {
	reg := regexp.MustCompile(sp.ResponseMatchRule)
	//Apply
	c.OnResponse(func(resp *colly.Response) {
		//匹配贪婪提取敏感路径和接口
		matches := reg.FindAllString(string(resp.Body), -1)
		for _, url := range matches {
			sp.ErrChannel <- common.LogBuild("spider", url, common.INFO)
			sp.sensitiveInterface = append(sp.sensitiveInterface, guessUrlLevel(url))
		}
	})
}

//大佬可以自己增加判断等级的方式
func guessUrlLevel(url string) SensitiveInterface {
	s := SensitiveInterface{
		Level: 0,
		API:   url,
	}
	reg := regexp.MustCompile(`\?\w+=`)
	if reg.MatchString(url) {
		s.Level++
	}
	reg = regexp.MustCompile(`file|download|upload|url`)
	if reg.MatchString(url) {
		s.Level++
	}
	reg = regexp.MustCompile(`http://|https://`)
	if reg.MatchString(url) {
		s.Level++
	}
	return s
}

func (sp *Spider) requestLinkExtract(c *colly.Collector) {
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

func (sp *Spider) setup(myUrl string) *colly.Collector {
	allowDomain := ""
	if sp.LocalLink {
		url, err := url.Parse(myUrl)
		if err != nil {
			sp.ErrChannel <- common.LogBuild("spider", err.Error(), common.FAULT)
			return nil
		}
		allowDomain = strings.Split(url.Host, ":")[0]
	}
	c := colly.NewCollector(
		colly.DisallowedURLFilters(
			regexp.MustCompile(sp.DisAllowedRequest),
		),
		colly.AllowedDomains(allowDomain),
	)
	c.MaxDepth = sp.MaxDeps
	return c
}
