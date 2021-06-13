package main

import (
	"github.com/b1gcat/DarkEye/common"
	"github.com/gocarina/gocsv"
	"github.com/olekukonko/tablewriter"
	"golang.org/x/time/rate"
	"os"
	"sort"
	"sync"
	"time"
)

var (
	outputRate = rate.NewLimiter(rate.Every(3*time.Second), 1)
	outputLk   sync.RWMutex
)

func (s *superScanRuntime) OutPut() {
	//排序
	outputLk.Lock()
	defer outputLk.Unlock()
	sort.Slice(s.result, func(i, j int) bool {
		return common.CompareIP(s.result[j].Ip, s.result[i].Ip) > 0
	})
	//写入文件
	csvFile, err := os.OpenFile(s.Output, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		common.Log(superScan, err.Error(), common.FAULT)
		return
	}
	defer csvFile.Close()

	if err := gocsv.Marshal(s.result, csvFile); err != nil { // Load clients from file
		common.Log(superScan, err.Error(), common.FAULT)
		return
	}

	if outputRate.Allow() {
		s.display()
	}
}

func (s *superScanRuntime) display() {
	t, err := tablewriter.NewCSV(os.Stdout, s.Output, true)
	if err != nil {
		return
	}
	t.SetAlignment(tablewriter.ALIGN_LEFT) // Set Alignment
	t.Render()
}
