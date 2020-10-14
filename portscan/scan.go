package main

import (
	"fmt"
	"github.com/zsdevX/DarkEye/common"
	"strconv"
	"time"
)

var (
	PortList = "80,80-89,8000-9090,1433,1521,3306,5432,50000,443,445,873,5984,6379,7001,7002,9200,9300,11211,27017,27018,50000,50070,50030,21,22,23,445,2601,2604,3389,21,22,23,25,53,69,80,80-89,110,135,139,143,161,389,443,445,512,513,514,873,1025,111,1080,1158,1433,1521,2082,2083,2222,2601,2604,3128,3306,3312,3311,3389,3690,4440,4848,5432,5900,5984,6082,6379,7001,7002,7778,8000-9090,8080,8080,8089,9090,8081,8083,8649,8888,9000,9043,9200,9300,10000,11211,27017,27018,28017,50000,50060,50070,50030"
)

func New(ip string) *Scan {
	return &Scan{
		Ip:                   ip,
		DefaultTimeOut:       3000,
		ActivePort:           "80",
		MinTimeOut:           300,
		PortRange:            PortList,
		Test:                 false,
		Title:                true,
		PortsHaveBeenScanned: make(map[int]bool, 0),
		PortsScannedOpened:   make([]PortInfo, 0),
		Callback:             callback,
		BarCallback:          barCallback,
	}
}

func (s *Scan) Run() {
	fromTo, _ := common.GetPortRange(s.PortRange)
	for _, p := range fromTo {
		i := p.From
		for i <= p.To {
			if _, ok := s.PortsHaveBeenScanned[i]; ok {
				s.BarCallback()
				continue
			}
			//检查端口是否有效
			s.Check(i)
			if !s.IsFireWallNotForbidden() {
				//被防火墙策略限制探测，等待恢复期（恢复期比较傻，需要优化）。
				time.Sleep(time.Second * 10)
				//恢复后从中断的端口重新检测
				continue
			}
			s.BarCallback()
			s.PortsHaveBeenScanned[i] = true
			i++
		}
	}
}

func (s *Scan) Check(p int) {
	port := strconv.Itoa(int(p))
	if common.IsAlive(s.Ip, port, s.TimeOut) {
		pi := PortInfo{}
		pi.Port = p
		if s.Title {
			pi.Server, pi.Title = common.GetHttpTitle("http", s.Ip+":"+port)
			if pi.Server == "" && pi.Title == "" {
				pi.Server, pi.Title = common.GetHttpTitle("https", s.Ip+":"+port)
			}
		}
		if s.Callback != nil {
			s.Callback(s.Ip, port, "Opened", pi.Server, pi.Title)
		}
		s.PortsScannedOpened = append(s.PortsScannedOpened, pi)
	}
}

func (s *Scan) IsFireWallNotForbidden() bool {
	//为0不矫正
	if s.ActivePort == "0" {
		return true
	}
	maxRetries := 3
	for maxRetries > 0 {
		if common.IsAlive(s.Ip, s.ActivePort, s.TimeOut) {
			return true
		}
		maxRetries --
	}
	return false
}

func (s *Scan) TimeOutTest() error {
	s.TimeOut = s.DefaultTimeOut
	if !s.Test {
		return nil
	}
	lastRate := 0
	for {
		if !s.IsFireWallNotForbidden() {
			if lastRate == 0 {
				return fmt.Errorf("网络质量差,默认timeout太低,默认timeout太低:-timeout 3000ms或放弃目标")
			}
			s.TimeOut = lastRate
			break
		} else {
			fmt.Println(fmt.Sprintf("[OK] 测试超时: %dms", s.TimeOut))
			if s.TimeOut <= s.MinTimeOut {
				break
			}
			lastRate = s.MinTimeOut
			s.TimeOut -= 50
			time.Sleep(time.Millisecond * time.Duration(s.TimeOut))
		}
	}
	return nil
}

func callback(a ...interface{}) {
	fmt.Println(a)
}

func barCallback() {
	fmt.Println("Bar callback")
}
