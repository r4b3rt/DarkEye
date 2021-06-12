package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/schollz/progressbar"
	"github.com/sirupsen/logrus"
	"github.com/b1gcat/DarkEye/common"
	"github.com/b1gcat/DarkEye/superscan"
	"github.com/b1gcat/DarkEye/superscan/plugins"
	"golang.org/x/time/rate"
	"os"
	"strings"
	"time"
)

type superScanRuntime struct {
	Attack                bool
	Ip                    string
	PortList              string
	TimeOut               int
	Thread                int
	Plugin                string
	PacketPerSecond       int
	User                  string
	Pass                  string
	WebSiteDomain         string
	ActivePort            string
	Output                string
	OnlyCheckAliveNetwork bool
	OnlyCheckAliveHost    bool

	MaxConcurrencyIp int
	Bar              *progressbar.ProgressBar
	PacketRate       *rate.Limiter
	scan             *superscan.Scan
	flagSet          *flag.FlagSet
	result           []analysisEntity
	ctx              context.Context
}

var (
	superScan               = "superScan"
	superScanRuntimeOptions = &superScanRuntime{
		flagSet: flag.NewFlagSet(superScan, flag.ExitOnError),
		result:  make([]analysisEntity, 0),
		ctx:     context.Background(),
	}
)

func (s *superScanRuntime) Start() {
	s.initializer()

	if s.OnlyCheckAliveNetwork || s.OnlyCheckAliveHost {
		scan := s.newScan("")
		scan.PingNet(s.Ip, s.OnlyCheckAliveHost)
		return
	}
	//初始化scan对象
	ips := common.ParseFileOrVariable(s.Ip)
	if len(ips) == 0 {
		common.Log("superScan.start", "目标空", common.ALERT)
		return
	}
	tot := 0
	scans := make([]*superscan.Scan, 0)
	for _, ip := range ips {
		base, start, end, err := common.GetIPRange(ip)
		if err != nil {
			common.Log(superScan, err.Error(), common.FAULT)
			return
		}
		for {
			nip := common.GenIP(base, start)
			if common.CompareIP(nip, end) > 0 {
				break
			}
			ss := s.newScan(nip)
			ss.ActivePort = "0"
			ss.Parent = s.ctx
			_, t := common.GetPortRange(ss.PortRange)
			tot += t
			scans = append(scans, ss)
			start++
		}
	}
	logrus.Info(fmt.Sprintf(
		"已加载%d个IP,共计%d个端口,启动每IP扫描端口线程数%d,同时可同时检测IP数量%d",
		len(scans), tot, s.Thread, s.MaxConcurrencyIp))
	plugins.Plugin()

	//建立进度条
	s.Bar = s.newBar(tot)
	if len(scans) == 1 {
		//单IP支持校正
		scans[0].ActivePort = s.ActivePort
	}
	task := common.NewTask(s.MaxConcurrencyIp, s.ctx)
	defer task.Wait("superScan")
	for _, sc := range scans {
		//Job
		if !task.Job() {
			break
		}
		go func(sup *superscan.Scan) {
			defer task.UnJob()
			sup.Run()
		}(sc)
	}
	s.display()
}

func (s *superScanRuntime) Init() {
	s.flagSet.BoolVar(&s.Attack, "attack", false, "发现漏洞即刻攻击")
	s.flagSet.StringVar(&s.Ip, "ip", "127.0.0.1", "a.b.c.1-a.b.c.255(多组逗号隔开)或文件(一行一组IP范围)")
	s.flagSet.StringVar(&s.PortList, "port-list", common.PortList, "端口范围,默认1000+常用端口")
	s.flagSet.IntVar(&s.TimeOut, "timeout", 3000, "网络超时请求(单位ms)")
	s.flagSet.IntVar(&s.Thread, "thread", 128, "每个IP爆破端口的线程数量")
	s.flagSet.IntVar(&s.PacketPerSecond, "pps", 0, "扫描工具整体发包频率 packets/秒")
	s.flagSet.StringVar(&s.Plugin, "plugin", "", "指定协议插件爆破")
	s.flagSet.StringVar(&s.User, "user", "", "用户名(多个逗号隔开)或文件(一行一个账号)")
	s.flagSet.StringVar(&s.Pass, "pass", "", "密码(多个逗号隔开)或文件(一行一个密码)")
	s.flagSet.StringVar(&s.ActivePort, "alive_port", "0", "使用已知开放的端口校正扫描行为。例如某服务器限制了IP访问频率，开启此功能后程序发现限制会自动调整保证扫描完整、准确")
	s.flagSet.BoolVar(&s.OnlyCheckAliveNetwork, "only-alive-network", false, "只检查活跃主机的网段(ping)")
	s.flagSet.BoolVar(&s.OnlyCheckAliveHost, "alive-host-check", false, "检查所有活跃主机(ping)")
	s.flagSet.StringVar(&s.Output, "output", "superScan.csv", "输出文件")
	s.flagSet.StringVar(&s.WebSiteDomain, "website-domain", "", "域名(多个逗号隔开)或文件(一行一个域名)")
	s.MaxConcurrencyIp = 32
	s.flagSet.Parse(os.Args[1:])

}

func (s *superScanRuntime) newScan(ip string) *superscan.Scan {
	return &superscan.Scan{
		Ip:          ip,
		TimeOut:     s.TimeOut,
		ActivePort:  s.ActivePort,
		PortRange:   s.PortList,
		Thread:      s.Thread,
		Callback:    s.myCallback,
		BarCallback: s.myBarCallback,
	}
}

func (s *superScanRuntime) initializer() {
	//设置自定义文件字典替代内置字典
	plugins.Config.UserList = common.ParseFileOrVariable(s.User)
	plugins.Config.PassList = common.ParseFileOrVariable(s.Pass)
	plugins.Config.WebSiteDomainList = common.ParseFileOrVariable(s.WebSiteDomain)
	//设置发包频率
	if s.PacketPerSecond > 0 {
		//每秒发包*mRateLimiter，缓冲10个
		s.PacketRate = rate.NewLimiter(rate.Every(1000000*time.Microsecond/time.Duration(s.PacketPerSecond)), 10)
	}
	plugins.Config.PPS = s.PacketRate
	plugins.Config.SelectPlugin = s.Plugin
	plugins.Config.ParentCtx = s.ctx
	plugins.Config.TimeOut = s.TimeOut
	plugins.Config.Attack = s.Attack
	plugins.Config.UpdateDesc = s.myBarChangeDesc
}

func (s *superScanRuntime) newBar(max int) *progressbar.ProgressBar {
	bar := progressbar.NewOptions(max,
		progressbar.OptionSetDescription(fmt.Sprintf("%-24s", "Cracking...")),
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
		Service: plg.Result.ServiceName,
	}
	if err := mapstructure.Decode(plg.Result.Output.Items(), &ent); err != nil {
		common.Log(superScan, err.Error(), common.FAULT)
		return
	}
	s.result = append(s.result, ent)
	s.OutPut()
}

func (s *superScanRuntime) myBarCallback(i int) {
	_ = s.Bar.Add(i)
}

func (s *superScanRuntime) myBarChangeDesc(a interface{}, args ...string) {
	plg := a.(*plugins.Plugins)
	ip := strings.Split(plg.TargetIp, ".")
	desc := args[0] + "://" + "*" + ip[2] + "." + ip[3] + "/" + args[1]
	b := fmt.Sprintf("%-24s", desc)
	if len(desc) > 24 {
		b = desc[:24]
	}
	s.Bar.Describe(b)
	_ = s.Bar.RenderBlank()
}
