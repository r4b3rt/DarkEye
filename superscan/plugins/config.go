package plugins

import (
	"context"
	"github.com/elastic/beats/libbeat/common/atomic"
	"github.com/zsdevX/DarkEye/superscan/dic"
	"golang.org/x/time/rate"
	"sync"
	"time"
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
	//HostName string   `json:",omitempty"`
	//UserName string   `json:",omitempty"`
	Ip     string `json:",omitempty"`
	Os     string
	Shares string `json:",omitempty"`
}

//Account add comment
type Account struct {
	Username string `json:",omitempty"`
	Password string `json:",omitempty"`
}

type tmpCache struct {
	urlPath string
	tls     bool
	cookie  string
}

//Plugins add comment
type Plugins struct {
	//Result
	PortOpened bool
	Cracked    []Account `json:",omitempty"`
	Web        Web       `json:",omitempty"`
	NetBios    NetBios   `json:",omitempty"`

	//Request
	TargetIp       string
	TargetPort     string
	TargetProtocol string

	//tmpValues
	TimeOut   int `json:"-"`
	highLight bool
	sync.RWMutex //protect 'Cracked'
	tmp tmpCache
}

//Config add comment
type Config struct {
	//自定义字典
	UserList []string
	PassList []string
	//发包速度限制
	Pps         *rate.Limiter
	RateWait    func(*rate.Limiter)
	UsingPlugin string
	Stop        atomic.Bool
	Ctx         context.Context
	Thread      int
}

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

type funcDesc struct {
	port string
	name string
	user []string
	pass []string
	doit func(*Plugins, *funcDesc)
}

var (
	checkFuncs    = make(map[string]funcDesc, 0)
	preCheckFuncs = make(map[string]funcDesc, 0)
	//GlobalConfig add comment
	GlobalConfig = Config{
		UsingPlugin: "",
		Thread:      2,
		RateWait: func(r *rate.Limiter) {
			if r == nil {
				return
			}
			for {
				if r.Allow() {
					break
				} else {
					time.Sleep(time.Millisecond * 10)
				}
			}
		},
	}
)

func init() {
	checkFuncs["ftp"] = funcDesc{
		name: "ftp",
		port: "21",
		doit: ftpCheck,
		user: dic.DIC_USERNAME_FTP,
		pass: dic.DIC_PASSWORD_FTP,
	}

	checkFuncs["memcached"] = funcDesc{
		name: "memcached",
		port: "11211",
		doit: memcachedCheck,
	}

	checkFuncs["mongodb"] = funcDesc{
		name: "mongodb",
		port: "27017",
		doit: mongoCheck,
		user: dic.DIC_USERNAME_MONGODB,
		pass: dic.DIC_PASSWORD_MONGODB,
	}

	checkFuncs["mysql"] = funcDesc{
		name: "mysql",
		port: "3306",
		doit: mysqlCheck,
		user: dic.DIC_USERNAME_MYSQL,
		pass: dic.DIC_PASSWORD_MYSQL,
	}

	checkFuncs["mssql"] = funcDesc{
		name: "mssql",
		port: "1433",
		doit: mssqlCheck,
		user: dic.DIC_USERNAME_SQLSERVER,
		pass: dic.DIC_PASSWORD_SQLSERVER,
	}

	checkFuncs["postgres"] = funcDesc{
		name: "postgres",
		port: "5432",
		doit: postgresCheck,
		user: dic.DIC_USERNAME_POSTGRESQL,
		pass: dic.DIC_PASSWORD_POSTGRESQL,
	}

	checkFuncs["redis"] = funcDesc{
		name: "redis",
		port: "6379",
		doit: redisCheck,
		user: dic.DIC_USERNAME_REDIS,
		pass: dic.DIC_PASSWORD_REDIS,
	}

	checkFuncs["smb"] = funcDesc{
		name: "smb",
		port: "445",
		doit: msbCheck,
		user: dic.DIC_USERNAME_SMB,
		pass: dic.DIC_PASSWORD_SMB,
	}

	checkFuncs["ssh"] = funcDesc{
		name: "ssh",
		port: "22",
		doit: sshCheck,
		user: dic.DIC_USERNAME_SSH,
		pass: dic.DIC_PASSWORD_SSH,
	}

	checkFuncs["web"] = funcDesc{
		name: "web",
		doit: webCheck,
	}
	///////// UDP check
	preCheckFuncs["ms17010"] = funcDesc{
		name: "ms17010",
		port: "445",
		doit: ms17010Check,
	}

	preCheckFuncs["netbios"] = funcDesc{
		name: "netbios",
		port: "135",
		doit: nbCheck,
	}

	preCheckFuncs["snmp"] = funcDesc{
		name: "snmp",
		port: "161",
		doit: snmpCheck,
	}
}
