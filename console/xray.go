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
	}
	wg := sync.WaitGroup{}
	wg.Add(2)
	//启动代理
	go func() {
		defer wg.Done()
		if runtime.GOOS == "windows" {
			cmd := exec.CommandContext(ctx, "CMD", "/C", xRayProgram,
				"webscan", "--listen", "127.0.0.1:"+x.proxyPort,
				"--json-output", "vulnerability.json")
			_ = runShell(cmd)
		} else {
			cmd := exec.CommandContext(ctx, xRayProgram,
				"webscan", "--listen", "127.0.0.1:"+x.proxyPort,
				"--json-output", "vulnerability.json")
			_ = runShell(cmd)
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
				time.Sleep(time.Second * 3)
				break
			}
			time.Sleep(time.Second * 1)
		}
		for _, vulnerable := range urls {
			if runtime.GOOS == "windows" {
				cmd := exec.CommandContext(ctx, "CMD", "/C",
					crawlerGo["name"],
					"-t", "10", "-c", x.chrome,
					"--request-proxy", "http://127.0.0.1:"+x.proxyPort, vulnerable)
				_ = runShell(cmd)
			} else {
				cmd := exec.CommandContext(ctx,
					crawlerGo["name"],
					"-t", "10", "-c", x.chrome,
					"--request-proxy", "http://127.0.0.1:"+x.proxyPort, vulnerable)
				_ = runShell(cmd)
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
		x.down(x.download, bin, ctx)
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
	crawler := crawlerGo[runtime.GOOS]
	fmt.Println("Download binary from", crawler)
	x.down(x.download, bin, ctx)
	if _, e := os.Stat(bin); e != nil {
		return fmt.Errorf(fmt.Sprintf("Err: %s", e.Error()))
	}

	s := strings.Split(crawler, "/")
	zipFile := s[len(s)-1]
	if _, e := os.Stat(zipFile); e != nil {
		return fmt.Errorf(fmt.Sprintf("Err: %s", e.Error()))
	}
	uz := unzip.New()
	files, e := uz.Extract(zipFile, "./")
	if e != nil {
		return fmt.Errorf(fmt.Sprintf("Err: %s", e.Error()))
	}
	fmt.Printf("Extracted files count: %d", len(files))
	fmt.Printf("Files list: %v", files)
	zipFile = strings.TrimSuffix(zipFile, ".zip")
	if runtime.GOOS == "windows" {
		zipFile += ".exe"
	}
	_ = os.Rename(zipFile, bin)
	if _, e := os.Stat(bin); e != nil {
		return fmt.Errorf(fmt.Sprintf("Err: %s", e.Error()))
	}
	fmt.Println(bin, "[OK]")
	return nil
}

func (x *xRayRuntime) down(url, save string, ctx context.Context) {
	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	resp, _ := http.DefaultClient.Do(req)
	defer resp.Body.Close()

	f, _ := os.OpenFile(save, os.O_CREATE|os.O_WRONLY, 0700)
	defer f.Close()

	bar := progressbar.DefaultBytes(
		resp.ContentLength,
		"downloading",
	)
	io.Copy(io.MultiWriter(f, bar), resp.Body)
}
