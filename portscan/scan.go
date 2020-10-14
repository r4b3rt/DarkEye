package main

import (
	"fmt"
	"github.com/zsdevX/DarkEye/common"
	"strconv"
	"time"
)

func New(ip string) *Scan {
	return &Scan{
		Ip:                   ip,
		DefaultTimeOut:       3000,
		ActivePort:           "80",
		MinTimeOut:           300,
		PortRange:            "80-8080,22,23",
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
