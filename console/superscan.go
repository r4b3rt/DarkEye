package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/schollz/progressbar"
	"github.com/zsdevX/DarkEye/common"
	"github.com/zsdevX/DarkEye/superscan"
	"github.com/zsdevX/DarkEye/superscan/plugins"
	"golang.org/x/time/rate"
	"os"
	"strings"
	"sync"
	"time"
)

type superScanRuntime struct {
	IpList                string
	PortList              string
	TimeOut               int
	Thread                int
	PluginThread          int
	Plugin                string
	PacketPerSecond       int
	UserList              string
	PassList              string
	ActivePort            string
	OnlyCheckAliveNetwork bool
	OnlyCheckAliveHost    bool

	MaxConcurrencyIp int
	Bar              *progressbar.ProgressBar
	PacketRate       *rate.Limiter
	scan             *superscan.Scan
	flagSet          *flag.FlagSet
	sync.RWMutex
}

var (
	superScanRuntimeOptions = &superScanRuntime{
		flagSet: flag.NewFlagSet("superScan", flag.ExitOnError),
	}
	mScans = make([]*superscan.Scan, 0)
)

func (a *superScanRuntime) compileArgs(cmd []string) error {
	if err := a.flagSet.Parse(splitCmd(cmd)); err != nil {
		return err
	}
	a.flagSet.Parsed()
	return nil
}

func (a *superScanRuntime) usage() {
	fmt.Println(fmt.Sprintf("Usage of %s:", superScan))
	fmt.Println("Options:")
	a.flagSet.VisitAll(func(f *flag.Flag) {
		var opt = "  -" + f.Name
		fmt.Println(opt)
		fmt.Println(fmt.Sprintf("		%v (default '%v')", f.Usage, f.DefValue))
	})
}

func superScanInitRunTime() {
	superScanRuntimeOptions.flagSet.StringVar(&superScanRuntimeOptions.IpList, "ip", "127.0.0.1", "a.b.c.1-a.b.c.255")
	superScanRuntimeOptions.flagSet.StringVar(&superScanRuntimeOptions.PortList, "port-list", common.PortList, "端口范围,默认1000+常用端口")
	superScanRuntimeOptions.flagSet.IntVar(&superScanRuntimeOptions.TimeOut, "timeout", 3000, "网络超时请求(单位ms)")
	superScanRuntimeOptions.flagSet.IntVar(&superScanRuntimeOptions.Thread, "thread", 128, "每个IP爆破端口的线程数量")
	superScanRuntimeOptions.flagSet.IntVar(&superScanRuntimeOptions.PluginThread, "plugin-thread", 2, "每个协议爆破弱口令的线程数量")
	superScanRuntimeOptions.flagSet.IntVar(&superScanRuntimeOptions.PacketPerSecond, "pps", 0, "扫描工具整体发包频率 packets/秒")
	superScanRuntimeOptions.flagSet.StringVar(&superScanRuntimeOptions.Plugin, "plugin", "", "指定协议插件爆破")
	superScanRuntimeOptions.flagSet.StringVar(&superScanRuntimeOptions.UserList, "user-list", "", "字符串(u1,u2,u3)或文件(一个账号一行）")
	superScanRuntimeOptions.flagSet.StringVar(&superScanRuntimeOptions.PassList, "pass-list", "", "字符串(p1,p2,p3)或文件（一个密码一行")
	superScanRuntimeOptions.flagSet.StringVar(&superScanRuntimeOptions.ActivePort, "alive_port", "0", "使用已知开放的端口校正扫描行为。例如某服务器限制了IP访问频率，开启此功能后程序发现限制会自动调整保证扫描完整、准确")
	superScanRuntimeOptions.flagSet.BoolVar(&superScanRuntimeOptions.OnlyCheckAliveNetwork, "only-alive-network", false, "只检查活跃主机的网段(ping)")
	superScanRuntimeOptions.flagSet.BoolVar(&superScanRuntimeOptions.OnlyCheckAliveHost, "alive-host-check", false, "检查所有活跃主机(ping)")
	superScanRuntimeOptions.MaxConcurrencyIp = 32
}

