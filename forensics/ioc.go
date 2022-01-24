package main

import (
	"encoding/json"
	"fmt"
	"github.com/b1gcat/go-libraries/utils"
	"github.com/shirou/gopsutil/process"
	"github.com/sirupsen/logrus"
	"net"
	"os/exec"
	"path/filepath"
	"time"
)

var (
	iocOutput = "output_ioc.txt"
)

type iocResult struct {
	Time    string `json:"time"`
	Pid     int32  `json:"pid"`
	Ppid    int32  `json:"ppid"`
	Name    string `json:"name"`
	Dir     string `json:"dir"`
	Cmdline string `json:"cmdline"`
	Ioc     string `json:"ioc"`
	IocIp   string `json:"ioc_ip,omitempty"`
}

func (c *config) iOC() {
	if !gConfig.ioc {
		c.logger.Warn("ioc disabled")
		return
	}
	c.logger.Info("start ioc checking...")
	defer c.logger.Info("ioc done!")

	iocList, err := utils.LoadListFromFile(c.iocFile)
	if err != nil {
		c.logger.Error("iOC.", err.Error())
		return
	}

	ioc := make(map[string]string, 0)
	for _, v := range iocList {
		switch {
		case v == "":
			break
		case net.ParseIP(v) == nil:
			ips, err := c.domain2ip(v)
			if err != nil {
				c.logger.Error("iOC.", err.Error())
				break
			}
			for _, ip := range ips {
				ioc[ip] = v
			}
		default:
			ioc[v] = v
		}
	}
	c.logger.Debug("ioc:", ioc)
	loop := 0
	tk := time.NewTicker(time.Second * time.Duration(c.iocInterval))
	c.logger.Info("loading ioc ip number:", len(ioc))
	for {
		select {
		case <-tk.C:
			loop++
			if err = c.iocRun(ioc, loop); err != nil {
				logrus.Error("iOC.", err.Error())
				return
			}
		}
	}
}

func (c *config) iocRun(ioc map[string]string, loop int) error {
	c.logger.Debug("iot run ", loop, " times")
	ps, err := process.Processes()
	if err != nil {
		return fmt.Errorf("iocRun.Processes:%v", err.Error())
	}
	for _, p := range ps {
		ns, err := p.Connections()
		if err != nil {
			return fmt.Errorf("iocRun.Connections:%v", err.Error())
		}
		for _, n := range ns {
			id := n.Laddr.IP
			if _, ok := ioc[id]; !ok {
				id = n.Raddr.IP
				if _, ok = ioc[id]; !ok {
					continue
				}
			}

			cmdline, _ := p.Cmdline()
			name, _ := p.Name()
			dir, _ := p.Cwd()
			ppid, _ := p.Ppid()
			if dir == "" {
				dir, _ = exec.LookPath(name)
			}
			r := &iocResult{
				Time:    time.Now().Format("2006/1/2 15:04:05"),
				Pid:     p.Pid,
				Ppid:    ppid,
				Name:    name,
				Cmdline: cmdline,
				Dir:     dir,
				Ioc:     ioc[id],
			}
			if net.ParseIP(r.Ioc) == nil {
				r.IocIp = id
			}
			c.logger.Info(r)
			body, err := json.MarshalIndent(r, "", "	")
			if err != nil {
				return fmt.Errorf("output.Marshal:%v", err.Error())
			}
			utils.WriteToFile(filepath.Join(c.output, iocOutput), string(body))
		}
	}
	return nil
}

func (c *config) domain2ip(domain string) ([]string, error) {
	ips, err := net.LookupHost(domain)
	if err != nil {
		return nil, fmt.Errorf("domain2ip:%v", err.Error())
	}
	logrus.Info("dns resolver:", domain, " ", ips)
	return ips, nil
}
