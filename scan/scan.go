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

const (
	Discovery int = iota
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

type IdListType map[string]int

var (
	IdList = IdListType{
		"tcp":       Discovery,
		"ping":      Discovery,
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

func (id IdListType) Id(name string) int {
	i, ok := id[name]
	if ok {
		return i
	}
	return Unknown
}

func (id IdListType) Name(sid int) string {
	for k, v := range id {
		if sid == v {
			return k
		}
	}
	return "unknown"
}

//New says @timeout:millisecond
func New(id, timeout int, args ...interface{}) (Scan, error) {
	u, p := genUsePass(IdList.Name(id))
	switch id {
	case Discovery:
		return NewDiscovery(timeout, args)
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
