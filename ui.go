package main

import (
	"fmt"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
	"os"
	"time"
)

var (
	programName    = "DarkEye"
	programDesc    = "白嫖神器"
	programVersion = "1.0." + fmt.Sprintf("%d%d%d%d%d\nhttps://github.com/zsdevX/DarkEye\n大橘Oo0\n84500316@qq.com",
		time.Now().Year(), time.Now().Month(), time.Now().Day(), time.Now().Hour(), time.Now().Minute())
)

func main() {
	runApp()
}

func runApp() {
	app := widgets.NewQApplication(len(os.Args), os.Args)
	app.SetWindowIcon(gui.NewQIcon5(":/qml/logo.ico"))

	//加载配置
	loadCfg()

	sysTray := NewQSystemTrayIconWithCustomSlot(nil)
	app.SetQuitOnLastWindowClosed(false)

	windowFofa = registerFofa(sysTray)
	windowFofa.Show()

	windowSecurityTails = registerSecurityTrails(sysTray)
	windowSecurityTails.Hide()

	sysTrayDaemon(sysTray, app)
	sysTray.Show()
	app.Exec()
}

func sysTrayDaemon(sysTray *QSystemTrayIconWithCustomSlot, app *widgets.QApplication) {
	sysTray.SetIcon(gui.NewQIcon5(":/qml/logo.png"))
	sysTray.SetToolTip(programDesc)

	sysTrayMenu := widgets.NewQMenu(nil)
	fofa := sysTrayMenu.AddAction("资产搜索神器")
	securitytrails := sysTrayMenu.AddAction("域名搜索神器")
	about := sysTrayMenu.AddAction("关于")
	quit := sysTrayMenu.AddAction("退出")
	sysTray.SetContextMenu(sysTrayMenu)

	sysTray.ConnectTriggerSlot(func() {
		sysTray.ShowMessage("信息", programVersion, widgets.QSystemTrayIcon__Information, 5000)
	})

	fofa.ConnectTriggered(func(bool) {
		windowFofa.Show()
	})

	securitytrails.ConnectTriggered(func(bool) {
		windowSecurityTails.Show()
	})

	about.ConnectTriggered(func(bool) {
		information := programVersion
		widgets.QMessageBox_Information(nil, "信息", information,
			widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
	})

	quit.ConnectTriggered(func(bool) {
		if widgets.QMessageBox_Information(nil, "信息", "客官清走，欢迎白嫖",
			widgets.QMessageBox__Ok, widgets.QMessageBox__Cancel) == widgets.QMessageBox__Ok {
			app.Quit()
		}
	})
}

type QSystemTrayIconWithCustomSlot struct {
	widgets.QSystemTrayIcon
	_ func() `slot:"triggerSlot"`
}

type CustomEditor struct {
	widgets.QTextEdit
	_ func(string) `signal:"updateTextFromGoroutine,auto(this.QTextEdit.Append)"`
}

func getWindowCtl() (chan string, chan bool, *CustomEditor) {
	logC := make(chan string, 128)
	runCtl := make(chan bool, 1)

	logP := NewCustomEditor(nil)
	logP.SetReadOnly(true)

	go func() {
		for {
			log := <-logC
			logP.UpdateTextFromGoroutine(log)
		}
	}()
	return logC, runCtl, logP
}
