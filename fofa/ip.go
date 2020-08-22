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

func (f *Fofa) doIp(ip string) {
	url := fmt.Sprintf("https://fofa.so/result?qbase64=%s&full=true",
		base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("ip=%s", ip))))
	cookie := "_fofapro_ars_session" + "=" + f.FofaSession
	userAgent := common.UserAgents[rand.Int()%len(common.UserAgents)]
	req := common.Http{
		Agent:   userAgent,
		Cookie:  cookie,
		Url:     url,
		TimeOut: time.Duration(5 + f.Interval),
		Method:  "GET",
		Referer: url,
	}
	//获取首页面
	body, err := req.Http()
	if err != nil {
		f.ErrChannel <- common.LogBuild("fofa.doip",
			fmt.Sprintf("收集IP%s失败:%s", ip, err.Error()), common.ALERT)
		return
	}
	//获取页数
	pageRe, err := regexp.Compile(`>(\d*)</a> <a class="next_page" rel="next"`)
	if err != nil {
		f.ErrChannel <- common.LogBuild("fofa.doip",
			fmt.Sprintf("收集IP%s失败:为匹配到页码", ip), common.FAULT)
		return
	}
	pageNum := pageRe.FindSubmatch(body)
	if len(pageNum) < 2 {
		f.ErrChannel <- common.LogBuild("fofa.doip",
			fmt.Sprintf("收集IP%s失败:匹配到错误页码", ip), common.FAULT)
		return
	}
	pageNr, _ := strconv.Atoi(string(pageNum[1]))
	//非授权的只能获取f.Pages=5页
	if pageNr > f.Pages {
		pageNr = f.Pages
	}
	//解析页面
	start := 1
	for {
		//学做人，防止fofa封
		time.Sleep(time.Second * time.Duration(common.GenHumanSecond(f.Interval)))
		//耗时点增加关停
		if common.ShouldStop() {
			break
		}
		f.parseIPHtml(ip, body, start)
		//NextPage
		start += 1
		if start > pageNr {
			break
		}
		url := fmt.Sprintf("https://fofa.so/result?qbase64=%s&full=true&page=%d",
			base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("ip=%s", ip))), start)
		req.Url = url
		//拉页面
		body, err = req.Http()
		if err != nil {
			f.ErrChannel <- common.LogBuild("fofa.doip", fmt.Sprintf("收集IP%s失败:%s", ip, err.Error()), common.ALERT)
			return
		}
	}
}

func (f *Fofa) parseIPHtml(ip string, body []byte, page int) {
	doc, err := htmlquery.Parse(bytes.NewReader(body))
	if err != nil {
		f.ErrChannel <- common.LogBuild("fofa.doip.parseIPHtml",
			fmt.Sprintf("解析失败:%s:%s", ip, err.Error()), common.ALERT)
		return
	}
	blocks := htmlquery.Find(doc, "//*[@class='right-list-view-item clearfix']")
	if len(blocks) == 0 {
		f.ErrChannel <- common.LogBuild("Fofa",
			fmt.Sprintf("%s:完成第%d页解析(无信息，请检查登录session是否过期或有效)", ip, page), common.ALERT)
		return
	}
	for _, blk := range blocks {
		node := ipNode{
			Ip: ip,
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
		node.Alive = common.IsAlive(node.Ip, node.Port)

		//保存结果
		f.ipNodes = append(f.ipNodes, node)
	}
	f.ErrChannel <- common.LogBuild("Fofa",
		fmt.Sprintf("%s:完成第%d页解析", ip, page), common.INFO)
}
