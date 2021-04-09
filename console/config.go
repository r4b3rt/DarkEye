package main

import (
	"context"
	"github.com/elastic/beats/libbeat/common/atomic"
)

//RequestContext add comment
type RequestContext struct {
	CmdArgs     []string
	ctx         context.Context
	cancel      context.CancelFunc
	running     atomic.Bool
	Interactive bool
}

type analysisEntity struct {
	ID      int64  `json:"id" gorm:"primaryKey"`
	Task    string `json:"task" gorm:"unique_index:UNIQ_hi;column:task"`
	Ip      string `json:"ip" gorm:"unique_index:UNIQ_hi;column:ip"`
	Port    string `json:"port" gorm:"unique_index:UNIQ_hi;column:port"`
	Service string `json:"service" gorm:"unique_index:UNIQ_hi;column:service"`

	Url             string `json:"url" gorm:"column:url"`
	Title           string `json:"title" gorm:"column:title"`
	WebServer       string `json:"web_server" gorm:"column:web_server"`
	WebResponseCode int32  `json:"http_code" gorm:"column:http_code"`

	Hostname  string
	Os        string
	Device    string
	Banner    string
	Version   string
	ExtraInfo string
	RDns      string
	Country   string

	NetBios     string `json:"netbios" gorm:"column:netbios"`
	WeakAccount string `json:"weak_account" gorm:"column:weak_account"`
	Vulnerable  string `json:"vulnerable" gorm:"column:vulnerable"`
}

func (analysisEntity) TableName() string {
	return "ent"
}
