package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/fatih/color"
	"github.com/k0kubun/go-ansi"
	"github.com/schollz/progressbar"
	"github.com/zsdevX/DarkEye/common"
	"github.com/zsdevX/DarkEye/superscan/plugins"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	"strings"
	"sync"
	"syscall"
)

var (
	mIp           = flag.String("ip", "127.0.0.1", "a.b.c.1-254")
	mTimeOut      = flag.Int("timeout", 3000, "单位ms")
	mThread       = flag.Int("thread", 32, "扫单IP线程数")
	mPortList     = flag.String("port-list", common.PortList, "端口范围,默认1000+常用端口")
	mNoTrust      = flag.Bool("no-trust", false, "由端口判定协议改为指纹方式判断协议,速度慢点")
	mPluginWorker = flag.Int("plugin-worker", 2, "单协议爆破密码时，线程个数")
	mActivePort   = flag.String("alive_port", "0", "使用已知开放的端口校正扫描行为。例如某服务器限制了IP访问频率，开启此功能后程序发现限制会自动调整保证扫描完整、准确")
	mMaxIPDetect  = 16
	mFile         *os.File
	mCsvWriter    *csv.Writer
	mFileName     string
	mBar          *progressbar.ProgressBar
)

var (
	mScans = make([]*Scan, 0)
	//记录文件
	mFileSync = sync.RWMutex{}
)

type BarValue struct {
	Key   string
	Value string
}

func init() {
	var err error
	_, mFile, mFileName, err = common.CreateCSV("portScan",
		[]string{"IP", "端口", "插件扫描信息"})
	if err != nil {
		panic("创建记录文件失败" + err.Error())
	}
	mCsvWriter = csv.NewWriter(mFile)
	fmt.Println("记录结果将保存在", mFileName)
}

func main() {
	color.Red(common.Banner)
	color.Yellow("\n一键端口发现、POC检测、弱口令检测\n\n")
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())
	//  debug/pprof/
	go func() {
		fmt.Println(http.ListenAndServe("localhost:10000", nil))
	}()

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
	//设置max file
	rLimit := syscall.Rlimit{
		Cur: 65535,
		Max: 65535,
	}
	_ = syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	_ = syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	color.Green(fmt.Sprintf("已加载%d个IP,共计%d个端口,启动检测线程数%d,同时可检测IP数量%d,系统资源上限为%v",
		len(mScans), tot, *mThread, mMaxIPDetect, rLimit))
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
		PluginWorker:           *mPluginWorker,
		PortsScannedOpened:     make([]plugins.Plugins, 0),
		Callback:               myCallback,
		BarCallback:            myBarCallback,
		BarDescriptionCallback: myBarDescUpdate,
	}
}

func myBarDescUpdate(a string) {
	mBar.Describe(a)
	mBar.Add(0)

}

func myCallback(result []byte) {
	plg := plugins.Plugins{}
	_ = json.Unmarshal(result, &plg)
	mFileSync.Lock()
	defer mFileSync.Unlock()
	_ = mCsvWriter.Write([]string{plg.TargetIp, string(result)})
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
		progressbar.OptionFullWidth(),
	)

	_ = bar.RenderBlank()
	return bar
}
