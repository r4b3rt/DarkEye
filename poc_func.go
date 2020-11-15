package main

import (
	"fmt"
	"github.com/therecipe/qt/widgets"
	"github.com/zsdevX/DarkEye/common"
	"github.com/zsdevX/DarkEye/hack/poc"
	"github.com/zsdevX/DarkEye/ui"
	"time"
)

func LoadPoc(mainWindow *ui.MainWindow) {

	mainWindow.Poc_url.SetText(mConfig.Poc.Urls)
	mainWindow.Poc_reverse_url.SetText(mConfig.Poc.ReverseUrl)
	mainWindow.Poc_reverse_url_check.SetText(mConfig.Poc.ReverseCheckUrl)
	mainWindow.Poc_file.SetText(mConfig.Poc.FileName)

	logC, runCtl := logChannel(mainWindow.Poc_log)
	//Action
	mainWindow.Poc_start.ConnectClicked(func(bool) {
		//清空
		mainWindow.Poc_log.Clear()
		//保存配置
		mConfig.Poc = poc.NewConfig()
		mConfig.Poc.Urls = mainWindow.Poc_url.Text()
		mConfig.Poc.ReverseUrl = mainWindow.Poc_reverse_url.Text()
		mConfig.Poc.ReverseCheckUrl = mainWindow.Poc_reverse_url_check.Text()
		mConfig.Poc.FileName = mainWindow.Poc_file.Text()

		mConfig.Poc.ErrChannel = logC
		if err := saveCfg(); err != nil {
			logC <- common.LogBuild("UI", "保存配置失败:"+err.Error(), common.FAULT)
			return
		}
		//启动流程
		mainWindow.Poc_start.SetEnabled(false)
		mainWindow.Poc_stop.SetEnabled(true)
		common.StartIt(&mConfig.Poc.Stop)
		go func() {
			mConfig.Poc.Check()
			mainWindow.Poc_start.SetEnabled(true)
			mainWindow.Poc_stop.SetEnabled(false)
			runCtl <- false
		}()
	})

	mainWindow.Poc_stop.ConnectClicked(func(bool) {
		common.StopIt(&mConfig.Poc.Stop)
		mainWindow.Poc_stop.SetDisabled(true)
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
					mainWindow.Poc_stop.SetText(fmt.Sprintf("等待%d秒", 60-sec))
				}
				if stop {
					break
				}
			}
			mainWindow.Poc_start.SetEnabled(true)
			mainWindow.Poc_stop.SetText("停止")
		}()
	})

	mainWindow.Poc_import_urls.ConnectClicked(func(bool) {
		qFile := widgets.NewQFileDialog2(nil, "选择url列表文件", ".", "")
		fn := qFile.GetOpenFileName(nil, "文件", ".", "", "", widgets.QFileDialog__ReadOnly)
		if fn == "" {
			return
		}
		urls, err := common.ImportFiles(fn, mainWindow.Poc_url.Text())
		if err != nil {
			logC <- common.LogBuild("UI", "加载文件失败"+err.Error(), common.ALERT)
			return
		}
		mainWindow.Poc_url.SetText(urls)
		logC <- common.LogBuild("UI", "批量导入完成", common.INFO)
	})

	mainWindow.Poc_import_file.ConnectClicked(func(bool) {
		qFile := widgets.NewQFileDialog2(nil, "选择POC列表文件夹", "", "")
		qFile.SetFileMode(widgets.QFileDialog__DirectoryOnly)
		if qFile.Exec() != int(widgets.QDialog__Accepted) {
			return
		}
		mainWindow.Poc_file.SetText(qFile.SelectedFiles()[0])
		logC <- common.LogBuild("UI", "设置完成", common.INFO)
	})
	return
}
