package main

import (
	"flag"
	"fmt"
	"github.com/zsdevX/DarkEye/common"
	"os"
	"strings"
	"sync"
)

var (
	mPort       = flag.String("port", "1-65535", "端口格式参考Nmap")
	mIp         = flag.String("ip", "127.0.0.1", "a.b.c.1-254")
	mActivePort = flag.String("alive_port", "0", "已知开放的端口用来校正扫描")
	mRate       = flag.Int("rate", 2000, "端口之间的扫描间隔单位ms，也可用通过-rate_test")
	mMinRate    = flag.Int("min_rate", 100, "自动计算的速率不能低于min_rate")
	mTestRate   = flag.Bool("rate_test", false, "发包频率")
	mOutputFile = flag.String("output", "result.csv", "结果保存到该文件")
	mThread     = flag.Int("thread", 1, "结果保存到该文件")
)

var (
	mScans = make([]*Scan, 0)
)

func main() {
	help()
	flag.Parse()
	if *mRate < 300 {
		fmt.Println("[WARN]: 建议rate设置在300ms以上比较稳")
	}
	Start()
}

func Start() {
	go BarDescription()
	ips := strings.Split(*mIp, ",")
	tot := 0
	for _, ip := range ips {
		base, start, end, err := common.ParseNmapIP(ip)
		if err != nil {
			fmt.Println("ip格式错误:", err.Error())
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
	fmt.Println(fmt.Sprintf("加载%d个IP %d个端口", len(mScans), tot))
	Bar = NewBar(tot)
	wg := sync.WaitGroup{}
	wg.Add(*mThread)
	i := 0
	for {
		if i >= *mThread {
			break
		}
		go Run(&wg, i)
		i++
	}
	wg.Wait()
}

func Run(wg *sync.WaitGroup, id int) {
	i := 0
	max := len(mScans)
	defer wg.Done()
	for {
		if (*mThread*i)+id >= max {
			break
		}
		s := mScans[*mThread*i+id]
		_ = s.RateTest()
		s.Run()
		s.OutPut()
		i++
	}
}

func help() {
	fmt.Println(common.Banner)
	fmt.Println(common.ProgramVersion)
	fmt.Println("Examples: ")
	fmt.Println("./lowSpeedPortScan -alive_port 8443 -ip f.u.c.k -port 1-65535 -rate_test -output result.txt")
	fmt.Println("./lowSpeedPortScan -ip f.u.c.k,f.u.c.1-254 -port 1-65535 -rate 200 -output result.txt")
	fmt.Print("----------------\n\n")
}
