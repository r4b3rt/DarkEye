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
	taskId      int
}

type analysisEntity struct {
	ID      int64  `json:"id,omitempty" gorm:"primaryKey"`
	Task    string `json:"task,omitempty" gorm:"unique_index:UNIQ_hi;column:task"`
	Ip      string `json:"ip,omitempty" gorm:"unique_index:UNIQ_hi;column:ip"`
	Port    string `json:"port,omitempty" gorm:"unique_index:UNIQ_hi;column:port"`
	Service string `json:"service,omitempty" gorm:"unique_index:UNIQ_hi;column:service"`

	Url             string `json:"url,omitempty" gorm:"column:url"`
	Title           string `json:"title,omitempty" gorm:"column:title"`
	WebServer       string `json:"web_server,omitempty" gorm:"column:web_server"`
	WebResponseCode int32  `json:"http_code,omitempty" gorm:"column:http_code"`

	Hostname  string `json:"hostname,omitempty"`
	Os        string `json:"os,omitempty"`
	Device    string `json:"device,omitempty"`
	Banner    string `json:"banner,omitempty"`
	Version   string `json:"version,omitempty"`
	ExtraInfo string `json:"extra_info,omitempty" gorm:"column:extra_info"`
	RDns      string `json:"r_dns,omitempty" gorm:"column:r_dns"`
	Country   string `json:"country,omitempty"`
	Isp       string `json:"isp,omitempty"`

	NetBios     string `json:"netbios,omitempty" gorm:"column:netbios"`
	WeakAccount string `json:"weak_account,omitempty" gorm:"column:weak_account"`
	ExpHelper   string `json:"exp_helper,omitempty" gorm:"column:exp_helper"`
}

type crawler struct {
	ID     int64  `json:"id,omitempty" gorm:"primaryKey"`
	Target string `json:"target,omitempty"`
	Url    string `json:"url,omitempty"`
	Method string `json:"method,omitempty"`
	Data   string `json:"data,omitempty"`
}

func (analysisEntity) TableName() string {
	return "ent"
}

func (crawler) TableName() string {
	return "crawler"
}
