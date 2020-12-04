package plugins

import (
	"golang.org/x/time/rate"
	"sync"
)

type Web struct {
	Server string `json:",omitempty"`
	Title  string `json:",omitempty"`
	Code   int32  `json:",omitempty"`
}

type Account struct {
	Username string `json:",omitempty"`
	Password string `json:",omitempty"`
}

type NetBios struct {
	HostName  string `json:",omitempty"`
	UserName  string `json:",omitempty"`
	WorkGroup string `json:",omitempty"`
	Ip        []string
}

type Plugins struct {
	PortOpened bool
	Web        Web       `json:",omitempty"`
	Cracked    []Account `json:",omitempty"`
	Mysql      []Account `json:",omitempty"`
	NetBios    NetBios   `json:",omitempty"`

	RateLimiter    *rate.Limiter
	NoTrust        bool `json:",omitempty"`
	Worker         int  `json:",omitempty"`
	TargetIp       string
	TargetPort     string
	TargetProtocol string
	TimeOut        int `json:"-"`
	DescCallback   func(string)
	RateWait       func(*rate.Limiter)
	highLight      bool
	locker         sync.RWMutex
}

const (
	SSHSrv = iota
	MysqlSrv
	WEBSrv  //放到最后
	PluginNR
)

const (
	OKNone = iota
	OKWait
	OKTimeOut
	OKDone
	OKNext
	OKStop
)
