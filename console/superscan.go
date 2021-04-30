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
	"strconv"
	"strings"
	"sync"
	"time"
)

type superScanRuntime struct {
	Module
	parent *RequestContext

	Attack                bool
	IpList                string
	PortList              string
	TimeOut               int
	Thread                int
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
	cmd []string
}

var (
	superScan               = "superScan"
	superScanRuntimeOptions = &superScanRuntime{
		flagSet: flag.NewFlagSet(superScan, flag.ExitOnError),
	}
)

func (s *superScanRuntime) Start(parent context.Context) {
	s.initializer(parent)

	if s.OnlyCheckAliveNetwork || s.OnlyCheckAliveHost {
		scan := s.newScan("")
		scan.PingNet(s.IpList, s.OnlyCheckAliveHost)
		return
	}
	//解析变量
	ipList, err := analysisRuntimeOptions.Var("", s.IpList)
	if err != nil {
	} else {
		s.IpList = ""
		for _, v := range ipList {
			s.IpList += v + ","
		}
		s.IpList = strings.TrimSuffix(s.IpList, ",")
	}
	//初始化scan对象
	ips := strings.Split(s.IpList, ",")
	if len(ips) == 0 {
		common.Log("superScan.start", "目标空", common.ALERT)
		return
	}
	tot := 0
	scans := make([]*superscan.Scan, 0)
	for _, ip := range ips {
		base, start, end, err := common.GetIPRange(ip)
		if err != nil {
			common.Log(s.parent.CmdArgs[0], err.Error(), common.FAULT)
			return
		}
		for {
			nip := common.GenIP(base, start)
			if common.CompareIP(nip, end) > 0 {
				break
			}
			s := s.newScan(nip)
			s.ActivePort = "0"
			s.Parent = parent
			_, t := common.GetPortRange(s.PortRange)
			tot += t
			scans = append(scans, s)
			start++
		}
	}
	fmt.Println(fmt.Sprintf(
		"已加载%d个IP,共计%d个端口,启动每IP扫描端口线程数%d,同时可同时检测IP数量%d",
		len(scans), tot, s.Thread, s.MaxConcurrencyIp))
	plugins.SupportPlugin()

	//建立进度条
	s.Bar = s.newBar(tot)
	if len(scans) == 1 {
		//单IP支持校正
		scans[0].ActivePort = s.ActivePort
	}
	task := common.NewTask(s.MaxConcurrencyIp, parent)
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
}

func (s *superScanRuntime) Init(requestContext *RequestContext) {
	s.parent = requestContext
	s.flagSet.BoolVar(&s.Attack, "attack", false, "发现漏洞即刻攻击")
	s.flagSet.StringVar(&s.IpList, "ip", "127.0.0.1", "a.b.c.1-a.b.c.255")
	s.flagSet.StringVar(&s.PortList, "port-list", common.PortList, "端口范围,默认1000+常用端口")
	s.flagSet.IntVar(&s.TimeOut, "timeout", 3000, "网络超时请求(单位ms)")
	s.flagSet.IntVar(&s.Thread, "thread", 128, "每个IP爆破端口的线程数量")
	s.flagSet.IntVar(&s.PacketPerSecond, "pps", 0, "扫描工具整体发包频率 packets/秒")
	s.flagSet.StringVar(&s.Plugin, "plugin", "", "指定协议插件爆破")
	s.flagSet.StringVar(&s.UserList, "user-list", "", "字符串(u1,u2,u3)或文件(一个账号一行）")
	s.flagSet.StringVar(&s.PassList, "pass-list", "", "字符串(p1,p2,p3)或文件（一个密码一行")
	s.flagSet.StringVar(&s.ActivePort, "alive_port", "0", "使用已知开放的端口校正扫描行为。例如某服务器限制了IP访问频率，开启此功能后程序发现限制会自动调整保证扫描完整、准确")
	s.flagSet.BoolVar(&s.OnlyCheckAliveNetwork, "only-alive-network", false, "只检查活跃主机的网段(ping)")
	s.flagSet.BoolVar(&s.OnlyCheckAliveHost, "alive-host-check", false, "检查所有活跃主机(ping)")
	s.MaxConcurrencyIp = 32
}

