package main

import (
	"encoding/csv"
	"fmt"
	"github.com/k0kubun/go-ansi"
	"github.com/schollz/progressbar"
	"github.com/zsdevX/DarkEye/common"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	Bar       = &progressbar.ProgressBar{}
	BarDesc   = make(chan *BarValue, 64)
	mFileSync = sync.RWMutex{}
)

type BarValue struct {
	Key   string
	Value string
}

func NewBar(max int) *progressbar.ProgressBar {
	bar := progressbar.NewOptions(max,
		progressbar.OptionSetWriter(ansi.NewAnsiStdout()),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionThrottle(3000*time.Millisecond),
		progressbar.OptionSetDescription("Loading ..."),
		progressbar.OptionSetWriter(os.Stderr),
		progressbar.OptionSetWidth(10),
		progressbar.OptionShowCount(),
		progressbar.OptionShowIts(),
		progressbar.OptionOnCompletion(func() {
			_, _ = fmt.Fprint(os.Stderr, "\nDONE")
		}),
		progressbar.OptionSpinnerType(14),
		progressbar.OptionFullWidth(),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}))

	_ = bar.RenderBlank()
	return bar
}

func New(ip string) *Scan {
	return &Scan{
		Ip: ip,
	}
}

func (s *Scan) Run() {
	fromTo, _ := genFromTo(s.PortRange)
	for _, p := range fromTo {
		i := p.from
		for i <= p.to {
			if _, ok := s.PortsHaveBeenScanned[i]; ok {
				_ = Bar.Add(1)
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
				fmt.Println("[", s.Ip, port, "Opened", "]", pi.Server, pi.Title)
				s.PortsScannedOpened = append(s.PortsScannedOpened, pi)
			}
			if !s.AliveTest() {
				time.Sleep(time.Second * 5)
				i++
				continue
			}
			_ = Bar.Add(1)
			s.PortsHaveBeenScanned[i] = true
			i++
		}
	}
}

func (s *Scan) OutPut(f *os.File) {
	mFileSync.Lock()
	defer mFileSync.Unlock()

	w := csv.NewWriter(f)

	for _, p := range s.PortsScannedOpened {
		_ = w.Write([]string{s.Ip, strconv.Itoa(p.Port), p.Server, p.Title})
	}
	w.Flush()
}

func genFromTo(portRange string) ([]FromTo, int) {
	res := make([]FromTo, 0)
	tot := 0
	ports := strings.Split(portRange, ",")
	for _, port := range ports {
		from := 0
		to := 0
		fromTo := strings.Split(port, "-")
		from, _ = strconv.Atoi(fromTo[0])
		to = from
		if len(fromTo) == 2 {
			to, _ = strconv.Atoi(fromTo[1])
		}
		a := FromTo{
			from: from,
			to:   to,
		}
		res = append(res, a)
		tot += 1 + to - from
	}
	return res, tot
}

func (s *Scan) AliveTest() bool {
	//不校正
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
	s.initParams(s.Ip)
	if s.PortsHaveBeenScanned == nil {
		s.PortsHaveBeenScanned = make(map[int]bool, 0)
	}
	if s.PortsScannedOpened == nil {
		s.PortsScannedOpened = make([]PortInfo, 0)
	}
	return nil
}

func (s *Scan) initParams(ip string) {
	s.Ip = ip
	s.DefaultTimeOut = *mTimeOut
	s.ActivePort = *mActivePort
	s.MinTimeOut = mMinTimeOut
	s.PortRange = *mPort
	s.Test = *mTestTimeOut
	s.Title = *mTitle
}
