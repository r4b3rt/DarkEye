package main

import (
	"fmt"
	"github.com/k0kubun/go-ansi"
	"github.com/schollz/progressbar"
	"github.com/zsdevX/DarkEye/common"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var (
	Bar     = &progressbar.ProgressBar{}
	BarDesc = make(chan *BarValue, 64)
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

func BarDescription() {
	description := make(map[string]string, 0)
	for {
		desc := ""
		i := 0
		value := <-BarDesc
		description[value.Key] = value.Value
		for k, v := range description {
			i++
			desc += fmt.Sprintf("%d.%s:%s\n", i, k, v)
		}
		Bar.Describe(desc)
	}
}

func New(ip string) *Scan {
	return &Scan{
		Ip: ip,
	}
}

func (s *Scan) Run() {
	fromTo, tot := genFromTo(s.PortRange)
	find := 0
	blk := 0
	for _, p := range fromTo {
		i := p.from
		for i <= p.to {
			if len(s.PortsScannedOpened)  != 0 {
				BarDesc <- &BarValue{
					Value: fmt.Sprintf("[+%d/-%d/%d/%d/%dms] opened:[%s]",
						find, blk, tot, i, s.Rate, genPortList(s.PortsScannedOpened)),
					Key: s.Ip,
				}
			}
			if _, ok := s.PortsHaveBeenScanned[i]; ok {
				_ = Bar.Add(1)
				continue
			}
			port := strconv.Itoa(int(i))
			if common.IsAlive(s.Ip, port, s.Rate) {
				s.PortsScannedOpened = append(s.PortsScannedOpened, i)
				find++
				_ = s.saveCfg()
			}
			if !s.AliveTest() {
				blk++
				time.Sleep(time.Second * 5)
				i++
				continue
			}
			_ = Bar.Add(1)
			s.PortsHaveBeenScanned[i] = true
			i++
			_ = s.saveCfg()
		}
	}
	_ = s.saveCfg()
}

func genPortList(ports []int) string {
	portList := ""
	for k, p := range ports {
		if k != 0 {
			portList += " "
		}
		portList += strconv.Itoa(p)
	}
	return portList
}
func (s *Scan) OutPut() {
	filename := filepath.Join(mBasedir, s.Ip+"."+*mOutputFile)
	if len(s.PortsScannedOpened) == 0 {
		_ = common.SaveFile(fmt.Sprintf("%s开放端口:无", s.Ip), filename)
	} else {
		for _, p := range s.PortsScannedOpened {
			_ = common.SaveFile(fmt.Sprintf("%s开放端口:%d", s.Ip, p), filename)
		}
	}
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
		if common.IsAlive(s.Ip, s.ActivePort, s.Rate) {
			return true
		}
		maxRetries --
	}
	return false
}

func (s *Scan) RateTest() error {
	s.Rate = s.DefaultRate
	if !s.Test {
		return nil
	}
	lastRate := 0
	for {
		if !s.AliveTest() {
			if lastRate == 0 {
				return fmt.Errorf("网络质量差,默认rate太低,烦请大佬手动调高rate参数:-rate 3000ms或放弃目标")
			}
			s.Rate = lastRate
			break
		} else {
			fmt.Println(fmt.Sprintf("[OK] Tested Rate: %dms", s.Rate))
			if s.Rate <= s.MinRate {
				break
			}
			lastRate = s.Rate
			s.Rate -= 50
			time.Sleep(time.Millisecond * time.Duration(s.Rate))
		}
	}
	return nil
}

func (s *Scan) InitConfig() error {
	_ = s.loadCfg()
	if !s.valid {
		s.initParams(s.Ip)
	}
	if s.PortsHaveBeenScanned == nil {
		s.PortsHaveBeenScanned = make(map[int]bool, 0)
	}
	if s.PortsScannedOpened == nil {
		s.PortsScannedOpened = make([]int, 0)
	}
	return nil
}

func (s *Scan) initParams(ip string) {
	s.Ip = ip
	s.DefaultRate = *mRate
	s.ActivePort = *mActivePort
	s.MinRate = *mMinRate
	s.PortRange = *mPort
	s.Test = *mTestRate
	s.PortsScannedOpened = nil
	s.PortsHaveBeenScanned = nil
	_ = s.saveCfg()
}
