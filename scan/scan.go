package scan

import (
	"context"
	"fmt"
	"github.com/b1gcat/DarkEye/dict"
	"strings"
)

type Scan interface {
	Identify(parent context.Context, ip, port string) bool
	Start(parent context.Context, ip, port string) (interface{}, error)
	Attack(parent context.Context, ip, port string) error
	Output() interface{}
}

type IdType int

const (
	Discovery IdType = iota
	DiscoHttp
	DiscoTcp
	DiscoPing
	DiscoNb
	DiscoEnd
	RiskStart
	Ssh
	Redis
	Mssql
	Ftp
	Memcached
	Mongodb
	Mysql
	Postgres
	RiskEnd
	Unknown
)

type IdListType map[string]IdType

var (
	IdList = IdListType{
		"tcp":       DiscoTcp,
		"ping":      DiscoPing,
		"http":      DiscoHttp,
		"nb":        DiscoNb,
		"ssh":       Ssh,
		"redis":     Redis,
		"mssql":     Mssql,
		"ftp":       Ftp,
		"memcached": Memcached,
		"mongodb":   Mongodb,
		"mysql":     Mysql,
		"postgres":  Postgres,
	}
)

func (id IdListType) String() string {
	r := make([]string, 0)
	for k := range IdList {
		r = append(r, k)
	}
	return strings.Join(r, ",")
}

func (id IdListType) Id(name string) IdType {
	i, ok := id[name]
	if ok {
		return i
	}
	return Unknown
}

func (id IdType) String() string {
	for k, v := range IdList {
		if id == v {
			return k
		}
	}
	return "unknown"
}

//New says @timeout:millisecond
func New(id IdType, timeout int, args ...interface{}) (Scan, error) {
	u, p := genUsePass(id.String())
	switch id {
	case DiscoPing:
		return NewDiscovery(timeout, DiscoPing)
	case DiscoTcp:
		return NewDiscovery(timeout, DiscoTcp)
	case DiscoHttp:
		return NewDiscovery(timeout, DiscoHttp)
	case DiscoNb:
		return NewDiscovery(timeout, DiscoHttp)
	case Ssh:
		return NewSSh(timeout, append(args, u, p))
	case Redis:
		return NewRedis(timeout, append(args, u, p))
	case Mssql:
		return NewMssql(timeout, append(args, u, p))
	case Ftp:
		return NewFtp(timeout, append(args, u, p))
	case Memcached:
		return NewMemCache(timeout, append(args, u, p))
	case Mongodb:
		return NewMongodb(timeout, append(args, u, p))
	case Mysql:
		return NewMysql(timeout, append(args, u, p))
	case Postgres:
		return NewPostgres(timeout, append(args, u, p))
	default:
		return nil, fmt.Errorf("不支持的扫描类型 %v", id)
	}
}

func genUsePass(name string) ([]string, []string) {
	user := fmt.Sprintf("dict_%s_username", name)
	pass := fmt.Sprintf("dict_%s_password", name)

	u, err := dict.Asset(user)
	if err != nil {
		return nil, nil
	}
	p, err := dict.Asset(pass)
	if err != nil {
		return nil, nil
	}

	ul := strings.Split(string(u), "\n")
	for k, v := range ul {
		if v == "空" {
			ul[k] = ""
		}
	}
	pl := strings.Split(string(p), "\n")
	for k, v := range pl {
		if v == "空" {
			pl[k] = ""
		}
	}

	return ul, pl
}
