package main

import (
	"github.com/b1gcat/DarkEye/supersearch/ui"
	"github.com/therecipe/qt/widgets"
	"os"
	"path/filepath"
	"runtime"
)

func preferenceInit() {
	workDir := "."
	p := ui.NewDialogPreference(nil)

	homeDir, _ := os.UserHomeDir()
	if runtime.GOOS == "windows" {
		defaultOutputDir = filepath.Join(homeDir, "Desktop")
	} else if runtime.GOOS == "darwin" {
		defaultOutputDir = filepath.Join(homeDir, "Desktop")
	}
	p.LineEditWorkDir.SetText(defaultOutputDir)
	p.LineEditWebHook.SetText("https://sctapi.ftqq.com/{KEY}.send")
	p.PushButtonWorkDir.ConnectClicked(func(bool) {
		qFile := widgets.NewQFileDialog2(nil, "选择文件目录", defaultOutputDir, "")
		workDir = qFile.GetExistingDirectory(nil, "目录", defaultOutputDir, widgets.QFileDialog__ShowDirsOnly)
	})
	p.PushButtonOK.ConnectClicked(func(bool) {
		webHook = p.LineEditWebHook.Text()
		if workDir != "" {
			defaultOutputDir = workDir
		}
		if webHook != "" {
			localWebHookInit()
		}
		widgets.QMessageBox_Information(nil, "信息", "设置完成",
			widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
	})
	p.Show()
}
