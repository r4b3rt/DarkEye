package main

import (
	"github.com/zsdevX/DarkEye/common"
	"github.com/zsdevX/DarkEye/fofa"
	"github.com/zsdevX/DarkEye/ui"
	"strconv"
)

func LoadFoFa(mainWindow *ui.MainWindow) {
	mainWindow.Fofa_session.SetText(mConfig.Fofa.FofaSession)
	mainWindow.Fofa_asset_ip.SetText(mConfig.Fofa.Ip)
	mainWindow.Fofa_interval.SetText(strconv.Itoa(mConfig.Fofa.Interval))

	logC, runCtl := logChannel(mainWindow.Fofa_log)
	//Action
	mainWindow.Fofa_start.ConnectClicked(func(bool) {
		//保存配置
		mConfig.Fofa = fofa.NewConfig()
		mConfig.Fofa.Ip = mainWindow.Fofa_asset_ip.Text()
		mConfig.Fofa.Interval, _ = strconv.Atoi(mainWindow.Fofa_interval.Text())
		mConfig.Fofa.FofaSession = mainWindow.Fofa_session.Text()
		mConfig.Fofa.ErrChannel = logC
		if err := saveCfg(); err != nil {
			logC <- common.LogBuild("UI", "保存配置失败:"+err.Error(), common.FAULT)
			return
		}
		//启动流程
		mainWindow.Fofa_start.SetEnabled(false)
		mainWindow.Fofa_stop.SetEnabled(true)
		common.StartIt(&mConfig.Fofa.Stop)
		go func() {
			mConfig.Fofa.Run()
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
	return
}
