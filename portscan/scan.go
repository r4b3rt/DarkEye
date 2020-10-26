package main

import (
	"fmt"
	"github.com/zsdevX/DarkEye/common"
	"strconv"
	"sync"
	"time"
)

var (
	PortList = "80,80-89,8000-9090,1433,1521,3306,5432,50000,443,445,873,5984,6379,7001,7002,9200,9300,11211,27017,27018,50000,50070,50030,21,22,23,445,2601,2604,3389,21,22,23,25,53,69,80,80-89,110,135,139,143,161,389,443,445,512,513,514,873,1025,111,1080,1158,1433,1521,2082,2083,2222,2601,2604,3128,3306,3312,3311,3389,3690,4440,4848,5432,5900,5984,6082,6379,7001,7002,7778,8000-9090,8080,8080,8089,9090,8081,8083,8649,8888,9000,9043,9200,9300,10000,11211,27017,27018,28017,50000,50060,50070,50030"
)

func New(ip string) *Scan {
	return &Scan{
		Ip:                 ip,
		DefaultTimeOut:     3000,
		ActivePort:         "80",
		MinTimeOut:         300,
		PortRange:          PortList,
		Test:               false,
		Title:              true,
		PortsScannedOpened: make([]PortInfo, 0),
		Callback:           callback,
		BarCallback:        barCallback,
	}
}

func (s *Scan) Run() {
	fromTo, _ := common.GetPortRange(s.PortRange)
	for _, p := range fromTo {
		if p.To-p.From > s.PortRangeThresholds {
			wg := sync.WaitGroup{}
			wg.Add(s.ThreadNumber)
			inc := (p.To - p.From) / s.ThreadNumber
			i := 0
			f := p.From
			for i < s.ThreadNumber {
				//端口扫描任务
				go func(a, b int) {
					defer wg.Done()
					s._run(a, b)
				}(f, f+inc)

				f += inc + 1
				i++
				if i == s.ThreadNumber-1 {
					inc = p.To - f
				}
			}
			wg.Wait()
		} else {
			s._run(p.From, p.To)
		}
	}
}

func (s *Scan) _run(from, to int) {
	for from <= to {
		//检查端口是否有效
		s.Check(from)
		if !s.IsFireWallNotForbidden() {
			//被防火墙策略限制探测，等待恢复期（恢复期比较傻，需要优化）。
			time.Sleep(time.Second * 10)
			//恢复后从中断的端口重新检测
			continue
		}
		from++
	}
}

func (s *Scan) Check(p int) {
	s.BarCallback(1)
	port := strconv.Itoa(int(p))
	if !common.IsAlive(s.Ip, port, s.TimeOut) {
		return
	}
	pi := PortInfo{}
	pi.Port = p
	if s.Title {
		pi.Server, pi.Title = common.GetHttpTitle("http", s.Ip+":"+port)
		if pi.Server == "" && pi.Title == "" {
			pi.Server, pi.Title = common.GetHttpTitle("https", s.Ip+":"+port)
		}
	}

	s.lock.Lock()
	defer s.lock.Unlock()
	s.PortsScannedOpened = append(s.PortsScannedOpened, pi)
	s.Callback(s.Ip, port, "Opened", pi.Server, pi.Title)
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
	fmt.Println(a...)
}

func barCallback(i int) {
	fmt.Println("Bar callback")
}
