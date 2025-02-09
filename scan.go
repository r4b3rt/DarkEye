package main

import (
	"context"
	"fmt"
	"github.com/b1gcat/DarkEye/scan"
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
		if c._scanJustCount(sc, len(port)) {
			continue
		}
		if quit.Load().(bool) {
			break
		}
		sc.p.Add(1)
		go func(tip net.IP) {
			defer sc.p.Done()
			c._scanningOne(ctx, tip, port, sc,
				func(l interface{}) {
					switch sc.action {
					case actionDiscoNet:
						once.Do(func() {
							c.output(sc.sid.String(), ipc, "alive")
							quit.Store(true)
							cancel()
						})
						return
					default:
						c.output(l)
						return //continue for disco
					}
				})
		}(ip)
	}
}

func (c *config) _scanningOne(ctx context.Context, f net.IP, port []string,
	sc *myScan, cb func(interface{})) {

	switch sc.sid {
	case scan.DiscoNb:
		fallthrough
	case scan.DiscoPing:
		defer sc.bar.Inc()
		r, err := sc.s.Start(ctx, f.String(), "0")
		if err != nil || r == nil {
			logrus.Debug("_scanningOne.not.found:", err)
			return
		}
		//found
		cb(r)
		return
	}

	pp := EzPool(c.maxThreadForEachIPScan)
	defer pp.Wait()

	for _, _p := range port {
		select {
		case <-ctx.Done():
			pp.Close()
			return
		default:
		}
		pp.Add(1)
		go func(p string) {
			defer func() {
				sc.bar.Inc()
				pp.Done()
			}()
			switch {
			case sc.sid >= scan.DiscoEnd:
				if !sc.s.Identify(ctx, f.String(), p) {
					logrus.Debug("_scanningOne.ident.fail:", f.String(), ":", p)
					return
				}
				fallthrough
			default:
				r, err := sc.s.Start(ctx, f.String(), p)
				if err != nil || r == nil {
					logrus.Debug("_scanningOne.not.found:", err)
					return
				}
				//found callback
				cb(r)
			}
		}(_p)
	}
}

func (c *config) _scanJustCount(sc *myScan, ports int) bool {
	if sc.action != actionNone {
		return false
	}
	switch {
	case sc.sid == scan.DiscoPing || sc.sid == scan.DiscoNb:
		sc.total++ //ip
	default:
		sc.total += int64(ports)
	}
	return true
}
