package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/artdarek/go-unzip/pkg/unzip"
	"github.com/schollz/progressbar"
	"github.com/zsdevX/DarkEye/common"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type xRayRuntime struct {
	Module
	parent *RequestContext

	download  string
	url       string
	proxyPort string
	chrome    string
	flagSet   *flag.FlagSet
	cmd       []string
}

type xRayTarget struct {
	Url    string `json:"url"`
	Method string `json:"method"`
	Data   string `json:"data"`
}

type xRayReq struct {
	Req []xRayTarget `json:"req_list"`
}

var (
	xRayProgram        = "xRay"
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

func (x *xRayRuntime) Start(parent context.Context) {
	//准备工具
	preCtx, _ := context.WithCancel(parent)
	if err := x.prepare(preCtx); err != nil {
		common.Log("xRayRuntime.prepare", err.Error(), common.INFO)
	}
	x.parent.taskId ++
	//爬取数据
	cCtx, _ := context.WithCancel(parent)
	if err := x.crawler(cCtx); err != nil {
		common.Log("xRayRuntime.crawler", err.Error(), common.INFO)
		return
	}
	//模拟请求
	go func() {
		loop := 0
		for loop < 10 {
			if common.IsAlive(parent, "127.0.0.1", x.proxyPort, 1000) == common.Alive {
				time.Sleep(time.Second * 5)
				break
			}
			common.Log("xRayRuntime.proxyServer.checkProxy", "等待proxy, 1秒后尝试", common.INFO)
			time.Sleep(time.Second * 1)
			loop ++
		}
		sCtx, _ := context.WithCancel(parent)
		if err := x.simulate(sCtx); err != nil {
			common.Log("xRayRuntime.simulate", err.Error(), common.INFO)
		}
	}()
	//开始访问
	pCtx, _ := context.WithCancel(parent)
	x.proxyServer(pCtx)
}

func (x *xRayRuntime) Init(requestContext *RequestContext) {
	x.parent = requestContext
	x.flagSet.StringVar(&x.download,
		"download", "https://ghproxy.com/https://github.com/zsdevX/helper/releases/download/1/xray_darwin_amd64", "xRay binary file")
	x.flagSet.StringVar(&x.url,
		"url", "http://vuln.com.cn", "vulnerable website or file list")
	x.flagSet.StringVar(&x.proxyPort,
		"proxy-port", "7777", "被动扫描代理端口")
	x.flagSet.StringVar(&x.chrome,
		"chrome", "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome", "Chrome path")
}

func (x *xRayRuntime) ValueCheck(value string) (bool, error) {
	if v, ok := xRayValueCheck[value]; ok {
		if isDuplicateArg(value, x.parent.CmdArgs) {
			return false, fmt.Errorf("参数重复")
		}
		return v, nil
	}
	return false, fmt.Errorf("无此参数")
}

func (x *xRayRuntime) CompileArgs(cmd []string, os []string) error {
	if cmd != nil {
		if err := x.flagSet.Parse(splitCmd(cmd)); err != nil {
			return err
		}
		x.flagSet.Parsed()
	} else {
		if err := x.flagSet.Parse(os); err != nil {
			return err
		}
	}
	return nil
}

func (x *xRayRuntime) Usage() {
	fmt.Println(fmt.Sprintf("Usage of %s:", xRayProgram))
	fmt.Println("Options:")
	x.flagSet.VisitAll(func(f *flag.Flag) {
		var opt = "  -" + f.Name
		fmt.Println(opt)
		fmt.Println(fmt.Sprintf("		%v (default '%v')", f.Usage, f.DefValue))
	})
}

func (x *xRayRuntime) proxyServer(ctx context.Context) {
	defer func() {
		common.Log("xRayRuntime.proxyServer", "结束", common.ALERT)
	}()
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.CommandContext(ctx, "CMD", "/C", xRayProgram,
			"webscan", "--listen", "127.0.0.1:"+x.proxyPort,
			"--json-output", "vulnerability.json")
	} else {
		cmd = exec.CommandContext(ctx, "./"+xRayProgram,
			"webscan", "--listen", "127.0.0.1:"+x.proxyPort,
			"--json-output", "vulnerability.json")
	}
	if err := runShell(xRayProgram, cmd, false); err != nil {
		common.Log("xRayRuntime.proxyServer", err.Error(), common.FAULT)
	}
}

func (x *xRayRuntime) simulate(ctx context.Context) error {
	cr, err := analysisRuntimeOptions.getCrawler()
	if err != nil {
		return err
	}
	for _, t := range cr {
		proxy := "http://127.0.0.1:" + x.proxyPort
		if strings.HasPrefix(t.Url, "https://") {
			proxy = "https://127.0.0.1:" + x.proxyPort
		}
		req := common.HttpRequest{
			Method: t.Method,
			Url:    t.Url,
			Body:   []byte(t.Data),
			Proxy:  proxy,
			Headers: map[string]string{
				"User-Agent": common.UserAgents[0],
			},
		}
		if _, e := req.Go(); e != nil {
			common.Log("xRayRuntime.proxyServer", e.Error(), common.ALERT)
		}

	}
	return nil
}

