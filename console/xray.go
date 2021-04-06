package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/artdarek/go-unzip/pkg/unzip"
	"github.com/schollz/progressbar"
	"github.com/zsdevX/DarkEye/common"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

type xRayRuntime struct {
	download  string
	url       string
	proxyPort string
	chrome    string
	flagSet   *flag.FlagSet
}

var (
	xRayRuntimeOptions = &xRayRuntime{
		flagSet: flag.NewFlagSet(xRayProgram, flag.ExitOnError),
	}
	crawlerGo = map[string]string{
		"name":    "crawlerGo",
		"darwin":  "https://ghproxy.com/https://github.com/0Kee-Team/crawlergo/releases/download/v0.4.0/crawlergo_darwin_amd64.zip",
		"linux":   "https://ghproxy.com/https://github.com/0Kee-Team/crawlergo/releases/download/v0.4.0/crawlergo_linux_amd64.zip",
		"windows": "https://ghproxy.com/https://github.com/0Kee-Team/crawlergo/releases/download/v0.4.0/crawlergo_windows_amd64.zip",
	}
)

func xRayInitRunTime() {
	xRayRuntimeOptions.flagSet.StringVar(&xRayRuntimeOptions.download,
		"download", "https://ghproxy.com/https://github.com/zsdevX/helper/releases/download/1/xray_darwin_amd64", "xRay binary file")
	xRayRuntimeOptions.flagSet.StringVar(&xRayRuntimeOptions.url,
		"url", "http://vuln.com.cn", "vulnerable website or file list")
	xRayRuntimeOptions.flagSet.StringVar(&xRayRuntimeOptions.proxyPort,
		"proxy-port", "7777", "被动扫描代理端口")
	xRayRuntimeOptions.flagSet.StringVar(&xRayRuntimeOptions.proxyPort,
		"chrome", "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome", "Chrome path")

}

func (x *xRayRuntime) compileArgs(cmd []string) error {
	if err := x.flagSet.Parse(splitCmd(cmd)); err != nil {
		return err
	}
	x.flagSet.Parsed()
	return nil
}

func (x *xRayRuntime) usage() {
	fmt.Println(fmt.Sprintf("Usage of %s:", xRayProgram))
	fmt.Println("Options:")
	x.flagSet.VisitAll(func(f *flag.Flag) {
		var opt = "  -" + f.Name
		fmt.Println(opt)
		fmt.Println(fmt.Sprintf("		%v (default '%v')", f.Usage, f.DefValue))
	})
}

func (x *xRayRuntime) start(ctx context.Context) {
	if err := x.prepare(ctx); err != nil {
		fmt.Println(err)
		return
	}
	wg := sync.WaitGroup{}
	wg.Add(2)
	//启动代理
	go func() {
		defer wg.Done()
		var cmd *exec.Cmd
		if runtime.GOOS == "windows" {
			cmd = exec.CommandContext(ctx, "CMD", "/C", xRayProgram,
				"webscan", "--listen", "127.0.0.1:"+x.proxyPort,
				"--json-output", "vulnerability.json")
		} else {
			cmd = exec.CommandContext(ctx, xRayProgram,
				"webscan", "--listen", "127.0.0.1:"+x.proxyPort,
				"--json-output", "vulnerability.json")
		}
		if err := runShell(cmd); err != nil {
			fmt.Println("Err:", err.Error())
		}
	}()

	go func() {
		defer wg.Done()
		urls := make([]string, 0)
		if _, e := os.Stat(x.url); e != nil {
			urls = common.GenDicFromFile(x.url)
		} else {
			urls = append(urls, x.url)
		}
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}
			if common.IsAlive("127.0.0.1", x.proxyPort, 1000) == common.Alive {
				fmt.Println("Warn", "未检测到代理等待3秒重试")
				time.Sleep(time.Second * 3)
				break
			}
			time.Sleep(time.Second * 1)
		}
		for _, vulnerable := range urls {
			var cmd *exec.Cmd
			if runtime.GOOS == "windows" {
				cmd = exec.CommandContext(ctx, "CMD", "/C",
					crawlerGo["name"],
					"-t", "10", "-c", x.chrome,
					"--request-proxy", "http://127.0.0.1:"+x.proxyPort, vulnerable)
			} else {
				cmd = exec.CommandContext(ctx,
					crawlerGo["name"],
					"-t", "10", "-c", x.chrome,
					"--request-proxy", "http://127.0.0.1:"+x.proxyPort, vulnerable)
			}
			if err := runShell(cmd); err != nil {
				fmt.Println("Err:", err.Error())
			}
		}
	}()
	wg.Wait()
}

