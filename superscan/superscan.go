package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/fatih/color"
	"github.com/schollz/progressbar"
	"github.com/zsdevX/DarkEye/common"
	"github.com/zsdevX/DarkEye/superscan/plugins"
	"golang.org/x/time/rate"
	"runtime"
	"time"

	//	_ "net/http/pprof"
	"os"
	"strings"
	"sync"
)

var (
	mIp                    = flag.String("ip", "127.0.0.1", "a.b.c.1-254")
	mTimeOut               = flag.Int("timeout", 3000, "单位ms")
	mThread                = flag.Int("thread", 32, "扫单IP线程数")
	mPortList              = flag.String("port-list", common.PortList, "端口范围,默认1000+常用端口")
	mUserList              = flag.String("user-file", "", "用户名字典文件")
	mPassList              = flag.String("pass-file", "", "密码字典文件")
	mU                     = flag.String("U", "", "用户名字典:root,test")
	mP                     = flag.String("P", "", "密码:123456,1q2w3e")
	mNoTrust               = flag.Bool("no-trust", false, "由端口判定协议改为指纹方式判断协议,速度慢点")
	mPluginWorker          = flag.Int("plugin-worker", 2, "单协议爆破密码时，线程个数")
	mRateLimiter           = flag.Int("pps", 0, "扫描工具整体发包频率n/秒, 该选项可避免线程过多发包会占有带宽导致丢失目标端口")
	mActivePort            = flag.String("alive_port", "0", "使用已知开放的端口校正扫描行为。例如某服务器限制了IP访问频率，开启此功能后程序发现限制会自动调整保证扫描完整、准确")
	mListPlugin            = flag.Bool("list-plugin", false, "列出支持的爆破协议")
	mPocReverse            = flag.String("reverse-url", "qvn0kc.ceye.io", "CEye 标识")
	mPocReverseCheck       = flag.String("reverse-check-url", "http://api.ceye.io/v1/records?token=066f3d242991929c823ac85bb60f4313&type=http&filter=", "CEye API")
	mOnlyCheckAliveNetwork = flag.Bool("only-check-alive", false, "检查有活跃主机的网段")

	mMaxIPDetect = 32
	mFile        *os.File
	mCsvWriter   *csv.Writer
	mFileName    string
	mBar         *progressbar.ProgressBar
	mPps         *rate.Limiter
)

var (
	mScans = make([]*Scan, 0)
	//记录文件
	mFileSync = sync.RWMutex{}
)

func main() {
	//  debug/pprof
	/*
		go func() {
			fmt.Println(http.ListenAndServe("localhost:10000", nil))
		}()
	*/
	color.Yellow("超级弱口令、系统Vulnerable检测\n")
	flag.Parse()
	if *mListPlugin {
		plugins.SupportPlugin()
		return
	}
	loadPlugins()
	color.Red(common.Banner)
	common.SetRLimit()
	runtime.GOMAXPROCS(runtime.NumCPU())

	//活跃网段检测
	if *mOnlyCheckAliveNetwork {
		networkCheck()
		return
	}
	recordInit()
	Start()
}

//Start add comment
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
			nip := common.GenIP(base, start)
			if common.CompareIP(nip, end) > 0 {
				break
			}
			s := newScan(nip)
			s.ActivePort = "0"
			_, t := common.GetPortRange(s.PortRange)
			tot += t
			mScans = append(mScans, s)
			start++
		}
	}
	color.Green(fmt.Sprintf("已加载%d个IP,共计%d个端口,启动检测线程数%d,同时可检测IP数量%d",
		len(mScans), tot, *mThread, mMaxIPDetect))
	//建立进度条
	mBar = newBar(tot)
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

func networkCheck() {
	s := newScan("")
	s.PingNet(*mIp)
}

//修改插件参数
func loadPlugins() {
	//设置poc的反弹检测参数
	plugins.GlobalConfig.ReverseUrl = *mPocReverse
	plugins.GlobalConfig.ReverseCheckUrl = *mPocReverseCheck
	//设置自定义文件字典替代内置字典
	plugins.GlobalConfig.UserList = common.GenDicFromFile(*mUserList)
	plugins.GlobalConfig.PassList = common.GenDicFromFile(*mPassList)
	//设置字典替代内置字典
	if *mU != "" {
		plugins.GlobalConfig.UserList = strings.Split(*mU, ",")
	}
	if *mP != "" {
		plugins.GlobalConfig.PassList = strings.Split(*mP, ",")
	}
	//设置发包频率
	if *mRateLimiter > 0 {
		//每秒发包*mRateLimiter，缓冲10个
		mPps = rate.NewLimiter(rate.Every(1000000*time.Microsecond/time.Duration(*mRateLimiter)), 10)
		color.Green("rate limit enable <= %v pps\n", mPps.Limit())
	}
	plugins.GlobalConfig.Pps = mPps
}

//NewScan add comment
func newScan(ip string) *Scan {
	return &Scan{
		Ip:                     ip,
		TimeOut:                *mTimeOut,
		ActivePort:             *mActivePort,
		PortRange:              *mPortList,
		ThreadNumber:           *mThread,
		NoTrust:                *mNoTrust,
		PluginWorker:           *mPluginWorker,
		Callback:               myCallback,
		BarCallback:            myBarCallback,
		BarDescriptionCallback: myBarDescUpdate,
	}
}

func myBarDescUpdate(a string) {
	b := fmt.Sprintf("%-32s", a)
	if len(a) > 32 {
		b = a[:(32-3)] + "..."
	}
	mBar.Describe(b)
}

func myCallback(a interface{}) {
	plg := a.(*plugins.Plugins)
	mFileSync.Lock()
	defer mFileSync.Unlock()
	ck, _ := json.Marshal(plg.Cracked)
	_ = mCsvWriter.Write([]string{plg.TargetIp, plg.TargetPort, plg.TargetProtocol,
		string(ck)})
	mCsvWriter.Flush()
}

func myBarCallback(i int) {
	_ = mBar.Add(i)
}

func newBar(max int) *progressbar.ProgressBar {
	bar := progressbar.NewOptions(max,
		progressbar.OptionSetDescription(fmt.Sprintf("%-32s", "Cracking...")),
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

func recordInit() {
	var err error
	mFileName, _ = common.Write2CSV("superScan", nil)
	f, err := os.Create(mFileName)
	if err != nil {
		return
	}
	_, _ = f.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM
	w := csv.NewWriter(f)
	_ = w.Write([]string{"IP", "端口", "协议", "插件信息"})
	mCsvWriter = csv.NewWriter(f)
	color.Yellow("记录结果将保存在%s", mFileName)
}
