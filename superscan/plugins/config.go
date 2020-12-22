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
	HostName string   `json:",omitempty"`
	UserName string   `json:",omitempty"`
	Ip       []string `json:",omitempty"`
	Shares   []string `json:",omitempty"`
}

type Poc struct {
	Desc string `json:",omitempty"`
}

type Account struct {
	Username string `json:",omitempty"`
	Password string `json:",omitempty"`
	Web      `json:",omitempty"`
	NetBios  `json:",omitempty"`
	Poc      `json:",omitempty"`
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

	NoTrust        bool `json:",omitempty"`
	Worker         int  `json:",omitempty"`
	TargetIp       string
	TargetPort     string
	TargetProtocol string
	TimeOut        int `json:"-"`
	DescCallback   func(string)
	highLight      bool
	sync.RWMutex //protect 'Cracked'
}

type Config struct {
	//自定义字典
	UserList []string
	PassList []string
	//Poc反弹验证地址
	ReverseUrl      string
	ReverseCheckUrl string
	//发包速度限制
	Pps      *rate.Limiter
	RateWait func(*rate.Limiter)
}

//check list
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

//pre-check list
const (
	NetBiosPre = iota
	Ms17010Pre
	SnmpPre
	PluginPreCheckNR
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
