package main

import (
	"github.com/b1gcat/DarkEye/supersearch/ui"
	"github.com/therecipe/qt/widgets"
	"os"
	"runtime"
)


func runApp() {
	//加载配置
	runtime.GOMAXPROCS(runtime.NumCPU())
	app := widgets.NewQApplication(len(os.Args), os.Args)
	//app.SetWindowIcon(gui.NewQIcon5(":/qml/logo.ico"))
	app.SetQuitOnLastWindowClosed(false)
	app.SetStyle2("fusion")
	//初始化窗口
	mainWin := ui.NewMainWindow(nil)
	//初始化数据
	initMainWin(mainWin)
	mainWin.Show()

	widgets.QApplication_Exec()
}

func initMainWin(mainWin *ui.MainWindow) {
	mainWin.ActionZoomEye.ConnectTriggered(func(bool) {
		zoomEyeInit()
	})
	mainWin.ActionXray.ConnectTriggered(func(bool) {
		xrayInit()
	})
	mainWin.ActionAbout.ConnectTriggered(func(bool) {
		widgets.QMessageBox_Information(nil, "信息", "https://github.com/b1gcat/DarkEye",
			widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
	})
	mainWin.ActionPreference.ConnectTriggered(func(bool) {
		preferenceInit()
	})
}

func newLogChannel(view *widgets.QTextEdit) chan string {
	logC := make(chan string, 128)
	view.SetReadOnly(true)
	go func() {
		for {
			log := <-logC
			view.Append(log)
		}
	}()
	return logC
}