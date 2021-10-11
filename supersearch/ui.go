package main

import (
	"bytes"
	"github.com/b1gcat/DarkEye/supersearch/ui"
	"github.com/noborus/trdsql"
	"github.com/therecipe/qt/widgets"
	"os"
	"runtime"
)

var (
	programName = "DarkEye"
)

func main() {
	runApp()
}

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
	mainWin.PlainTextEditMain.SetPlainText(``)

	mainWin.ActionZoomEye.ConnectTriggered(func(bool) {
		zoomEyeInit()
	})

	mainWin.ActionREADME.ConnectTriggered(func(bool) {
		widgets.QMessageBox_Information(nil, "信息", "https://github.com/b1gcat/DarkEye",
			widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
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

func output(d []byte, f string) error {
	r := bytes.NewBuffer(d)
	importer, err := trdsql.NewBufferImporter("any", r, trdsql.InFormat(trdsql.JSON))
	if err != nil {
		return err
	}
	fp, err := os.Create(f)
	if err != nil {
		return err
	}
	defer fp.Close()

	writer := trdsql.NewWriter(trdsql.OutFormat(trdsql.CSV),
		trdsql.OutHeader(true),
		trdsql.OutStream(fp))
	trd := trdsql.NewTRDSQL(importer, trdsql.NewExporter(writer))
	err = trd.Exec("select * from any")
	if err != nil {
		return err
	}
	return nil
}
