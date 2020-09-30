package main

import (
	"fmt"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
	"github.com/zsdevX/DarkEye/common"
	"github.com/zsdevX/DarkEye/securitytrails"
	"time"
)

var (
	//功能界面
	windowSecurityTails = &widgets.QMainWindow{}
)

func registerSecurityTrails() (window *widgets.QMainWindow) {
	window = widgets.NewQMainWindow(nil, 0)
	//Input
	apiKey := widgets.NewQLineEdit(nil)
	apiKey.SetPlaceholderText("apiKey")
	apiKey.SetAlignment(core.Qt__AlignHCenter)
	if mConfig.SecurityTrails.ApiKey != "" {
		apiKey.SetText(mConfig.SecurityTrails.ApiKey)
	}

	dnsServer := widgets.NewQLineEdit(nil)
	dnsServer.SetPlaceholderText("DNS服务器(格式:8.8.8.8:53)")
	dnsServer.SetAlignment(core.Qt__AlignHCenter)
	if mConfig.SecurityTrails.DnsServer != "" {
		dnsServer.SetText(mConfig.SecurityTrails.DnsServer)
	}

	queries := widgets.NewQLineEdit(nil)
	queries.SetPlaceholderText("a.com,b.com,c.com")
	queries.SetAlignment(core.Qt__AlignHCenter)

	checkBox := widgets.NewQCheckBox2("解析域名为IP", nil)
	checkBox.SetChecked(true)

	//Log
	logC, runCtl, inputLog := getWindowCtl()

	btnOpen := widgets.NewQPushButton2("启动", nil)
	btnStop := widgets.NewQPushButton2("停止", nil)
	btnStop.SetDisabled(true)

	widgetC := widgets.NewQWidget(nil, 0)
	widgetC.SetLayout(widgets.NewQHBoxLayout())
	widgetC.Layout().AddWidget(queries)
	widgetC.Layout().AddWidget(checkBox)

	widgetD := widgets.NewQWidget(nil, 0)
	widgetD.SetLayout(widgets.NewQHBoxLayout())
	widgetD.Layout().AddWidget(apiKey)
	widgetD.Layout().AddWidget(dnsServer)

	widgetA := widgets.NewQWidget(nil, 0)
	widgetA.SetLayout(widgets.NewQVBoxLayout())
	widgetA.Layout().AddWidget(inputLog)

	widgetB := widgets.NewQWidget(nil, 0)
	widgetB.SetLayout(widgets.NewQHBoxLayout())
	widgetB.Layout().AddWidget(btnOpen)
	widgetB.Layout().AddWidget(btnStop)

	widget := widgets.NewQWidget(nil, 0)
	widget.SetLayout(widgets.NewQVBoxLayout())
	widget.Layout().AddWidget(widgetC)
	widget.Layout().AddWidget(widgetD)
	widget.Layout().AddWidget(widgetA)
	widget.Layout().AddWidget(widgetB)

	window.SetMinimumSize2(650, 480)
	window.SetWindowTitle(programDesc)
	window.SetCentralWidget(widget)
	window.AutoFillBackground()
	window.SetWindowFlags(core.Qt__Dialog)

	//Action
	btnOpen.ConnectClicked(func(bool) {
		//保存配置
		mConfig.SecurityTrails = securitytrails.NewConfig()
		mConfig.SecurityTrails.Queries = queries.Text()
		mConfig.SecurityTrails.ApiKey = apiKey.Text()
		mConfig.SecurityTrails.DnsServer = dnsServer.Text()
		mConfig.SecurityTrails.IpCheck = checkBox.IsChecked()

		mConfig.SecurityTrails.ErrChannel = logC
		if err := saveCfg(); err != nil {
			logC <- common.LogBuild("UI", "保存配置失败:"+err.Error(), common.FAULT)
			return
		}
		//启动流程
		common.StartIt(&mConfig.SecurityTrails.Stop)

		go func() {
			mConfig.SecurityTrails.Run()
			btnStop.SetDisabled(true)
			btnOpen.SetDisabled(false)
			runCtl <- false
		}()
		btnStop.SetDisabled(false)
		btnOpen.SetDisabled(true)
		window.SetWindowState(core.Qt__WindowNoState)
	})

	btnStop.ConnectClicked(func(bool) {
		common.StopIt(&mConfig.SecurityTrails.Stop)
		btnStop.SetDisabled(true)
		//异步处理等待结束避免界面卡顿
		go func() {
			sec := 0
			stop := false
			tick := time.NewTicker(time.Second)
			for {
				select {
				case <-runCtl:
					stop = true
				case <-tick.C:
					sec ++
					btnStop.SetText(fmt.Sprintf("等待%d秒", 60-sec))
				}
				if stop {
					break
				}
			}
			btnOpen.SetDisabled(false)
			btnStop.SetText("停止")
		}()
	})
	return
}
