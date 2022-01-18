package main

import (
	"github.com/vbauerster/mpb/v7"
	"github.com/vbauerster/mpb/v7/decor"
	"sync"
	"time"
)

type progress struct {
	enabled  bool
	progress *mpb.Progress
}

type bar struct {
	sync.RWMutex
	p   *progress
	bar *mpb.Bar
}

func (p *progress) init(enable bool) {
	if !enable {
		return
	}
	p.enabled = enable
	p.progress = mpb.New()
}

func (p *progress) Add(name string, total int64) *bar {
	mb := bar{
		p: p,
	}
	if !p.enabled {
		return &mb
	}
	mb.bar = p.progress.AddBar(total,
		mpb.PrependDecorators(
			decor.Name(name, decor.WC{W: len(name) + 1, C: decor.DidentRight}),
			decor.CountersNoUnit("%d/%d", decor.WCSyncWidth),
		),
		mpb.AppendDecorators(decor.Percentage(decor.WC{W: 5})),
	)
	return &mb
}

func (b *bar) Inc() {
	if !b.p.enabled {
		return
	}
	b.Lock()
	defer b.Unlock()
	b.bar.Increment()
	b.bar.DecoratorEwmaUpdate(time.Since(time.Now()))
}
