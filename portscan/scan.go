package main

import (
	"fmt"
	"github.com/zsdevX/DarkEye/common"
	"strconv"
	"time"
)

func New(ip string) *Scan {
	return &Scan{
		Ip: ip,
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
			port := strconv.Itoa(int(i))
			if common.IsAlive(s.Ip, port, s.TimeOut) {
				pi := PortInfo{}
				pi.Port = i
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
			if !s.AliveTest() {
				time.Sleep(time.Second * 5)
				i++
				continue
			}
			s.BarCallback()
			s.PortsHaveBeenScanned[i] = true
			i++
		}
	}
}

func (s *Scan) AliveTest() bool {
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
		if !s.AliveTest() {
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

func (s *Scan) InitConfig() error {
	s.DefaultTimeOut = *mTimeOut
	s.ActivePort = *mActivePort
	s.MinTimeOut = mMinTimeOut
	s.PortRange = *mPort
	s.Test = *mTestTimeOut
	s.Title = *mTitle

	if s.PortsHaveBeenScanned == nil {
		s.PortsHaveBeenScanned = make(map[int]bool, 0)
	}
	if s.PortsScannedOpened == nil {
		s.PortsScannedOpened = make([]PortInfo, 0)
	}
	return nil
}