func (s *superScanRuntime) ValueCheck(value string) (bool, error) {
	if v, ok := superScanValueCheck[value]; ok {
		if isDuplicateArg(value, s.parent.CmdArgs) {
			return false, fmt.Errorf("参数重复")
		}
		return v, nil
	}
	return false, fmt.Errorf("无此参数")
}

func (s *superScanRuntime) CompileArgs(cmd []string, os []string) error {
	if cmd != nil {
		if err := s.flagSet.Parse(splitCmd(cmd)); err != nil {
			return err
		}
		s.flagSet.Parsed()
	} else {
		if err := s.flagSet.Parse(os); err != nil {
			return err
		}
	}
	return nil
}

func (a *superScanRuntime) saveCmd(cmd []string) {
	a.cmd = cmdSave(cmd)
}

func (a *superScanRuntime) restoreCmd() []string {
	return cmdRestore(a.cmd)
}

func (a *superScanRuntime) Usage() {
	fmt.Println(fmt.Sprintf("Usage of %s:", superScan))
	fmt.Println("Options:")
	a.flagSet.VisitAll(func(f *flag.Flag) {
		var opt = "  -" + f.Name
		fmt.Println(opt)
		fmt.Println(fmt.Sprintf("		%v (default '%v')", f.Usage, f.DefValue))
	})
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

func (s *superScanRuntime) initializer(parent context.Context) {
	//设置自定义文件字典替代内置字典
	if s.UserList != "" {
		if _, e := os.Stat(s.UserList); e != nil {
			plugins.Config.UserList = common.GenDicFromFile(s.UserList)
		} else {
			plugins.Config.UserList = strings.Split(s.UserList, ",")
		}
	}
	if s.PassList != "" {
		if _, e := os.Stat(s.PassList); e != nil {
			plugins.Config.PassList = common.GenDicFromFile(s.PassList)
		} else {
			plugins.Config.PassList = strings.Split(s.PassList, ",")
		}
	}
	//设置发包频率
	if s.PacketPerSecond > 0 {
		//每秒发包*mRateLimiter，缓冲10个
		s.PacketRate = rate.NewLimiter(rate.Every(1000000*time.Microsecond/time.Duration(s.PacketPerSecond)), 10)
	}
	plugins.Config.PPS = s.PacketRate
	plugins.Config.SelectPlugin = s.Plugin
	plugins.Config.ParentCtx = parent
	plugins.Config.TimeOut = s.TimeOut
	plugins.Config.Attack = s.Attack
	s.parent.taskId ++
}

func (s *superScanRuntime) newBar(max int) *progressbar.ProgressBar {
	bar := progressbar.NewOptions(max,
		progressbar.OptionSetDescription("[Cracking ...]"),
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
		Task:      strconv.Itoa(s.parent.taskId),
		Ip:        plg.TargetIp,
		Port:      plg.TargetPort,
		Service:   plg.Result.ServiceName,
		ExpHelper: plg.Result.ExpHelp,
	}
	if plg.Result.Web.Url != "" {
		ent.Url = plg.Result.Web.Url
		ent.Title = plg.Result.Web.Title
		ent.WebServer = plg.Result.Web.Server
		ent.WebResponseCode = plg.Result.Web.Code
	}
	if plg.Result.NetBios.Net != "" ||
		plg.Result.NetBios.Os != "" ||
		plg.Result.NetBios.Shares != "" {
		ent.NetBios = fmt.Sprintf(" ['%s' '%s' '%s', '%s', '%s']",
			plg.Result.NetBios.Net, plg.Result.NetBios.Hw,
			plg.Result.NetBios.Shares,
			plg.Result.NetBios.Domain,
			plg.Result.NetBios.UserName)
	}
	if plg.Result.Cracked.Username != "" ||
		plg.Result.Cracked.Password != "" {
		ent.WeakAccount = fmt.Sprintf(
			"[%s/%s]", plg.Result.Cracked.Username, plg.Result.Cracked.Password)
	}
	analysisRuntimeOptions.upInsertEnt(&ent)
	analysisRuntimeOptions.PrintCurrentTaskResult()

}

func (s *superScanRuntime) myBarCallback(i int) {
	_ = s.Bar.Add(i)
}
