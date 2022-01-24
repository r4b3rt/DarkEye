package main

import (
	"flag"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"time"
)

type config struct {
	//ioc
	iocFile     string //ip or domain
	ioc         bool
	iocInterval int

	//gather
	noGather bool

	verbose bool
	output  string
	logger  *logrus.Logger
}

var (
	Version = "v1.0.0"
	gConfig = &config{
		logger: logrus.New(),
	}
)

func main() {
	fmt.Println(Version)
	gConfig.initialize()
	gConfig.gatherInfo()
	gConfig.iOC()

	logrus.Info("CTRL+C to quit")
	for {
		time.Sleep(time.Second * 60)
	}
}

func (c *config) initialize() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		c.logger.Error("filepath.Abs:", err.Error())
	}
	flag.StringVar(&c.iocFile, "ioc-file", filepath.Join(dir, "ioc.txt"), "情报ip或域名")
	flag.IntVar(&c.iocInterval, "ioc-interval", 5, "ioc检查周期")
	flag.BoolVar(&c.ioc, "ioc", false, "开启ioc检查")

	flag.BoolVar(&c.noGather, "no-gather", false, "开启收集信息")

	flag.StringVar(&c.output, "output", filepath.Join(dir, "output"), "输出目录")
	flag.BoolVar(&c.verbose, "v", false, "开启调试")
	flag.Parse()

	if c.verbose {
		c.logger.SetLevel(logrus.DebugLevel)
	}
	c.logger.SetFormatter(&logrus.TextFormatter{
		DisableColors:          false,
		TimestampFormat:        "200612150405",
		DisableLevelTruncation: false,
	})

	os.Mkdir(c.output, 0600)
	c.logger.Info("Writing data to:", c.output)
}
