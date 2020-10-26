package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"github.com/k0kubun/go-ansi"
	"github.com/schollz/progressbar"
	"github.com/zsdevX/DarkEye/common"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	mPort                = flag.String("port", "", "端口格式参考Nmap,默认为常用端口")
	mIp                  = flag.String("ip", "127.0.0.1", "a.b.c.1-254")
	mActivePort          = flag.String("alive_port", "0", "使用已知开放的端口校正扫描行为。例如某服务器限制了IP访问频率，开启此功能后程序发现限制会自动调整保证扫描完整、准确")
	mTimeOut             = flag.Int("timeout", 2000, "扫描过程中每个端口的timeout时间；可以用-timeout_test参数来自动确认")
	mTestTimeOut         = flag.Bool("timeout_test", false, "自动获取超时时间，互联网环境建议使用")
	mThread              = flag.Int("thread", 1, "该参数可以控制每个线程扫描IP个数")
	mTitle               = flag.Bool("title", false, "获取标题，http/https有效")
	mPortRangeThresHolds = flag.Int("port-range-thresholds", 1000, "端口范围大于阀值会触发多线程扫描，线程数通过mThread获取")
	mMinTimeOut          = 100 //ms
)

var (
	mScans = make([]*Scan, 0)
	//进度条
	Bar     = &progressbar.ProgressBar{}
	BarDesc = make(chan *BarValue, 64)
	//记录文件
	mFileSync = sync.RWMutex{}
)

type BarValue struct {
	Key   string
	Value string
}

func main() {
	fmt.Println(common.Banner)
	if len(os.Args) == 1 {
		help()
		return
	}
	flag.Parse()
	if !*mTitle {
		fmt.Println("***未开启指纹功能***", "查看帮助如何开启./portscan -h")
	}
	Start()
}

func Start() {
	_, f, fileName, err := common.CreateCSV("portScan",
		[]string{"IP", "端口", "中间件", "标题"})
	if err != nil {
		fmt.Println("创建记录文件失败", err.Error())
		return
	}
	fmt.Println("记录结果将保存在", fileName)
	defer f.Close()
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
			s := NewScan(fmt.Sprintf("%s.%d", base, start))
			_, t := common.GetPortRange(s.PortRange)
			tot += t
			mScans = append(mScans, s)
			start++
		}
	}
	fmt.Println(fmt.Sprintf("加载%d个IP,%d个端口", len(mScans), tot))
	//建立进度条
	Bar = NewBar(tot)
	//创建任务
	wg := sync.WaitGroup{}
	wg.Add(*mThread)
	i := 0
	for {
		if i >= *mThread {
			break
		}
		go Run(f, &wg, i)
		i++
	}
	wg.Wait()
}

func NewScan(ip string) *Scan {
	portList := PortList
	if *mPort != "" {
		portList = *mPort
	}
	return &Scan{
		Ip:                  ip,
		DefaultTimeOut:      *mTimeOut,
		ActivePort:          *mActivePort,
		MinTimeOut:          mMinTimeOut,
		PortRange:           portList,
		Test:                *mTestTimeOut,
		Title:               *mTitle,
		PortsScannedOpened:  make([]PortInfo, 0),
		Callback:            myCallback,
		BarCallback:         myBarCallback,
		PortRangeThresholds: *mPortRangeThresHolds,
		ThreadNumber:        *mThread,
	}
}

func Run(file *os.File, wg *sync.WaitGroup, id int) {
	i := 0
	max := len(mScans)
	defer wg.Done()
	for (*mThread*i)+id < max {

		s := mScans[*mThread*i+id]
		_ = s.TimeOutTest()
		if len(mScans) != 1 {
			//多IP暂时不支持校正
			s.ActivePort = "0"
		}
		s.Run()
		OutPut(file, s)
		i++
	}
}

func OutPut(f *os.File, s *Scan) {
	mFileSync.Lock()
	defer mFileSync.Unlock()

	w := csv.NewWriter(f)
	for _, p := range s.PortsScannedOpened {
		_ = w.Write([]string{s.Ip, strconv.Itoa(p.Port), p.Server, p.Title})
	}
	w.Flush()
}

func myCallback(a ...interface{}) {
	fmt.Println(a...)
}

func myBarCallback(i int) {
	_ = Bar.Add(i)
}

func NewBar(max int) *progressbar.ProgressBar {
	bar := progressbar.NewOptions(max,
		progressbar.OptionSetWriter(ansi.NewAnsiStdout()),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionThrottle(3000*time.Millisecond),
		progressbar.OptionSetDescription("Loading ..."),
		progressbar.OptionSetWriter(os.Stderr),
		progressbar.OptionSetWidth(10),
		progressbar.OptionShowCount(),
		progressbar.OptionShowIts(),
		progressbar.OptionOnCompletion(func() {
			_, _ = fmt.Fprint(os.Stderr, "\nDONE")
		}),
		progressbar.OptionSpinnerType(14),
		progressbar.OptionFullWidth(),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}))

	_ = bar.RenderBlank()
	return bar
}

func help() {
	fmt.Println("Example1: ")
	fmt.Println("./portscan -alive_port 8443 -ip f.u.c.k -port 1-65535 -timeout_test")
	fmt.Println("Example2: ")
	fmt.Println("./portscan -ip f.u.c.k,f.u.c.1-254 -port 1-65535 -title -thread 16 -timeout 200")
}
