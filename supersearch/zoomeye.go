package main

import (
	"context"
	"encoding/json"
	"github.com/b1gcat/DarkEye/common"
	"github.com/b1gcat/DarkEye/supersearch/ui"
	"github.com/b1gcat/DarkEye/supersearch/zoomeye"
	"path/filepath"
	"strconv"
	"time"
)

func zoomEyeInit() {
	ctx, cancel := context.WithCancel(context.Background())
	z := ui.NewZoomEye(nil)
	log := newLogChannel(z.TextEditLog)
	//设置搜索类型联动
	zoomEyeTypeChanged(z)
	z.ComboBoxType.ConnectCurrentTextChanged(func(string) {
		zoomEyeTypeChanged(z)
	})
	//设置搜索帮助
	z.ToolButton.ConnectClicked(func(bool) {
		log <- `
Samples:
Server: Apache/2.4.49
app:"apache web server 2.4.49 2.4.50"

Get More: 
host dork:[app,version,device,port,city,country,asn,banner,timestamp,*]
web dork: [app,headers,keywords,title,site,city,country,webapp,component,framework,server,waf,os,timestamp,*]
https://www.zoomeye.org/
`
	})
	//设置域名提示
	z.LineEditDomain.SetPlaceholderText("0x0x.com")
	//设置搜索
	z.PushButtonStart.ConnectClicked(func(bool) {
		z.PushButtonStart.SetEnabled(false)
		go func() {
			n, _ := strconv.Atoi(z.SpinBoxNumber.Text())
			m := zoomeye.Run(ctx, n,
				z.Key.Text(),
				z.ComboBoxType.CurrentText()+"@"+z.PlainTextEditSearch.ToPlainText(),
				z.LineEditFacet.Text(),
				log)
			if m != nil {
				f := filepath.Join(defaultOutputDir, "zoomEye_"+time.Now().Format("20060102150405")+".csv")
				d, _ := json.Marshal(m)
				if err := output(d, f); err != nil {
					common.LogUi(err.Error(), log, common.FAULT)
				} else {
					common.LogUi("保存文件:"+f, log, common.FAULT)
				}
			}
			z.PushButtonStart.SetEnabled(true)
		}()
	})
	z.PushButtonStop.ConnectClicked(func(bool) {
		cancel()
	})
	z.Show()
}

func zoomEyeTypeChanged(z *ui.ZoomEye) {
	if z.ComboBoxType.CurrentText() == "Host" {
		z.LineEditFacet.SetText("app,device,service,os,port,country,city")
	} else {
		z.LineEditFacet.SetText("webapp,component,framework,server,waf,os,country")
	}
}
