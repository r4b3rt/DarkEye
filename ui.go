package main

import (
	"fmt"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
	"github.com/zsdevX/DarkEye/common"
	"github.com/zsdevX/DarkEye/ui"
	"os"
	"runtime"
	"time"
)

var (
	programName = "DarkEye"
)

func main() {
	runApp()
}

func runApp() {
	//加载配置
	_ = loadCfg()
	runtime.GOMAXPROCS(runtime.NumCPU())
	app := widgets.NewQApplication(len(os.Args), os.Args)
	app.SetWindowIcon(gui.NewQIcon5(":/qml/logo.ico"))
	app.SetQuitOnLastWindowClosed(false)
	//初始化窗口
	mainWin := ui.NewMainWindow(nil)
	//初始化数据
	initMainWin(mainWin)
	mainWin.SetStyleSheet(QSS)
	//托盘图标初始化
	sysTray := NewQSystemTrayIconWithCustomSlot(nil)
	sysTrayDaemon(sysTray, mainWin, app)
	//显示
	sysTray.Show()
	mainWin.Show()
	//通知
	sysTray.TriggerSlot()
	widgets.QApplication_Exec()
}

func initMainWin(mainWin *ui.MainWindow) {
	//FoFa
	LoadAsset(mainWin)
	//subDomain
	LoadSubDomain(mainWin)
	//Spider
	LoadSpider(mainWin)
}

