package plugins

import (
	"golang.org/x/time/rate"
	"sync"
)

type Web struct {
	Server string `json:",omitempty"`
	Title  string `json:",omitempty"`
	Code   int32  `json:",omitempty"`
	Url    string `json:",omitempty"`
	Tls    bool   `json:",omitempty"`
}

type NetBios struct {
	HostName  string `json:",omitempty"`
	UserName  string `json:",omitempty"`
	WorkGroup string `json:",omitempty"`
	Ip        []string
}

type Poc struct {
	Desc string `json:",omitempty"`
}

type Account struct {
	Username  string `json:",omitempty"`
	Password  string `json:",omitempty"`
	Web       `json:",omitempty"`
	NetBios   `json:",omitempty"`
	Poc       `json:",omitempty"`
	PingAlive string
}

type tmpCache struct {
	urlPath string
	tls     bool
	cookie  string
}

type Plugins struct {
	PortOpened bool
	Cracked    []Account `json:",omitempty"`
	tmp        tmpCache

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
	sync.RWMutex //protect 'Cracked'
}

const (
	SSHSrv = iota
	MysqlSrv
	RedisSrv
	FtpSrv
	MongoSrv
	MemoryCacheSrv
	PostgresSrv
	MSSQLSrv
	SmbSrv
	WEBSrv  //放到最后
	PluginNR
)

const (
	OKNa = iota
	OKWait
	OKTimeOut
	OKDone
	OKNext
	OKForbidden
	OKNoauth
	OKStop
)
