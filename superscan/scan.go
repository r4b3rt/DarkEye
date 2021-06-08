package superscan

import (
	"fmt"
	"github.com/orcaman/concurrent-map"
	"github.com/zsdevX/DarkEye/common"
	"github.com/zsdevX/DarkEye/superscan/plugins"
	"strconv"
	"time"
)

//Run add comment
func (s *Scan) Run() {
	task := common.NewTask(s.Thread, s.Parent)
	defer task.Wait("scan")

	task.Job()
	go func() {
		defer task.UnJob()
		s.preCheck()
	}()
	fromTo, _ := common.GetPortRange(s.PortRange)
	for _, p := range fromTo {
		for p.From <= p.To {
			//Task Terminate
			if !task.Job() {
				return
			}
			go func(port int) {
				defer task.UnJob()
				s.job(port)
			}(p.From)

			p.From++
		}
	}
}

func (s *Scan) job(port int) {
	for {
		s.Check(port)
		if !s.isFireWallNotForbidden() {
			//被防火墙策略限制探测，等待恢复期（恢复期比较傻，需要优化）。
			time.Sleep(time.Second * 10)
			//恢复后从中断的端口重新检测
			continue
		}
		break
	}
}

//Check add comment
func (s *Scan) Check(p int) {
	defer func() {
		s.BarCallback(1)
	}()
	if s.ActivePort != "0" {
		if common.IsAlive(s.Parent, s.Ip, strconv.Itoa(p), s.TimeOut) == common.Alive {
			common.Log("with.ActivePort", s.Ip+":"+strconv.Itoa(p), common.INFO)
		}
		//开启防火墙检测仅判断端口，不爆破
		return
	}
	plg := plugins.Plugins{
		TargetIp:   s.Ip,
		TargetPort: strconv.Itoa(p),
		Result: plugins.Result{
			Output: cmap.New(),
		},
	}
	plg.Check()
	s.report(&plg)
}

func (s *Scan) preCheck() {
	if s.ActivePort != "0" {
		//开启防火墙检测仅判断端口，不探测
		return
	}
	plg := plugins.Plugins{
		TargetIp: s.Ip,
		Result: plugins.Result{
			Output: cmap.New(),
		},
	}
	plg.PreCheck()
	s.report(&plg)
}

func (s *Scan) report(plg *plugins.Plugins) {
	if !plg.Result.PortOpened {
		return
	}
	s.Callback(plg)
}

func (s *Scan) isFireWallNotForbidden() bool {
	//为0不矫正
	if s.ActivePort == "0" {
		return true
	}
	maxRetries := 3
	for maxRetries > 0 {
		if common.IsAlive(s.Parent, s.Ip, s.ActivePort, s.TimeOut) == common.Alive {
			return true
		}
		maxRetries--
	}
	return false
}

//New add comment
func New(ip string) *Scan {
	return &Scan{
		Ip:          ip,
		ActivePort:  "80",
		Thread:      1,
		PortRange:   common.PortList,
		Callback:    func(v interface{}) { fmt.Println(v) },
		BarCallback: func(int) {},
	}
}
