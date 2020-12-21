package main

import (
	"github.com/zsdevX/DarkEye/common"
	"github.com/zsdevX/DarkEye/fofa"
	"github.com/zsdevX/DarkEye/ui"
	"github.com/zsdevX/DarkEye/zoomeye"
	"strconv"
)

func LoadAsset(mainWindow *ui.MainWindow) {
	assetBanner(mainWindow)
	mainWindow.Fofa_session.SetText(mConfig.Fofa.FofaSession)
	mainWindow.Fofa_asset_ip.SetText(mConfig.Fofa.Ip)
	mainWindow.Fofa_interval.SetText(strconv.Itoa(mConfig.Fofa.Interval))
	mainWindow.Zoomeuye_search.SetText(mConfig.Zoomeye.Query)
	mainWindow.Zoomeye_key.SetCurrentText(mConfig.Zoomeye.ApiKey)
	mainWindow.Zoomeye_radioButton.SetCheckable(true)

	logC, runCtl := logChannel(mainWindow.Fofa_log)
	//Action
	mainWindow.Fofa_start.ConnectClicked(func(bool) {
		//保存配置
		if mainWindow.Fofa_radioButton.IsChecked() {
			mConfig.Fofa = fofa.NewConfig()
			mConfig.Fofa.Ip = mainWindow.Fofa_asset_ip.Text()
			mConfig.Fofa.Interval, _ = strconv.Atoi(mainWindow.Fofa_interval.Text())
			mConfig.Fofa.FofaSession = mainWindow.Fofa_session.Text()
			mConfig.Fofa.ErrChannel = logC
			if err := saveCfg(); err != nil {
				logC <- common.LogBuild("UI", "保存配置失败:"+err.Error(), common.FAULT)
				return
			}
			common.StartIt(&mConfig.Fofa.Stop)
		} else {
			mConfig.Zoomeye = zoomeye.New()
			mConfig.Zoomeye.ApiKey = mainWindow.Zoomeye_key.CurrentText()
			mConfig.Zoomeye.Query = mainWindow.Zoomeuye_search.Text()
			mConfig.Zoomeye.ErrChannel = logC
			if err := saveCfg(); err != nil {
				logC <- common.LogBuild("UI", "保存配置失败:"+err.Error(), common.FAULT)
				return
			}
		}
		//启动流程
		mainWindow.Fofa_start.SetEnabled(false)
		mainWindow.Fofa_stop.SetEnabled(true)
		go func() {
			if mainWindow.Fofa_radioButton.IsChecked() {
				mConfig.Fofa.Run()
			} else {
				mConfig.Zoomeye.Run()
			}
			mainWindow.Fofa_start.SetEnabled(true)
			mainWindow.Fofa_stop.SetEnabled(false)
			runCtl <- false
		}()
	})

	mainWindow.Fofa_stop.ConnectClicked(func(bool) {
		common.StopIt(&mConfig.Fofa.Stop)

		go func() {
			gracefulStop(mainWindow.Fofa_start, mainWindow.Fofa_stop, runCtl)
		}()

	})

	mainWindow.Fofa_clear.ConnectClicked(func(bool) {
		mainWindow.Fofa_log.Clear()
		mainWindow.Fofa_log.SetText("")
	})
}

func assetBanner(mainWindow *ui.MainWindow) {
	mainWindow.Fofa_log.SetText(`
	关于ZoomEye: 
		* 速度快!
		* 搜索语法请参考 https://www.zoomeye.org
		* API-KEY登录 https://www.zoomeye.org/profile 获取,每月免费10w/资源!
`)
}
