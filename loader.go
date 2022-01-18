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

	loaders, err := c.readLoaders()
	if err != nil {
		return err
	}

	c.progress.init(c.bar)
	wg := sync.WaitGroup{}
	for _, loader := range loaders {
		var my *myScan
		my, err = c.scanInit(loader)
		if err != nil {
			return err
		}
		if c.bar {
			my.sid = scan.Nothing
			c.run(my) //calc total ip
		}
		my.bar = c.progress.Add(loader.String(), my.total)
		my.total = 0
		my.action = myActionList.Id(c.action)
		my.sid = loader
		wg.Add(1)
		go func(sc *myScan) {
			defer wg.Done()
			c.run(my)
		}(my)
	}
	wg.Wait()

	return err
}

func (c *config) scanInit(sid scan.IdType) (*myScan, error) {
	var err error

	my := &myScan{
		p:  EzPool(c.maxThreadForEachScan),
	}
	my.s, err = scan.New(sid, c.timeout)
	l := logrus.New()
	l.SetLevel(logrus.GetLevel())
	switch {
	case sid > scan.DiscoHttp && c.action == actionIpHost.String():
		my.s.Setup(l, readList(c.host))
	case sid > scan.RiskStart && sid <= scan.RiskEnd && (c.user != "" || c.pass != ""):
		my.s.Setup(l, readList(c.user), readList(c.pass))
	default:
		my.s.Setup(l)
	}

	return my, err
}

func (c *config) readLoaders() ([]scan.IdType, error) {
	r := make([]scan.IdType, 0)
	loaders := strings.Split(c.loaders, ",")
	for _, l := range loaders {
		id := scan.IdList.Id(l)
		if id == scan.Unknown {
			return nil, fmt.Errorf("unkown loader %v", id)
		}
		switch c.action {
		case actionDiscoNet.String():
			if id == scan.DiscoNb || id == scan.DiscoHttp {
				logrus.Info("ignoring scan:", id.String())
				continue
			}
			fallthrough
		case actionDiscoHost.String():
			if id > scan.DiscoEnd {
				logrus.Info("ignoring scan:", id.String())
				continue
			}
		case actionRisk.String():
			if id <= scan.RiskStart || id >= scan.RiskEnd {
				logrus.Info("ignoring scan:", id.String())
				continue
			}
		case actionIpHost.String():
			r = append(r, scan.DiscoHttp)
			logrus.Info("force enabled:", scan.DiscoHttp.String())
			return r, nil
		case actionLocalInfo.String():
			fallthrough
		default:
			return nil, fmt.Errorf("not support action:%v", c.action)
		}

		r = append(r, id)
	}
	for _, l := range r {
		logrus.Info("enabled:", l.String())
	}
	return r, nil
}
