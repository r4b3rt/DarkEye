package main

import (
	"github.com/zsdevX/DarkEye/common"
	"github.com/zsdevX/DarkEye/subdomain"
	"github.com/zsdevX/DarkEye/ui"
	"strings"
)

func subDomainUICtl(mainWindow *ui.MainWindow, brute bool) {
	mainWindow.St_apikeylist.SetDisabled(brute)
	mainWindow.St_domain_brute_rate.SetDisabled(!brute)
	mainWindow.St_domain_brute.SetDisabled(!brute)
}

//LoadSubDomain add comment
func LoadSubDomain(mainWindow *ui.MainWindow) {
	mainWindow.St_apikeylist.SetCurrentText(mConfig.SubDomain.ApiKey)
	mainWindow.St_dns.SetText(mConfig.SubDomain.DnsServer)
	mainWindow.St_domain.SetText(mConfig.SubDomain.Queries)
	mainWindow.St_domain.SetText(mConfig.SubDomain.BruteRate)
	mainWindow.St_domain.SetText(mConfig.SubDomain.BruteLength)
	logC, runCtl := logChannel(mainWindow.St_log)

	subDomainUICtl(mainWindow, mConfig.SubDomain.Brute)
	mainWindow.St_domain_brute_mode.ConnectClicked(func(brute bool) {
		subDomainUICtl(mainWindow, brute)
	})
	//Action
	mainWindow.St_start.ConnectClicked(func(bool) {
		//保存配置
		mConfig.SubDomain = subdomain.NewConfig()
		mConfig.SubDomain.DnsServer = mainWindow.St_dns.Text()
		mConfig.SubDomain.Queries = mainWindow.St_domain.Text()
		mConfig.SubDomain.BruteRate = mainWindow.St_domain.Text()
		mConfig.SubDomain.BruteLength = mainWindow.St_domain.Text()
		mainWindow.St_domain_brute_mode.IsChecked()

		if mainWindow.St_apikeylist.CurrentText() != mConfig.SubDomain.ApiKey {
			mainWindow.St_apikeylist.AddItems([]string{mainWindow.St_apikeylist.CurrentText()})
		}
		mConfig.SubDomain.ApiKey = strings.Trim(mainWindow.St_apikeylist.CurrentText(), "\n")

		mConfig.SubDomain.ErrChannel = logC
		if err := saveCfg(); err != nil {
			logC <- common.LogBuild("UI", "保存配置失败:"+err.Error(), common.FAULT)
			return
		}
		//启动流程
		mainWindow.St_start.SetEnabled(false)
		mainWindow.St_stop.SetEnabled(true)
		common.StartIt(&mConfig.SubDomain.Stop)
		go func() {
			mConfig.SubDomain.Run()
			mainWindow.St_start.SetEnabled(true)
			mainWindow.St_stop.SetEnabled(false)
			runCtl <- false
		}()
	})

	mainWindow.St_stop.ConnectClicked(func(bool) {
		common.StopIt(&mConfig.SubDomain.Stop)
		go func() {
			gracefulStop(mainWindow.St_start, mainWindow.St_stop, runCtl)
		}()

	})

	mainWindow.St_clear.ConnectClicked(func(bool) {
		mainWindow.St_log.Clear()
		mainWindow.St_log.SetText("")
	})
	return
}
