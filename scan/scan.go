package scan

import (
	"context"
	"fmt"
	"github.com/b1gcat/DarkEye/dict"
	"github.com/sirupsen/logrus"
	"strings"
)

type Scan interface {
	Identify(parent context.Context, ip, port string) bool
	Start(parent context.Context, ip, port string) (interface{}, error)
	Attack(parent context.Context, ip, port string) error
	Setup(args ...interface{})
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
		"tcp":        DiscoTcp,
		"ping":       DiscoPing,
		"http":       DiscoHttp,
		"nb":         DiscoNb,
		"ssh":        Ssh,
		"redis":      Redis,
		"mssql":      Mssql,
		"ftp":        Ftp,
		"memcached":  Memcached,
		"mongodb":    Mongodb,
		"mysql":      Mysql,
		"postgresql": Postgres,
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

/*New says
@timeout:millisecond
*/
func New(id IdType, timeout int) (Scan, error) {
	var s Scan
	var err error

	switch id {
	case DiscoPing:
		fallthrough
	case DiscoTcp:
		fallthrough
	case DiscoHttp:
		fallthrough
	case DiscoNb:
		s, err = NewDiscovery(timeout, id)
	case Ssh:
		s, err = NewSSh(timeout)
	case Redis:
		s, err = NewRedis(timeout)
	case Mssql:
		s, err = NewMssql(timeout)
	case Ftp:
		s, err = NewFtp(timeout)
	case Memcached:
		s, err = NewMemCache(timeout)
	case Mongodb:
		s, err = NewMongodb(timeout)
	case Mysql:
		s, err = NewMysql(timeout)
	case Postgres:
		s, err = NewPostgres(timeout)
	default:
		return nil, fmt.Errorf("不支持的扫描类型 %v", id)
	}
	if err != nil {
		return nil, err
	}
	switch {
	case id >= RiskStart && id <= RiskEnd:
		u, p := genUsePass(id.String())
		s.Setup(logrus.New(), u, p)
	case id >= Discovery && id <= DiscoEnd:
		s.Setup(logrus.New())
	}
	return s, err
}

func genUsePass(name string) ([]string, []string) {
	user := fmt.Sprintf("dic_username_%s.txt", name)
	pass := fmt.Sprintf("dic_password_%s.txt", name)

	u, err := dict.Asset(user)
	if err != nil {
		return nil, nil
	}
	p, err := dict.Asset(pass)
	if err != nil {
		return nil, nil
	}

	us := strings.ReplaceAll(string(u), "\r", "")
	ul := strings.Split(us, "\n")
	for k, v := range ul {
		if v == "空" {
			ul[k] = ""
		}
	}
	ps := strings.ReplaceAll(string(p), "\r", "")
	pl := strings.Split(ps, "\n")
	for k, v := range pl {
		if v == "空" {
			pl[k] = ""
		}
	}

	return ul, pl
}
