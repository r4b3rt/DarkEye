package main

import (
	"flag"
	"fmt"
	"github.com/schollz/progressbar"
	"github.com/zsdevX/DarkEye/common"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var (
	mPort       = flag.String("port", "1-65535", "端口格式参考Nmap")
	mIp         = flag.String("ip", "127.0.0.1", "a.b.c.d（不做扫C，扫C自己想办法或使用nmap --scan-delay 1000ms, 不准")
	mActivePort = flag.String("alive_port", "80", "已知开放的端口用来校正扫描")
	mSpeed      = flag.Int("speed", 2000, "端口之间的扫描间隔单位ms，也可用通过-test_speed自动计算")
	mMinSpeed   = flag.Int("min_speed", 100, "自动计算的速率不能低于min_speed")
	mTestSpeed  = flag.Bool("speed_test", false, "检测防火墙限制频率")
	mOutputFile = flag.String("output", "result.txt", "结果保存到该文件")
	mHowTo = flag.Bool("examples", false, "显示使用示例")
)

var (
	scanCfg = Scan{}
)

func main() {
	flag.Parse()

	if *mHowTo {
		fmt.Println("./lowSpeedPortScan -alive_port 8443 -ip f.u.c.k -port 1-65535 -speed_test -output result.txt")
		return
	}

	if err := initConfig(); err != nil {
		fmt.Println(err.Error())
		os.Exit(0)
	}
	s := &scanCfg
	s.SpeedTest()
	s.Run()
	s.OutPut()
	fmt.Println("完成:", filepath.Join(common.BaseDir, *mOutputFile))
}

func (s *Scan) Run() {
	fromTo, tot := genFromTo(s.PortRange)

	fmt.Println(fmt.Sprintf("开启检测共计%d个端口", tot))

	bar := progressbar.Default(int64(tot), fmt.Sprintf("[0/0/%d] scanning", tot))
	find := 0
	blk := 0
	for _, p := range fromTo {
		i := p.from
		for i <= p.to {
			if _, ok := s.PortsHaveBeenScanned[i]; ok {
				_ = bar.Add(1)
				continue
			}
			port := strconv.Itoa(int(i))
			if common.IsAlive(s.Ip, port, s.Speed) {
				s.PortsScannedOpened = append(s.PortsScannedOpened, i)
				find++
				bar.Describe(fmt.Sprintf("[%d/%d/%d] opened:[%s]", find, blk, tot, genPortList(s.PortsScannedOpened)))
				_ = saveCfg()
			}
			if !s.AliveTest() {
				blk++
				bar.Describe(fmt.Sprintf("[%d/%d/%d] opened:[%s]", find, blk, tot, genPortList(s.PortsScannedOpened)))
				time.Sleep(time.Second * 5)
				continue
			}
			_ = bar.Add(1)
			s.PortsHaveBeenScanned[i] = true
			i++
			_ = saveCfg()
		}
	}
	_ = saveCfg()
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
	filename := filepath.Join(common.BaseDir, *mOutputFile)
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
	maxRetries := 3
	for maxRetries > 0 {
		if common.IsAlive(s.Ip, s.ActivePort, s.DefaultSpeed) {
			return true
		}
		maxRetries --
	}
	return false
}

func (s *Scan) SpeedTest() error {
	s.Speed = s.DefaultSpeed
	if !s.Test {
		s.Speed = s.DefaultSpeed
		return nil
	}
	lastSpeed := 0
	for {
		if !s.AliveTest() {
			if lastSpeed == 0 {
				return fmt.Errorf("网络质量差,默认Speed太低,烦请大佬手动调高speed参数:-speed 3000ms或放弃目标")
			}
			s.Speed = lastSpeed
			break
		} else {
			fmt.Println(fmt.Sprintf("[OK] Tested Speed: %dms", s.Speed))
			if s.Speed <= s.MinSpeed {
				break
			}
			lastSpeed = s.Speed
			s.Speed -= 50
			time.Sleep(time.Millisecond * time.Duration(s.Speed))
		}
	}
	fmt.Println(fmt.Sprintf("测试频率为:%dms", s.Speed))
	return nil
}

func initConfig() error {
	_ = loadCfg()
	if scanCfg.valid {
		ans := ""
		fmt.Printf("检测到已经保存上一次的扫描记录,是否继续使用（yes/no）: ")
		n, _ := fmt.Scanln(&ans)
		if n != 1 {
			return fmt.Errorf("输入错误")
		}
		if ans == "no" {
			initParams()
		} else if ans == "yes" {
		} else {
			return fmt.Errorf("输入错误")
		}
	} else {
		initParams()
	}
	if scanCfg.PortsHaveBeenScanned == nil {
		scanCfg.PortsHaveBeenScanned = make(map[int]bool, 0)
	}
	if scanCfg.PortsScannedOpened == nil {
		scanCfg.PortsScannedOpened = make([]int, 0)
	}
	return nil
}

func initParams() {
	scanCfg.DefaultSpeed = *mSpeed
	scanCfg.ActivePort = *mActivePort
	scanCfg.Ip = *mIp
	scanCfg.MinSpeed = *mMinSpeed
	scanCfg.PortRange = *mPort
	scanCfg.Test = *mTestSpeed
	_ = saveCfg()
}