func (s *superScanRuntime) start(ctx context.Context) {
	s.initializer(ctx)

	if s.OnlyCheckAliveNetwork || s.OnlyCheckAliveHost {
		scan := s.newScan("")
		scan.PingNet(s.IpList, s.OnlyCheckAliveHost)
		return
	}
	//初始化scan对象
	ips := strings.Split(s.IpList, ",")
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
			s := s.newScan(nip)
			s.ActivePort = "0"
			_, t := common.GetPortRange(s.PortRange)
			tot += t
			mScans = append(mScans, s)
			start++
		}
	}
	fmt.Println(fmt.Sprintf(
		"已加载%d个IP,共计%d个端口,启动每IP扫描端口线程数%d,启动每协议爆破线程数量%d,同时可同时检测IP数量%d",
		len(mScans), tot, s.Thread, s.PluginThread, s.MaxConcurrencyIp))
	fmt.Println("支持的爆破协议:")
	plugins.SupportPlugin()

	//建立进度条
	s.Bar = s.newBar(tot)
	if len(mScans) == 1 {
		//单IP支持校正
		mScans[0].ActivePort = s.ActivePort
	}
	limiter := make(chan int, s.MaxConcurrencyIp)
	wg := sync.WaitGroup{}
	wg.Add(len(mScans))
	for _, sc := range mScans {
		//Job
		go func(s0 *superscan.Scan) {
			defer wg.Done()
			//Cancel
			if plugins.ShouldStop() {
				return
			}
			limiter <- 1
			s0.Run()
			<-limiter
		}(sc)
	}
	wg.Wait()
}

func (s *superScanRuntime) newScan(ip string) *superscan.Scan {
	return &superscan.Scan{
		Ip:           ip,
		TimeOut:      superScanRuntimeOptions.TimeOut,
		ActivePort:   superScanRuntimeOptions.ActivePort,
		PortRange:    superScanRuntimeOptions.PortList,
		ThreadNumber: superScanRuntimeOptions.Thread,
		Callback:     s.myCallback,
		BarCallback:  s.myBarCallback,
	}
}

func (s *superScanRuntime) initializer(ctx context.Context) {
	//设置自定义文件字典替代内置字典
	if s.UserList != "" {
		if _, e := os.Stat(s.UserList); e != nil {
			plugins.GlobalConfig.UserList = common.GenDicFromFile(s.UserList)
		} else {
			plugins.GlobalConfig.UserList = strings.Split(s.UserList, ",")
		}
	}
	if s.PassList != "" {
		if _, e := os.Stat(s.PassList); e != nil {
			plugins.GlobalConfig.PassList = common.GenDicFromFile(s.PassList)
		} else {
			plugins.GlobalConfig.PassList = strings.Split(s.PassList, ",")
		}
	}
	//设置发包频率
	if s.PacketPerSecond > 0 {
		//每秒发包*mRateLimiter，缓冲10个
		s.PacketRate = rate.NewLimiter(rate.Every(1000000*time.Microsecond/time.Duration(s.PacketPerSecond)), 10)
	}
	plugins.GlobalConfig.Pps = s.PacketRate
	plugins.GlobalConfig.UsingPlugin = s.Plugin
	plugins.GlobalConfig.Ctx = ctx
	plugins.GlobalConfig.Thread = s.PluginThread
}

func (s *superScanRuntime) newBar(max int) *progressbar.ProgressBar {
	bar := progressbar.NewOptions(max,
		progressbar.OptionSetDescription("["+s.IpList+"]"),
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

func (s *superScanRuntime) myCallback(a interface{}) {
	plg := a.(*plugins.Plugins)
	ent := analysisEntity{
		Ip:      plg.TargetIp,
		Port:    plg.TargetPort,
		Service: plg.TargetProtocol,
		Os:      plg.NetBios.Os,
		NetBios: fmt.Sprintf(
			"[Ip:'%s' Shares:'%s']", plg.NetBios.Ip, plg.NetBios.Shares),
		Url:             plg.Web.Url,
		Title:           plg.Web.Title,
		WebServer:       plg.Web.Server,
		WebResponseCode: plg.Web.Code,
	}
	for k, v := range plg.Cracked {
		ent.WeakAccount += fmt.Sprintf("%d:%s/%s.Username;", k, v.Username, v.Password)
	}
	analysisRuntimeOptions.createOrUpdate(&ent)
}

func (s *superScanRuntime) myBarCallback(i int) {
	_ = s.Bar.Add(i)
}
