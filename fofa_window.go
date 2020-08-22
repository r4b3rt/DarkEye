package main

import (
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
	"github.com/zsdevX/DarkEye/common"
	"strconv"
	"time"
	"fmt"
)

var (
	//fofa功能界面
	windowFofa = &widgets.QMainWindow{}
)

func registerFofa(sysTray *QSystemTrayIconWithCustomSlot) (window *widgets.QMainWindow) {
	window = widgets.NewQMainWindow(nil, 0)
	//Input
	ip := widgets.NewQLineEdit(nil)
	ip.SetPlaceholderText("IP（Nmap格式）")
	ip.SetAlignment(core.Qt__AlignHCenter)

	Interval := widgets.NewQLineEdit(nil)
	Interval.SetPlaceholderText("检索间隔（建议10秒）")
	Interval.SetAlignment(core.Qt__AlignHCenter)

	session := widgets.NewQLineEdit(nil)
	session.SetPlaceholderText("Fofa session")
	session.SetAlignment(core.Qt__AlignHCenter)

	//Log
	inputLog := NewCustomEditor(nil)
	inputLog.SetReadOnly(true)
	listenLog(inputLog)

	btnOpen := widgets.NewQPushButton2("启动", nil)
	btnStop := widgets.NewQPushButton2("停止", nil)
	btnStop.SetDisabled(true)

	widgetA := widgets.NewQWidget(nil, 0)
	widgetA.SetLayout(widgets.NewQVBoxLayout())
	widgetA.Layout().AddWidget(ip)
	widgetA.Layout().AddWidget(Interval)
	widgetA.Layout().AddWidget(session)
	widgetA.Layout().AddWidget(inputLog)

	widgetB := widgets.NewQWidget(nil, 0)
	widgetB.SetLayout(widgets.NewQHBoxLayout())
	widgetB.Layout().AddWidget(btnOpen)
	widgetB.Layout().AddWidget(btnStop)

	widget := widgets.NewQWidget(nil, 0)
	widget.SetLayout(widgets.NewQVBoxLayout())
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
		mConfig.Fofa.Ip = ip.Text()
		mConfig.Fofa.Interval, _ = strconv.Atoi(Interval.Text())
		mConfig.Fofa.FofaSession = session.Text()
		mConfig.Fofa.ErrChannel = logC
		if err := saveCfg(); err != nil {
			sendUILog(common.LogBuild("UI", "保存配置失败:"+err.Error(), common.FAULT))
			return
		}
		//启动流程
		common.StartIt()
		go func() {
			sysTray.TriggerSlot()
		}()
		go func() {
			mConfig.Fofa.Run()
			btnStop.SetDisabled(true)
			btnOpen.SetDisabled(false)
			runCtl <- false
		}()
		btnStop.SetDisabled(false)
		btnOpen.SetDisabled(true)
		window.SetWindowState(core.Qt__WindowNoState)
	})

	btnStop.ConnectClicked(func(bool) {
		common.StopIt()
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
