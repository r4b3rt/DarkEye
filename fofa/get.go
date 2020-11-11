package fofa

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/antchfx/htmlquery"
	"github.com/zsdevX/DarkEye/common"
	"math/rand"
	"regexp"
	"strconv"
	"time"
)

func (f *Fofa) get(query string) {
	cookie := f.FofaSession
	userAgent := common.UserAgents[rand.Int()%len(common.UserAgents)]
	req := common.Http{
		Agent:   userAgent,
		Cookie:  cookie,
		TimeOut: time.Duration(5 + f.Interval),
		Method:  "GET",
	}
	//获取首页面
	body, err := f.fetchBody(&req, query, 1)
	if err != nil {
		return
	}
	defer func() {
		//学做人，防止fofa封
		time.Sleep(time.Second * time.Duration(common.GenHumanSecond(f.Interval)))
	}()
	//获取页数
	pageRe, err := regexp.Compile(`>(\d*)</a> <a class="next_page" rel="next"`)
	if err != nil {
		f.ErrChannel <- common.LogBuild("fofa.get",
			fmt.Sprintf("%s: 无页码", query), common.ALERT)
		return
	}
	pageNr := 1
	pageNum := pageRe.FindSubmatch(body)
	if len(pageNum) < 1 {
		pageNr = 0
	} else {
		pageNr, _ = strconv.Atoi(string(pageNum[1]))
	}
	//非授权的只能获取f.Pages=5页
	if pageNr > f.Pages {
		pageNr = f.Pages
	}
	//解析页面
	start := 1
	for {
		if common.ShouldStop(&f.Stop) {
			break
		}
		//学做人，防止fofa封
		time.Sleep(time.Second * time.Duration(common.GenHumanSecond(f.Interval)))
		if f.parseHtml(query, body, start) {
			//解析页面遇到不可恢复的情况立刻终止，提高效率
			break
		}
		//下一页
		start += 1
		if start > pageNr {
			break
		}
		body, err = f.fetchBody(&req, query, start)
		if err != nil {
			return
		}
	}
}

func (f *Fofa) fetchBody(req *common.Http, query string, page int) (body []byte, err error) {
	req.Url = f.genUrl(query, page)
	retry := 0
	for {
		body, err = req.Http()
		if err != nil {
			errMsg := err.Error()
			if errMsg == "Bad status 429" {
				retry++
				f.ErrChannel <- common.LogBuild("fofa.get",
					fmt.Sprintf("%s: 提取信息失败，失败原因fofa session过期或刷新过快, 尝试等待1分钟重试, 当前尝试次数%d次",
						req.Url, retry), common.FAULT)
				time.Sleep(time.Second * 60)
			} else {
				retry++
				f.ErrChannel <- common.LogBuild("fofa.get",
					fmt.Sprintf("因网络质量问题%s获取信息失败 尝试次数（第%d次):%s", query, retry, errMsg), common.ALERT)
				time.Sleep(time.Second * time.Duration(common.GenHumanSecond(f.Interval)))
			}
			continue
		}
		return
	}
}

func (f *Fofa) genUrl(query string, page int) string {
	url := fmt.Sprintf("https://fofa.so/result?qbase64=%s&full=true&page=%d",
		base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("ip=%s", query))), page)
	return url
}

func (f *Fofa) parseHtml(query string, body []byte, page int) (stop bool) {
	stop = false
	doc, err := htmlquery.Parse(bytes.NewReader(body))
	if err != nil {
		f.ErrChannel <- common.LogBuild("fofa.get.parseIPHtml",
			fmt.Sprintf("%s: %s", query, err.Error()), common.ALERT)
		return
	}
	blocks := htmlquery.Find(doc, "//*[@class='right-list-view-item clearfix']")
	if len(blocks) == 0 {
		f.ErrChannel <- common.LogBuild("fofa",
			fmt.Sprintf("%s: 完成第%d页解析(无信息)", query, page), common.INFO)
		stop = true
		return
	}
	for _, blk := range blocks {
		node := ipNode{
			Ip: query,
		}
		//获取超链接
		items := htmlquery.Find(blk, "//*[@class='re-domain']/a[@href]")
		if items == nil {
			items = htmlquery.Find(blk, "//*[@class='re-domain']")
			node.Domain = common.TrimUseless(htmlquery.InnerText(items[0]))
		} else {
			node.Domain = htmlquery.SelectAttr(items[0], "href")
		}
		//获取网站标题
		items = htmlquery.Find(blk, "//*[@class='fl box-sizing']/div[2]")
		node.Title = htmlquery.InnerText(items[0])
		//获取中间件
		items = htmlquery.Find(blk, "//*[@class='com-tag-wrap clearfix']/a[@title]")
		for k, item := range items {
			if k != 0 {
				node.Server += "|"
			}
			node.Server += htmlquery.SelectAttr(item, "title")
		}
		//获取指纹
		items = htmlquery.Find(blk, "//*[@class='scroll-wrap-res']")
		node.Finger = common.TrimUseless(htmlquery.InnerText(items[0]))

		//获取端口
		items = htmlquery.Find(blk, "//*[@class='re-port ar']/a")
		node.Port = common.TrimUseless(htmlquery.InnerText(items[0]))

		//检查端口是否有效
		node.Alive = common.IsAlive(node.Ip, node.Port, 2000)

		//保存结果
		f.ipNodes = append(f.ipNodes, node)
	}
	f.ErrChannel <- common.LogBuild("fofa",
		fmt.Sprintf("%s: 完成第%d页解析", query, page), common.INFO)
	return
}
