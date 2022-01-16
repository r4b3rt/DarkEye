package main

import (
	"fmt"
	"github.com/b1gcat/DarkEye/scan"
	"github.com/sirupsen/logrus"
	"strings"
	"sync"
)

func (c *config) loader() error {
	logrus.Info("start action:", c.action)
	var err error
	defer logrus.Info("stop")
	discoNet := false
	switch myActionList.Id(c.action) {
	case actionDiscoNet:
		discoNet = true
		fallthrough
	case actionDiscoHost:
		err = c.scanStart(scan.Discovery, scan.DiscoEnd, discoNet)
	case actionRisk:
		err = c.scanStart(scan.RiskStart, scan.RiskEnd, discoNet)
	case actionLocalInfo:
	default:
		err = fmt.Errorf("not support action %v", c.action)
	}

	return err
}

func (c *config) scanStart(start, end int, discoNet bool) error {
	loaders, err := c.readLoaders()
	if err != nil {
		return err
	}

	wg := sync.WaitGroup{}
	for start < end {
		start++
		if _, ok := loaders[start]; !ok {
			logrus.Info("ignore:", scan.IdList.Name(start))
			continue
		}
		var my *myScan
		switch {
		case start < scan.DiscoEnd:
			my, err = c.scanInit(start, scan.IdList.Name(start))
		case start > scan.RiskStart && start < scan.RiskEnd:
			my, err = c.scanInit(start)
		default:
			return fmt.Errorf("unknown scan id %v", start)
		}
		if err != nil {
			return err
		}
		my.discoNet = discoNet
		my.sid = start
		wg.Add(1)
		go func(sc *myScan) {
			defer wg.Done()
			c.run(my)
		}(my)
	}
	wg.Wait()
	return nil
}

func (c *config) scanInit(sid int, args ...interface{}) (*myScan, error) {
	var err error

	my := &myScan{
		p: EzPool(c.maxThreadForEach),
	}
	my.s, err = scan.New(sid, c.timeout, args)
	return my, err
}

func (c *config) readLoaders() (map[int]string, error) {
	r := make(map[int]string, 0)
	loaders := strings.Split(c.loaders, ",")
	if loaders[0] == "all" {
		loaders = strings.Split(scan.IdList.String(), ",")
	}
	for _, l := range loaders {
		id := scan.IdList.Id(l)
		if id == scan.Unknown {
			return nil, fmt.Errorf("unkown loader %v", id)
		}
		if k, ok := r[id]; ok {
			logrus.Warn("overwrite ", k, "=>", l)
		}
		r[id] = l
	}
	return r, nil
}
