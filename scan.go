package main

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"net"
	"sync"
	"sync/atomic"
)

func (c *config) run(sc *myScan) {
	ips := readList(c.ip)
	for _, ip := range ips {
		if err := c.scanning(ip, sc); err != nil {
			logrus.Error("scanning:", err.Error())
			return
		}
	}
}

func (c *config) scanning(ip string, sc *myScan) error {
	ips, err := splitIp2C(ip)
	if err != nil {
		return err
	}

	port := portSplit(c.port)
	if port == nil {
		return fmt.Errorf("no port found")
	}

	//Each C
	for _, ipc := range ips {
		ipCs, err := splitIpC2Ip(ipc) //Each ip in C
		if err != nil {
			return err
		}
		c._scanning(ipc, ipCs, port, sc)
	}
	sc.p.Wait()
	return nil
}

func (c *config) _scanning(ipc string, ipCs []net.IP, port []string, sc *myScan) {
	ctx, cancel := context.WithCancel(c.ctx)
	once := sync.Once{}
	var quit atomic.Value

	quit.Store(false)
	for _, ip := range ipCs {
		if quit.Load().(bool) {
			break
		}
		sc.p.Add(1)
		go func(tip net.IP) {
			defer sc.p.Done()
			c._scanningOne(ctx, tip, port, sc,
				func(l interface{}) (stop bool) {
					//!discovery
					if sc.disco == "" {
						c.output(l)
						stop = true //stop if found one
						return
					}
					//discovery host or Net
					if sc.discoNet {
						stop = true //stop if found one ip in C
						once.Do(func() {
							c.output(ipc, "alive")
							quit.Store(true)
							cancel()
						})
						return
					}
					c.output(l)
					return //continue for disco
				})
		}(ip)
	}
}

func (c *config) _scanningOne(ctx context.Context, f net.IP, port []string,
	sc *myScan, cb func(interface{}) bool) {

	if sc.disco == "ping" {
		r, err := sc.s.Start(ctx, f.String(), "0")
		if err != nil || r == nil {
			logrus.Debug("_scanningOne.not.found:", err)
			return
		}
		//found
		cb(r)
		return
	}

	for _, p := range port {
		if sc.disco == "" {
			if !sc.s.Identify(ctx, f.String(), p) {
				logrus.Debug("_scanningOne.ident.fail:", f.String(), ":", p)
				return
			}
		}
		r, err := sc.s.Start(ctx, f.String(), p)
		if err != nil || r == nil {
			logrus.Debug("_scanningOne.not.found:", err)
			continue
		}
		//found
		if cb(r) {
			return
		}
	}
}
