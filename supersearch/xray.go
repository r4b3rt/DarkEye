package main

import (
	"context"
	"fmt"
	"github.com/b1gcat/DarkEye/supersearch/ui"
	"github.com/therecipe/qt/widgets"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func xrayInit() {
	ctx, cancel := context.WithCancel(context.Background())
	x := ui.NewDialogXray(nil)
	//设置默认插件
	x.LineEditPlugins.SetText(
		"baseline,brute-force,cmd-injection,crlf-injection," +
			"dirscan,jsonp,path-traversal,redirect,sqldet,ssrf," +
			"upload,xss,xxe,fastjson,shiro,struts,thinkphp,phantasm")
	//设置pocs目录
	x.PushButtonPocs.ConnectClicked(func(bool) {
		qFile := widgets.NewQFileDialog2(nil, "选择poc文件目录", defaultOutputDir, "")
		fn := qFile.GetExistingDirectory(nil, "目录", defaultOutputDir, widgets.QFileDialog__ShowDirsOnly)
		x.LineEditPocs.SetText(fn)
	})

	//自定义xray目录
	x.PushButtonXrayWorkDir.ConnectClicked(func(bool) {
		qFile := widgets.NewQFileDialog2(nil, "选择文件目录", x.LineEditWorkDir.Text(), "")
		fn := qFile.GetExistingDirectory(nil, "目录", defaultOutputDir, widgets.QFileDialog__ShowDirsOnly)
		if fn != "" {
			x.LineEditWorkDir.SetText(fn)
		} else {
			x.LineEditWorkDir.SetText(defaultOutputDir)
		}
	})

	//选择目标目录
	x.ComboBoxTargetType.ConnectActivated2(func(s string) {
		if s == "文件" {
			qFile := widgets.NewQFileDialog2(nil, "选择文件", defaultOutputDir, "")
			fn := qFile.GetOpenFileName(nil, "文件", ".", "", "", widgets.QFileDialog__ReadOnly)
			if fn == "" {
				widgets.QMessageBox_Information(nil, "信息", "未选择任何目标文件",
					widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
				x.ComboBoxOutput.SetCurrentText("自定义")
				return
			}
			x.LineEditTarget.SetText(fn)
			x.LineEditTarget.SetEnabled(false)
		} else {
			x.LineEditTarget.SetEnabled(true)
			x.LineEditTarget.Clear()
		}
	})

	x.LineEditListen.SetText("1234")

	//设置Action
	x.PushButtonStop.ConnectClicked(func(bool) {
		cancel()
	})

	x.PushButtonStart.ConnectClicked(func(bool) {
		x.PushButtonStart.SetEnabled(false)
		defer func() {
			x.PushButtonStart.SetEnabled(true)
		}()
		xRayRun(ctx, x)
	})

	x.Show()
}

func xRayRun(ctx context.Context, x *ui.DialogXray) {
	activeMode := x.TabWidget.CurrentIndex() == 0
	runXray := "xray ws"
	if x.LineEditPlugins.Text() != "" {
		runXray += " --plugins " + x.LineEditPlugins.Text()
	}
	if x.LineEditPocs.Text() != "" {
		runXray += " -p " + strings.TrimRight(x.LineEditPocs.Text(), "/") + "*.yaml"
	}
	f := filepath.Join(defaultOutputDir,
		"xray_"+
			time.Now().Format("20060102150405")+
			"."+strings.Split(x.ComboBoxOutput.CurrentText(), "-")[1])
	runXray += " --" + x.ComboBoxOutput.CurrentText() + " " + f

	if !activeMode {
		xrayRunPassive(ctx, runXray, x)
	} else {
		xrayRunActive(ctx, runXray, x)
	}

	if _, err := os.Stat(f); err == nil {
		widgets.QMessageBox_Information(nil, "记录", f,
			widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
	} else {
		widgets.QMessageBox_Information(nil, "记录", "无记录生成",
			widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
	}
}

func xrayRunActive(ctx context.Context, runXray string, x *ui.DialogXray) {
	defer func() {
		widgets.QMessageBox_Information(nil, "信息", "结束xrayRunActive",
			widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
	}()
	if x.LineEditExtra.Text() != "" {
		runXray += " " + x.LineEditExtra.Text()
	}

	if x.ComboBoxTargetType.CurrentText() == "单一URL" {
		if x.LineEditTarget.Text() == "" {
			widgets.QMessageBox_Information(nil, "记录", "未设置目标",
				widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
			return
		}
		runXray += " --url " + x.LineEditTarget.Text()
	} else {
		runXray += " --url-file " + x.LineEditTarget.Text()
	}
	goXray(ctx, runXray, x)
}

func xrayRunPassive(ctx context.Context, runXray string, x *ui.DialogXray) {
	defer func() {
		widgets.QMessageBox_Information(nil, "信息", "结束xrayRunPassive",
			widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
	}()
	runXray += " --listen 0.0.0.0:" + x.LineEditListen.Text()
	select {
	case <-ctx.Done():
		return
	default:
		//fixme: macbook 下过一段时间被动模式会退出
		goXray(ctx, runXray, x)
	}
}

func goXray(ctx context.Context, runXray string, x *ui.DialogXray) {
	var c *exec.Cmd
	ctx2, _ := context.WithCancel(ctx)
	if runtime.GOOS == "windows" {
		c = exec.CommandContext(ctx2, "CMD", "/C", runXray)
	} else if runtime.GOOS == "darwin" {
		c = exec.Command(`osascript`, "-s", "h", "-e",
			fmt.Sprintf(`tell application "Terminal" to do script "cd %s && %s"`,
				x.LineEditWorkDir.Text(), "./"+runXray))
	} else {
		widgets.QMessageBox_Information(nil, "!", "暂时不支持Linux",
			widgets.QMessageBox__Ok, widgets.QMessageBox__Ok)
		return
	}
	c.Dir = x.LineEditWorkDir.Text()

	if err := c.Run(); err != nil {
		return
	}
}