func (x *xRayRuntime) crawler(ctx context.Context) error {
	defer func() {
		common.Log("xRayRuntime.crawler", "结束", common.ALERT)
	}()
	var urls []string
	var err error
	urls, err = analysisRuntimeOptions.Var("", "$URL")
	if err != nil {
		if _, e := os.Stat(x.url); e != nil {
			urls = append(urls, x.url)
		} else {
			urls = common.GenDicFromFile(x.url)
		}
	}
	if len(urls) == 0 {
		return fmt.Errorf("目标空")
	}
	//先清空上一次结果
	analysisRuntimeOptions.cleanCrawler()
	for _, vulnerable := range urls {
		_ = os.Remove(runShellOutput)
		var cmd *exec.Cmd
		if runtime.GOOS == "windows" {
			cmd = exec.CommandContext(ctx, "CMD", "/C",
				crawlerGo["name"],
				"-t", "10", "-f", "smart", "--fuzz-path", "-c", x.chrome,
				"--output-mode", "json", "--wait-dom-content-loaded-timeout", "10s", vulnerable)
		} else {
			cmd = exec.CommandContext(ctx,
				"./"+crawlerGo["name"],
				"-t", "10", "-f", "smart", "--fuzz-path", "-c", x.chrome,
				"--output-mode", "json", "--wait-dom-content-loaded-timeout", "10s", vulnerable)
		}
		if err := runShell(crawlerGo["name"], cmd, true); err != nil {
			return err
		}
		if err := x.saveCrawler(vulnerable); err != nil {
			return err
		}
	}
	return nil
}

func (x *xRayRuntime) saveCrawler(target string) error {
	out, err := ioutil.ReadFile(runShellOutput)
	if err != nil {
		return err
	}
	tmpJson := bytes.Split(out, []byte("--[Mission Complete]--"))
	if len(tmpJson) != 2 {
		return fmt.Errorf("无 crawler result")
	}
	urlJson := tmpJson[1]
	targets := xRayReq{}
	if err = json.Unmarshal(urlJson, &targets); err != nil {
		return err
	}
	for _, t := range targets.Req {
		analysisRuntimeOptions.upInsertCrawler(&crawler{
			Target: target,
			Url:    t.Url,
			Method: t.Method,
			Data:   t.Data,
		})
	}
	return nil
}

func (x *xRayRuntime) prepare(ctx context.Context) error {
	//检查授权文件是否存在（狗头）
	if _, e := os.Stat("xray-license.lic"); e != nil {
		return fmt.Errorf("将xray的license文件(xray-license.lic)放到这里")
	}
	//检查xRay是否存在
	bin := xRayProgram
	if runtime.GOOS == "windows" {
		bin += ".exe"
	}
	if _, e := os.Stat(bin); e != nil {
		common.Log("xRayRuntime.prepare", "Download binary from "+x.download, common.INFO)
		if e = x.down(x.download, bin, ctx); e != nil {
			return e
		}
		if _, e := os.Stat(bin); e != nil {
			return fmt.Errorf(fmt.Sprintf("Err: %s", e.Error()))
		}
	}
	common.Log("xRayRuntime.prepare."+bin, "ok", common.INFO)
	//下载爬虫
	bin = crawlerGo["name"]
	if runtime.GOOS == "windows" {
		bin += ".exe"
	}
	if _, e := os.Stat(bin); e == nil {
		common.Log("xRayRuntime.prepare."+bin, "ok", common.INFO)
		return nil
	}
	//未发现则下载
	crawler := crawlerGo[runtime.GOOS]
	s := strings.Split(crawler, "/")
	zipFile := s[len(s)-1]
	common.Log("xRayRuntime.prepare", "Download binary from "+crawler, common.INFO)
	if e := x.down(crawler, zipFile, ctx); e != nil {
		return e
	}
	if _, e := os.Stat(zipFile); e != nil {
		return e
	}
	uz := unzip.New()
	destDir := strings.TrimSuffix(zipFile, ".zip")
	_, e := uz.Extract(zipFile, destDir)
	if e != nil {
		return e
	}
	targetFile := "crawlergo"
	if runtime.GOOS == "windows" {
		targetFile += ".exe"
	}
	_ = os.Rename(filepath.Join(destDir, targetFile), bin)
	if _, e := os.Stat(bin); e != nil {
		return e
	}
	common.Log("xRayRuntime.prepare."+bin, "OK", common.INFO)
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
