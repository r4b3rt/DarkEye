package plugins

import (
	"golang.org/x/time/rate"
	"sync"
)

//Web add comment
type Web struct {
	Server string `json:",omitempty"`
	Title  string `json:",omitempty"`
	Code   int32  `json:",omitempty"`
	Url    string `json:",omitempty"`
	Tls    bool   `json:",omitempty"`
}

//NetBios add comment
type NetBios struct {
	HostName string   `json:",omitempty"`
	UserName string   `json:",omitempty"`
	Ip       []string `json:",omitempty"`
	Shares   []string `json:",omitempty"`
}

//Poc add comment
type Poc struct {
	Desc string `json:",omitempty"`
}

//Account add comment
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

//Plugins add comment
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
	sync.RWMutex   //protect 'Cracked'
}

//Config add comment
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
	//SSHSrv add comment
	SSHSrv = iota
	//MysqlSrv add comment
	MysqlSrv
	//RedisSrv add comment
	RedisSrv
	//FtpSrv add comment
	FtpSrv
	//MongoSrv add comment
	MongoSrv
	//MemoryCacheSrv add comment
	MemoryCacheSrv
	//PostgresSrv add comment
	PostgresSrv
	//MSSQLSrv add comment
	MSSQLSrv
	//SmbSrv add comment
	SmbSrv
	//WEBSrv add comment
	WEBSrv //放到最后
	//PluginNr add comment
	PluginNR
)

//pre-check list
const (
	//NetBiosPre add comment
	NetBiosPre = iota
	//Ms17010Pre add comment
	Ms17010Pre
	//SnmpPre add comment
	SnmpPre
	//PluginPreCheckNR add comment
	PluginPreCheckNR
)

const (
	//OKNa add comment
	OKNa = iota
	//OKWait add comment
	OKWait
	//OKTimeOut add comment
	OKTimeOut
	//OKDone add comment
	OKDone
	//OKNext add comment
	OKNext
	//OKForbidden add comment
	OKForbidden
	//OKNoAuth add comment
	OKNoAuth
	//OKStop add comment
	OKStop
)
