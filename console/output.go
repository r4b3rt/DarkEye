package main

import (
	"github.com/gocarina/gocsv"
	"github.com/olekukonko/tablewriter"
	"github.com/zsdevX/DarkEye/common"
	"os"
	"sort"
)

func (s *superScanRuntime) OutPut() {
	//排序
	sort.Slice(s.result, func(i, j int) bool {
		return s.result[j].Ip > s.result[i].Ip
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

	t, _ := tablewriter.NewCSV(os.Stdout, s.Output, true)
	t.SetAlignment(tablewriter.ALIGN_LEFT) // Set Alignment
	t.Render()
}
