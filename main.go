package main

import (
	"context"
	"flag"
	"runtime"
	"strings"

	"github.com/b1gcat/DarkEye/scan"
	"github.com/sirupsen/logrus"
)

type config struct {
	ip      string
	port    string
	loaders string
	action  string
	timeout int

	host string //ip-host collision
	user string //user dict for risk
	pass string //pass dict for risk

	outfile  string //output file
	bar      bool   //progress
	progress *progress

	maxThreadForEachScan   int
	maxThreadForEachIPScan int
	debug                  bool

	//cache
	ctx    context.Context
	cancel context.CancelFunc
}

type myScan struct {
	s      scan.Scan
	p      *pool //for ip
	action actionType
	sid    scan.IdType //scan id
	bar    *bar
	total  int64
}

var (
	gConfig = &config{
		progress: &progress{},
	}
	Version = "v5.0.0"
)

func main() {
	initialize()
	logrus.Info(Version)
	gConfig.loader()
}

func initialize() {
	flag.StringVar(&gConfig.action, "action", actionDiscoNet.String(),
		"Format: "+myActionList.String())
	flag.StringVar(&gConfig.loaders, "loader", scan.IdList.String(),
		"Support loader: "+scan.IdList.String())
	flag.StringVar(&gConfig.ip, "ip", "127.0.0.1-254",
		"Format: 1.1.1.1-254,2.1.1.1,3.1.1.1/24")
	flag.StringVar(&gConfig.host, "host", "www.baidu.com",
		"Format: www.a.com,www.b.com OR host.txt")
	flag.StringVar(&gConfig.user, "user", "",
		"Format: user1,user2 OR user.txt")
	flag.StringVar(&gConfig.pass, "pass", "",
		"Format: pass1,pass2 OR pass.txt")
	flag.StringVar(&gConfig.outfile, "o", "output.txt",
		"output to file")
	flag.BoolVar(&gConfig.bar, "bar", false,
		"show progress for each loader")
	flag.IntVar(&gConfig.timeout,
		"timeout", 3000, "Format: 2000")
	flag.StringVar(&gConfig.port,
		"p",
		"21,23,80-89,389,443,445,512,513,514,873,"+
			"1025,111,1433,1521,2082,2083,2222,2375,2601,2604,3128,3306,3312,3311,3389,"+
			"4430-4445,5432,5900,5984,6082,6379,"+
			"7001,7002,7778,8000-9090,8080,8083,8649,8888,9200,9300,10000,11211,27017,27018,28017,50000,50070,50030,"+
			"554,53,110,1080,22",
		"Format: 80,80-81")
	flag.BoolVar(&gConfig.debug,
		"v", false, "for debugger")
	flag.IntVar(&gConfig.maxThreadForEachScan,
		"t", 32, "thread for every service")
	flag.IntVar(&gConfig.maxThreadForEachIPScan,
		"tt", 100, "thread for every ip")
	flag.Parse()

	if gConfig.debug {
		logrus.SetLevel(logrus.DebugLevel)
	}
	if runtime.NumCPU() > 1 {
		runtime.GOMAXPROCS(runtime.NumCPU() - 1)
	}
	setNoFiles()
	gConfig.ctx, gConfig.cancel = context.WithCancel(context.Background())
}

type actionType int

const (
	actionNone actionType = iota
	actionDiscoNet
	actionDiscoHost
	actionIpHost
	actionRisk
	actionLocalInfo
	actionUnknown
)

type actionList map[actionType]string

var (
	myActionList = actionList{
		actionDiscoNet:  "net",
		actionDiscoHost: "host",
		actionIpHost:    "ip-host",
		actionRisk:      "risk",
		actionLocalInfo: "local-info",
	}
)

func (a actionList) String() string {
	r := make([]string, 0)
	for _, l := range myActionList {
		r = append(r, l)
	}
	return strings.Join(r, ",")
}

func (a actionList) Id(name string) actionType {
	for k, v := range myActionList {
		if v == name {
			return k
		}
	}
	return actionUnknown
}

func (a actionType) String() string {
	if n, ok := myActionList[a]; ok {
		return n
	}
	return "unknown"
}
