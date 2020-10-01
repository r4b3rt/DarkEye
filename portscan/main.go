package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"github.com/zsdevX/DarkEye/common"
	"os"
	"strings"
	"sync"
)

var (
	mPort        = flag.String("port", "1-65535", "端口格式参考Nmap")
	mIp          = flag.String("ip", "127.0.0.1", "a.b.c.1-254")
	mActivePort  = flag.String("alive_port", "0", "使用已知开放的端口校正扫描行为。例如某服务器限制了IP访问频率，开启此功能后程序发现限制会自动调整保证扫描完整、准确")
	mTimeOut     = flag.Int("timeout", 2000, "扫描过程中每个端口的timeout时间；可以用-timeout_test参数来自动确认")
	mTestTimeOut = flag.Bool("timeout_test", false, "自动获取超时时间，互联网环境建议使用")
	mThread      = flag.Int("thread", 1, "仅扫描多个IP时有效，该参数可以控制每个线程扫描IP个数")
	mTitle       = flag.Bool("title", false, "获取标题，http/https有效")
	mMinTimeOut  = 100 //ms
)

var (
	mScans = make([]*Scan, 0)
)

func main() {
	help()
	flag.Parse()
	Start()
}

func Start() {
	ips := strings.Split(*mIp, ",")
	tot := 0
	for _, ip := range ips {
		base, start, end, err := common.ParseNmapIP(ip)
		if err != nil {
			fmt.Println("IP格式错误:", err.Error())
			return
		}
		for {
			if start > end {
				break
			}
			s := New(fmt.Sprintf("%s.%d", base, start))
			if err := s.InitConfig(); err != nil {
				fmt.Println(s.Ip, err.Error())
				os.Exit(0)
			}
			_, t := genFromTo(s.PortRange)
			tot += t
			mScans = append(mScans, s)
			start++
		}
	}
	fmt.Println(fmt.Sprintf("加载%d个IP,%d个端口", len(mScans), tot))
	Bar = NewBar(tot)
	wg := sync.WaitGroup{}
	wg.Add(*mThread)
	//先创建文件，创建失败结束
	f, err := os.Create(common.GenFileName("port_scan"))
	if err != nil {
		_, _ = fmt.Fprint(os.Stderr, "初始化失败", err.Error())
		return
	}
	defer f.Close()
	_, _ = f.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM
	w := csv.NewWriter(f)
	_ = w.Write([]string{"IP", "端口", "中间件", "标题"})

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

func Run(file *os.File, wg *sync.WaitGroup, id int) {
	i := 0
	max := len(mScans)
	defer wg.Done()
	for {
		if (*mThread*i)+id >= max {
			break
		}
		s := mScans[*mThread*i+id]
		_ = s.TimeOutTest()
		if len(mScans) != 1 {
			//多IP暂时不支持校正
			s.ActivePort = "0"
		}
		s.Run()
		s.OutPut(file)
		i++
	}
}

func help() {
	fmt.Println(common.Banner)
	fmt.Println(common.ProgramVersion)
	fmt.Println("Example1: ")
	fmt.Println("./portscan -alive_port 8443 -ip f.u.c.k -port 1-65535 -timeout_test")
	fmt.Println("Example2: ")
	fmt.Println("./portscan -ip f.u.c.k,f.u.c.1-254 -port 1-65535")
	fmt.Print("----------------\n\n")
}
