package main

import (
	"fmt"
	"github.com/zsdevX/DarkEye/common"
	"github.com/zsdevX/DarkEye/securitytrails"
	"github.com/zsdevX/DarkEye/ui"
	"time"
)

func LoadSecurityTrails(mainWindow *ui.MainWindow) {
	mainWindow.St_apikeylist.SetCurrentText(mConfig.SecurityTrails.ApiKey)
	mainWindow.St_dns.SetText(mConfig.SecurityTrails.DnsServer)
	mainWindow.St_domain.SetText(mConfig.SecurityTrails.Queries)
	logC, runCtl := logChannel(mainWindow.St_log)
	//Action
	mainWindow.St_start.ConnectClicked(func(bool) {
		//保存配置
		mConfig.SecurityTrails = securitytrails.NewConfig()
		mConfig.SecurityTrails.DnsServer = mainWindow.St_dns.Text()
		mConfig.SecurityTrails.Queries = mainWindow.St_domain.Text()

		if mainWindow.St_apikeylist.CurrentText() != mConfig.SecurityTrails.ApiKey {
			mainWindow.St_apikeylist.AddItems([]string{mainWindow.St_apikeylist.CurrentText()})
		}
		mConfig.SecurityTrails.ApiKey = mainWindow.St_apikeylist.CurrentText()

		mConfig.SecurityTrails.ErrChannel = logC
		if err := saveCfg(); err != nil {
			logC <- common.LogBuild("UI", "保存配置失败:"+err.Error(), common.FAULT)
			return
		}
		//启动流程
		mainWindow.St_start.SetEnabled(false)
		mainWindow.St_stop.SetEnabled(true)
		common.StartIt(&mConfig.SecurityTrails.Stop)
		go func() {
			mConfig.SecurityTrails.Run()
			mainWindow.St_start.SetEnabled(true)
			mainWindow.St_stop.SetEnabled(false)
			runCtl <- false
		}()
	})

	mainWindow.St_stop.ConnectClicked(func(bool) {
		common.StopIt(&mConfig.SecurityTrails.Stop)
		mainWindow.St_stop.SetDisabled(true)
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
					mainWindow.St_stop.SetText(fmt.Sprintf("等待%d秒", 60-sec))
				}
				if stop {
					break
				}
			}
			mainWindow.St_start.SetEnabled(true)
			mainWindow.St_stop.SetText("停止")
		}()
	})

	mainWindow.St_clear.ConnectClicked(func(bool) {
		mainWindow.St_log.Clear()
	})
	return
}
