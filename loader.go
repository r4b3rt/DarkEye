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
	switch c.action {
	case "disco-net":
		fallthrough
	case "disco-host":
		err = c.scanStart(scan.Discovery, scan.DiscoEnd, c.action == "disco-net")
	case "risk":
		err = c.scanStart(scan.RiskStart, scan.RiskEnd, false)
	case "localhost":
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
		my, err := c.scanInit(start)
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