func sysTrayDaemon(sysTray *QSystemTrayIconWithCustomSlot, mainWin *ui.MainWindow, app *widgets.QApplication) {
	sysTray.SetIcon(gui.NewQIcon5(":/qml/logo.png"))
	sysTray.SetToolTip("白嫖神器")

	sysTrayMenu := widgets.NewQMenu(nil)
	fucker := sysTrayMenu.AddAction("信息收集神器")
	about := sysTrayMenu.AddAction("关于")
	quit := sysTrayMenu.AddAction("退出")
	sysTray.SetContextMenu(sysTrayMenu)

	sysTray.ConnectTriggerSlot(func() {
		sysTray.ShowMessage("信息", common.ProgramVersion, widgets.QSystemTrayIcon__Information, 5000)
	})

	fucker.ConnectTriggered(func(bool) {
		mainWin.Show()
	})

	about.ConnectTriggered(func(bool) {
		information := common.ProgramVersion
		widgets.QMessageBox_Information(nil, "信息", information,
			widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
	})

	quit.ConnectTriggered(func(bool) {
		if widgets.QMessageBox_Information(nil, "信息", "客官再见，欢迎白嫖",
			widgets.QMessageBox__Ok, widgets.QMessageBox__Cancel) == widgets.QMessageBox__Ok {
			app.Quit()
		}
	})
}

//外部应使用goroutine调用
func gracefulStop(start, stop *widgets.QPushButton, runCtl chan bool) {
	//终止任务时避免卡顿
	sec := 0
	jumpOut := false
	tick := time.NewTicker(time.Second)
	stop.SetDisabled(true)
	for {
		select {
		case <-runCtl:
			jumpOut = true
		case <-tick.C:
			sec++
			stop.SetText(fmt.Sprintf("等待%d秒", 60-sec))
		}
		if jumpOut {
			break
		}
	}
	start.SetEnabled(true)
	stop.SetText("停止")
}

//QSystemTrayIconWithCustomSlot add comment
type QSystemTrayIconWithCustomSlot struct {
	widgets.QSystemTrayIcon
	_ func() `slot:"triggerSlot"`
}

func logChannel(view *widgets.QTextEdit) (chan string, chan bool) {
	logC := make(chan string, 128)
	runCtl := make(chan bool, 1)
	view.SetReadOnly(true)
	go func() {
		for {
			log := <-logC
			view.Append(log)
		}
	}()
	return logC, runCtl
}

var (
	//QSS add comment
	QSS = `
QMainWindow {
	background-color:#1e1d23;
}
QDialog {
	background-color:#1e1d23;
}
QColorDialog {
	background-color:#1e1d23;
}
QTextEdit {
	background-color:#1e1d23;
	color: #a9b7c6;
}
QPlainTextEdit {
	selection-background-color:#007b50;
	background-color:#1e1d23;
	border-style: solid;
	border-top-color: transparent;
	border-right-color: transparent;
	border-left-color: transparent;
	border-bottom-color: transparent;
	border-width: 1px;
	color: #a9b7c6;
}
QPushButton{
	border-style: solid;
	border-top-color: transparent;
	border-right-color: transparent;
	border-left-color: transparent;
	border-bottom-color: transparent;
	border-width: 1px;
	border-style: solid;
	color: #a9b7c6;
	padding: 2px;
	background-color: #1e1d23;
}
QPushButton::default{
	border-style: inset;
	border-top-color: transparent;
	border-right-color: transparent;
	border-left-color: transparent;
	border-bottom-color: #04b97f;
	border-width: 1px;
	color: #a9b7c6;
	padding: 2px;
	background-color: #1e1d23;
}
QToolButton {
	border-style: solid;
	border-top-color: transparent;
	border-right-color: transparent;
	border-left-color: transparent;
	border-bottom-color: #04b97f;
	border-bottom-width: 1px;
	border-style: solid;
	color: #a9b7c6;
	padding: 2px;
	background-color: #1e1d23;
}
QToolButton:hover{
	border-style: solid;
	border-top-color: transparent;
	border-right-color: transparent;
	border-left-color: transparent;
	border-bottom-color: #37efba;
	border-bottom-width: 2px;
	border-style: solid;
	color: #FFFFFF;
	padding-bottom: 1px;
	background-color: #1e1d23;
}
QPushButton:hover{
	border-style: solid;
	border-top-color: transparent;
	border-right-color: transparent;
	border-left-color: transparent;
	border-bottom-color: #37efba;
	border-bottom-width: 1px;
	border-style: solid;
	color: #FFFFFF;
	padding-bottom: 2px;
	background-color: #1e1d23;
}
QPushButton:pressed{
	border-style: solid;
	border-top-color: transparent;
	border-right-color: transparent;
	border-left-color: transparent;
	border-bottom-color: #37efba;
	border-bottom-width: 2px;
	border-style: solid;
	color: #37efba;
	padding-bottom: 1px;
	background-color: #1e1d23;
}
QPushButton:disabled{
	border-style: solid;
	border-top-color: transparent;
	border-right-color: transparent;
	border-left-color: transparent;
	border-bottom-color: #808086;
	border-bottom-width: 2px;
	border-style: solid;
	color: #808086;
	padding-bottom: 1px;
	background-color: #1e1d23;
}
QLineEdit {
	border-width: 1px; border-radius: 4px;
	border-color: rgb(58, 58, 58);
	border-style: inset;
	padding: 0 8px;
	color: #a9b7c6;
	background:#1e1d23;
	selection-background-color:#007b50;
	selection-color: #FFFFFF;
}
QLabel {
	color: #a9b7c6;
}
QLCDNumber {
	color: #37e6b4;
}
QProgressBar {
	text-align: center;
	color: rgb(240, 240, 240);
	border-width: 1px; 
	border-radius: 10px;
	border-color: rgb(58, 58, 58);
	border-style: inset;
	background-color:#1e1d23;
}
QProgressBar::chunk {
	background-color: #04b97f;
	border-radius: 5px;
}
QMenuBar {
	background-color: #1e1d23;
}
QMenuBar::item {
	color: #a9b7c6;
  	spacing: 3px;
  	padding: 1px 4px;
  	background: #1e1d23;
}

QMenuBar::item:selected {
  	background:#1e1d23;
	color: #FFFFFF;
}
QMenu::item:selected {
	border-style: solid;
	border-top-color: transparent;
	border-right-color: transparent;
	border-left-color: #04b97f;
	border-bottom-color: transparent;
	border-left-width: 2px;
	color: #FFFFFF;
	padding-left:15px;
	padding-top:4px;
	padding-bottom:4px;
	padding-right:7px;
	background-color: #1e1d23;
}
QMenu::item {
	border-style: solid;
	border-top-color: transparent;
	border-right-color: transparent;
	border-left-color: transparent;
	border-bottom-color: transparent;
	border-bottom-width: 1px;
	border-style: solid;
	color: #a9b7c6;
	padding-left:17px;
	padding-top:4px;
	padding-bottom:4px;
	padding-right:7px;
	background-color: #1e1d23;
}
QMenu{
	background-color:#1e1d23;
}
QTabWidget {
	color:rgb(0,0,0);
	background-color:#1e1d23;
}
QTabWidget::pane {
		border-color: rgb(77,77,77);
		background-color:#1e1d23;
		border-style: solid;
		border-width: 1px;
    	border-radius: 6px;
}
QTabBar::tab {
	border-style: solid;
	border-top-color: transparent;
	border-right-color: transparent;
	border-left-color: transparent;
	border-bottom-color: transparent;
	border-bottom-width: 1px;
	border-style: solid;
	color: #808086;
	padding: 3px;
	margin-left:3px;
	background-color: #1e1d23;
}
QTabBar::tab:selected, QTabBar::tab:last:selected, QTabBar::tab:hover {
  	border-style: solid;
	border-top-color: transparent;
	border-right-color: transparent;
	border-left-color: transparent;
	border-bottom-color: #04b97f;
	border-bottom-width: 2px;
	border-style: solid;
	color: #FFFFFF;
	padding-left: 3px;
	padding-bottom: 2px;
	margin-left:3px;
	background-color: #1e1d23;
}

QCheckBox {
	color: #a9b7c6;
	padding: 2px;
}
QCheckBox:disabled {
	color: #808086;
	padding: 2px;
}

QCheckBox:hover {
	border-radius:4px;
	border-style:solid;
	padding-left: 1px;
	padding-right: 1px;
	padding-bottom: 1px;
	padding-top: 1px;
	border-width:1px;
	border-color: rgb(87, 97, 106);
	background-color:#1e1d23;
}
QCheckBox::indicator:checked {

	height: 10px;
	width: 10px;
	border-style:solid;
	border-width: 1px;
	border-color: #04b97f;
	color: #a9b7c6;
	background-color: #04b97f;
}
QCheckBox::indicator:unchecked {

	height: 10px;
	width: 10px;
	border-style:solid;
	border-width: 1px;
	border-color: #04b97f;
	color: #a9b7c6;
	background-color: transparent;
}
QRadioButton {
	color: #a9b7c6;
	background-color: #1e1d23;
	padding: 1px;
}
QRadioButton::indicator:checked {
	height: 10px;
	width: 10px;
	border-style:solid;
	border-radius:5px;
	border-width: 1px;
	border-color: #04b97f;
	color: #a9b7c6;
	background-color: #04b97f;
}
QRadioButton::indicator:!checked {
	height: 10px;
	width: 10px;
	border-style:solid;
	border-radius:5px;
	border-width: 1px;
	border-color: #04b97f;
	color: #a9b7c6;
	background-color: transparent;
}
QStatusBar {
	color:#027f7f;
}
QSpinBox {
	color: #a9b7c6;	
	background-color: #1e1d23;
}
QDoubleSpinBox {
	color: #a9b7c6;	
	background-color: #1e1d23;
}
QTimeEdit {
	color: #a9b7c6;	
	background-color: #1e1d23;
}
QDateTimeEdit {
	color: #a9b7c6;	
	background-color: #1e1d23;
}
QDateEdit {
	color: #a9b7c6;	
	background-color: #1e1d23;
}
QComboBox {
	color: #a9b7c6;	
	background: #1e1d23;
}
QComboBox:editable {
	background: #1e1d23;
	color: #a9b7c6;
	selection-background-color: #1e1d23;
}
QComboBox QAbstractItemView {
	color: #a9b7c6;	
	background: #1e1d23;
	selection-color: #FFFFFF;
	selection-background-color: #1e1d23;
}
QComboBox:!editable:on, QComboBox::drop-down:editable:on {
	color: #a9b7c6;	
	background: #1e1d23;
}
QFontComboBox {
	color: #a9b7c6;	
	background-color: #1e1d23;
}
QToolBox {
	color: #a9b7c6;
	background-color: #1e1d23;
}
QToolBox::tab {
	color: #a9b7c6;
	background-color: #1e1d23;
}
QToolBox::tab:selected {
	color: #FFFFFF;
	background-color: #1e1d23;
}
QScrollArea {
	color: #FFFFFF;
	background-color: #1e1d23;
}
QSlider::groove:horizontal {
	height: 5px;
	background: #04b97f;
}
QSlider::groove:vertical {
	width: 5px;
	background: #04b97f;
}
QSlider::handle:horizontal {
	background: qlineargradient(x1:0, y1:0, x2:1, y2:1, stop:0 #b4b4b4, stop:1 #8f8f8f);
	border: 1px solid #5c5c5c;
	width: 14px;
	margin: -5px 0;
	border-radius: 7px;
}
QSlider::handle:vertical {
	background: qlineargradient(x1:1, y1:1, x2:0, y2:0, stop:0 #b4b4b4, stop:1 #8f8f8f);
	border: 1px solid #5c5c5c;
	height: 14px;
	margin: 0 -5px;
	border-radius: 7px;
}
QSlider::add-page:horizontal {
    background: white;
}
QSlider::add-page:vertical {
    background: white;
}
QSlider::sub-page:horizontal {
    background: #04b97f;
}
QSlider::sub-page:vertical {
    background: #04b97f;
}
`
)
