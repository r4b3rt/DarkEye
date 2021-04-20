package plugins

import (
	"context"
	"github.com/zsdevX/DarkEye/superscan/dic"
	"golang.org/x/time/rate"
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
	Ip     string `json:",omitempty"`
	Os     string
	Shares string `json:",omitempty"`
}

//Account add comment
type Account struct {
	Username string `json:",omitempty"`
	Password string `json:",omitempty"`
}

type Result struct {
	PortOpened  bool
	Cracked     Account `json:",omitempty"`
	Web         Web     `json:",omitempty"`
	NetBios     NetBios `json:",omitempty"`
	ServiceName string  `json:",omitempty"`
	ExpHelp     string  `json:"，omitempty"`
}

//Plugins add comment
type Plugins struct {
	//Result
	Result Result
	Hit    bool

	//Request
	TargetIp   string
	TargetPort string
}

//Config add comment
type config struct {
	TimeOut int
	//自定义字典
	UserList []string
	PassList []string
	//发包速度限制
	PPS *rate.Limiter
	//选择插件
	SelectPlugin string
	//程序控制
	ParentCtx context.Context
	//扫服务线程
	ServiceThread int
	//Attack
	Attack    bool
	SshPubKey string
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
	//OKTerm add comment
	OKTerm
)

type Service struct {
	name    string
	port    string
	user    []string
	pass    []string
	thread  int
	parent  *Plugins
	check   func(*Service)
	connect func(context.Context, *Service, string, string) int
	vars    map[string]string
}

var (
	services    = make(map[string]Service, 0)
	preServices = make(map[string]Service, 0)
	//GlobalConfig add comment
	Config = config{
		ServiceThread: 4,
		ParentCtx:     context.Background(),
		TimeOut:       1000, //millions
		SshPubKey:     "\n\n\nssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQC4tvzpcxSRxx51aVBLeetsu6J/OsDJTyQGt5LcLtbQDHzctGLTVzaXQ+NXRPnGXmLzIZP8/dn7SeEKGhPJmruByUEmJkhBln/Flgp1CUDtX/RJ7q/YkFTHdHYyq1zVG75y2/VpMfEMwP87UD7teZjbSKKeuD1SFfrXbwIqZruiRuOHXSNilsm3wINj8ZwhnxRo7IFBXSwtGA4TqCno1ngaDTzwHT+PKLIGt2n/5V2S7R/+EYneBiLAhQJ0b9GmW35RRZGsoWYKGSmytmPjd81GpEojjynKu4jsB/6F+IU9aH45KYzOF44yOZOwodj7mVIHtdL7kTE5y2rzaNNZH32qw7wM35WaiLjvHsqt9GAcLs88OMy9PSFb/41IrQEDdldxjzKCfAOKku6X0s3V1MfZPSy+foIcEy1sgfFm52nWaogNuBim1sYkq9lipwN88NhrvJH43afYv8/qe3ik+rKumAh3OqgUv4jNFMjBjpqUp+XUyIFjBouIUy/ORIUXm5E= root@b17ed2775c27\n\n\n",
	}
)

func init() {
	services["ftp"] = Service{
		name:    "ftp",
		port:    "21",
		user:    dic.DIC_USERNAME_FTP,
		pass:    dic.DIC_PASSWORD_FTP,
		check:   ftpCheck,
		connect: ftpConn,
		thread:  1,
	}
	services["memcached"] = Service{
		name:    "memcached",
		port:    "11211",
		check:   memCachedCheck,
		connect: memCacheConn,
		thread:  1,
	}
	services["mongodb"] = Service{
		name:    "mongodb",
		port:    "27017",
		user:    dic.DIC_USERNAME_MONGODB,
		pass:    dic.DIC_PASSWORD_MONGODB,
		check:   mongoCheck,
		connect: mongodbConn,
		thread:  1,
	}
	services["mssql"] = Service{
		name:    "mssql",
		port:    "1433",
		user:    dic.DIC_USERNAME_SQLSERVER,
		pass:    dic.DIC_PASSWORD_SQLSERVER,
		check:   mssqlCheck,
		connect: mssqlConn,
		thread:  1,
	}
	services["mysql"] = Service{
		name:    "mysql",
		port:    "3306",
		user:    dic.DIC_USERNAME_MYSQL,
		pass:    dic.DIC_PASSWORD_MYSQL,
		check:   mysqlCheck,
		connect: mysqlConn,
		thread:  1,
	}
	services["postgres"] = Service{
		name:    "postgres",
		port:    "5432",
		user:    dic.DIC_USERNAME_POSTGRESQL,
		pass:    dic.DIC_PASSWORD_POSTGRESQL,
		check:   postgresCheck,
		connect: postgresConn,
		thread:  1,
	}
	services["redis"] = Service{
		name:    "redis",
		port:    "6379",
		user:    dic.DIC_USERNAME_REDIS,
		pass:    dic.DIC_PASSWORD_REDIS,
		check:   redisCheck,
		connect: redisConn,
		thread:  1,
	}
	services["smb"] = Service{
		name:    "smb",
		port:    "445",
		user:    dic.DIC_USERNAME_SMB,
		pass:    dic.DIC_PASSWORD_SMB,
		check:   smbCheck,
		connect: smbConn,
		thread:  1,
	}

	services["ssh"] = Service{
		name:    "ssh",
		port:    "22",
		user:    dic.DIC_USERNAME_SSH,
		pass:    dic.DIC_PASSWORD_SSH,
		check:   sshCheck,
		connect: sshConn,
		thread:  1,
	}
	services["web"] = Service{
		name:  "web",
		check: webCheck,
	}

	///////// pre check
	preServices["ms17010"] = Service{
		name:  "ms17010",
		port:  "445",
		check: ms17010Check,
	}
	preServices["netbios"] = Service{
		name:  "netbios",
		port:  "135",
		check: nbCheck,
	}

	preServices["snmp"] = Service{
		name:  "snmp",
		port:  "161",
		check: snmpCheck,
	}
}
