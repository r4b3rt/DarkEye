package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"github.com/fatih/color"
	"github.com/k0kubun/go-ansi"
	"github.com/schollz/progressbar"
	"github.com/zsdevX/DarkEye/common"
	"github.com/zsdevX/DarkEye/superscan/plugins"
	"golang.org/x/time/rate"
	"time"

	//	_ "net/http/pprof"
	"os"
	"runtime"
	"strings"
	"sync"
)

var (
	mIp           = flag.String("ip", "127.0.0.1", "a.b.c.1-254")
	mTimeOut      = flag.Int("timeout", 3000, "单位ms")
	mThread       = flag.Int("thread", 32, "扫单IP线程数")
	mPortList     = flag.String("port-list", common.PortList, "端口范围,默认1000+常用端口")
	mUserList     = flag.String("user-file", "", "用户名字典文件")
	mPassList     = flag.String("pass-file", "", "密码字典文件")
	mNoTrust      = flag.Bool("no-trust", false, "由端口判定协议改为指纹方式判断协议,速度慢点")
	mPluginWorker = flag.Int("plugin-worker", 2, "单协议爆破密码时，线程个数")
	mRateLimiter  = flag.Int("pps", 0, "扫描工具整体发包频率n/秒, 该选项可避免线程过多发包会占有带宽导致丢失目标端口")
	mActivePort   = flag.String("alive_port", "0", "使用已知开放的端口校正扫描行为。例如某服务器限制了IP访问频率，开启此功能后程序发现限制会自动调整保证扫描完整、准确")
	mListPlugin  = flag.Bool("list-plugin", false, "列出支持的爆破协议")
	mMaxIPDetect  = 16
	mFile         *os.File
	mCsvWriter    *csv.Writer
	mFileName     string
	mBar          *progressbar.ProgressBar
	mPps          *rate.Limiter
)

var (
	mScans = make([]*Scan, 0)
	//记录文件
	mFileSync = sync.RWMutex{}
)

func recordInit() {
	var err error
	_, mFile, mFileName, err = common.CreateCSV("superScan",
		[]string{"IP", "端口", "协议", "插件信息"})
	if err != nil {
		panic("创建记录文件失败" + err.Error())
	}
	mCsvWriter = csv.NewWriter(mFile)
	fmt.Println("记录结果将保存在", mFileName)
}

func main() {
	color.Yellow("超级弱口令、系统Vulnerable检测\n")
	flag.Parse()
	if *mListPlugin {
		plugins.SupportPlugin()
		return
	}
	plugins.SetDicByFile(*mUserList, *mPassList)
	runtime.GOMAXPROCS(runtime.NumCPU())
	if *mRateLimiter > 0 {
		//每秒发包*mRateLimiter，缓冲10个
		mPps = rate.NewLimiter(rate.Every(1000000*time.Microsecond/time.Duration(*mRateLimiter)), 10)
		color.Green("rate limit enable <= %v pps\n", mPps.Limit())
	}
	//  debug/pprof
	/*
	go func() {
		fmt.Println(http.ListenAndServe("localhost:10000", nil))
	}()
	*/
	color.Red(common.Banner)
	recordInit()
	Start()
}

func Start() {
	//初始化scan对象
	ips := strings.Split(*mIp, ",")
	tot := 0
	for _, ip := range ips {
		base, start, end, err := common.GetIPRange(ip)
		if err != nil {
			fmt.Println("IP格式错误:", err.Error())
			return
		}
		for {
			if start > end {
				break
			}
			s := NewScan(common.GenIP(base, start))
			s.ActivePort = "0"
			_, t := common.GetPortRange(s.PortRange)
			tot += t
			mScans = append(mScans, s)
			start++
		}
	}
	setRLimit()
	color.Green(fmt.Sprintf("已加载%d个IP,共计%d个端口,启动检测线程数%d,同时可检测IP数量%d",
		len(mScans), tot, *mThread, mMaxIPDetect))
	//建立进度条
	mBar = NewBar(tot)
	if len(mScans) == 1 {
		//单IP支持校正
		mScans[0].ActivePort = *mActivePort
	}
	limiter := make(chan int, mMaxIPDetect)
	wg := sync.WaitGroup{}
	wg.Add(len(mScans))
	for _, s := range mScans {
		go func(s *Scan) {
			limiter <- 1
			s.Run()
			<-limiter
			wg.Done()
		}(s)
	}
	wg.Wait()
	color.Red("Done")
}

func NewScan(ip string) *Scan {
	return &Scan{
		Ip:                     ip,
		TimeOut:                *mTimeOut,
		ActivePort:             *mActivePort,
		PortRange:              *mPortList,
		ThreadNumber:           *mThread,
		NoTrust:                *mNoTrust,
		Rate:                   mPps,
		PluginWorker:           *mPluginWorker,
		Callback:               myCallback,
		BarCallback:            myBarCallback,
		BarDescriptionCallback: myBarDescUpdate,
	}
}

func myBarDescUpdate(a string) {
	mBar.Describe(a)
	_ = mBar.RenderBlank()
}

func myCallback(a interface{}) {
	plg := a.(*plugins.Plugins)
	mFileSync.Lock()
	defer mFileSync.Unlock()
	_ = mCsvWriter.Write([]string{plg.TargetIp, plg.TargetPort, plg.TargetProtocol,
		fmt.Sprintf("%v", plg.Cracked)})
	mCsvWriter.Flush()
}

func myBarCallback(i int) {
	_ = mBar.Add(i)
}

func NewBar(max int) *progressbar.ProgressBar {
	bar := progressbar.NewOptions(max,
		progressbar.OptionSetWriter(ansi.NewAnsiStdout()),
		progressbar.OptionSetDescription("Loading ..."),
		progressbar.OptionSetWriter(os.Stderr),
		progressbar.OptionShowCount(),
		progressbar.OptionShowIts(),
		progressbar.OptionOnCompletion(func() {
			_, _ = fmt.Fprint(os.Stderr, "\n扫描任务完成")
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
	return bar
}