func (x *xRayRuntime) prepare(ctx context.Context) error {
	//检查授权文件是否存在（狗头）
	if _, e := os.Stat("xray-license.lic"); e != nil {
		return fmt.Errorf(fmt.Sprintf("Err: %s", e.Error()))
	}
	//检查xRay是否存在
	bin := xRayProgram
	if runtime.GOOS == "windows" {
		bin += ".exe"
	}
	if _, e := os.Stat(bin); e != nil {
		fmt.Println("Download binary from", x.download)
		if e = x.down(x.download, bin, ctx); e != nil {
			return e
		}
		if _, e := os.Stat(bin); e != nil {
			return fmt.Errorf(fmt.Sprintf("Err: %s", e.Error()))
		}
	}
	fmt.Println(bin, "[OK]")
	//下载爬虫
	bin = crawlerGo["name"]
	if runtime.GOOS == "windows" {
		bin += ".exe"
	}
	if _, e := os.Stat(bin); e == nil {
		fmt.Println(bin, "[OK]")
		return nil
	}
	//未发现则下载
	crawler := crawlerGo[runtime.GOOS]
	s := strings.Split(crawler, "/")
	zipFile := s[len(s)-1]
	fmt.Println("Download binary from", crawler)
	if e := x.down(crawler, zipFile, ctx); e != nil {
		return e
	}
	if _, e := os.Stat(zipFile); e != nil {
		return fmt.Errorf(fmt.Sprintf("Err: %s", e.Error()))
	}
	uz := unzip.New()
	destDir := strings.TrimSuffix(zipFile, ".zip")
	files, e := uz.Extract(zipFile, destDir)
	if e != nil {
		return fmt.Errorf(fmt.Sprintf("Err: %s", e.Error()))
	}
	fmt.Println(fmt.Sprintf("Extracted files count: %d", len(files)))
	fmt.Println(fmt.Sprintf("Files list: %v", files))
	targetFile := "crawlergo"
	if runtime.GOOS == "windows" {
		targetFile += ".exe"
	}
	_ = os.Rename(filepath.Join(destDir, targetFile), bin)
	if _, e := os.Stat(bin); e != nil {
		return fmt.Errorf(fmt.Sprintf("Err: %s", e.Error()))
	}
	fmt.Println(bin, "[OK]")
	return nil
}

func (x *xRayRuntime) down(url, save string, ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	f, err := os.OpenFile(save, os.O_CREATE|os.O_WRONLY, 0700)
	if err != nil {
		return err
	}
	defer f.Close()

	bar := progressbar.NewOptions64(
		resp.ContentLength,
		progressbar.OptionSetDescription("Downloading"),
		progressbar.OptionSetWriter(os.Stderr),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(10),
		progressbar.OptionThrottle(65*time.Millisecond),
		progressbar.OptionShowCount(),
		progressbar.OptionOnCompletion(func() {
			fmt.Fprint(os.Stderr, "\n")
		}),
		progressbar.OptionSpinnerType(14),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "=",
			SaucerHead:    ">",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
		progressbar.OptionFullWidth(),
	)
	_ = bar.RenderBlank()

	_, _ = io.Copy(io.MultiWriter(f, bar), resp.Body)
	return nil
}
