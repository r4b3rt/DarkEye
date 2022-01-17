package scan

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"github.com/sirupsen/logrus"
	"net"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type mysqlConf struct {
	timeout  int
	username []string
	password []string
	logger   *logrus.Logger
}

func NewMysql(timeout int) (Scan, error) {
	s := &mysqlConf{
		timeout:  timeout,
	}

	return s, nil
}

func (s *mysqlConf) Start(parent context.Context, ip, port string) (interface{}, error) {
	addr := net.JoinHostPort(ip, port)
	return weakPass(parent, "mysql", addr, s.username, s.password, s.crack)
}

func (s *mysqlConf) Setup(args ...interface{}) {
	s.username = args[0].([]string)
	s.password = args[1].([]string)
	s.logger = args[2].(*logrus.Logger)
}

func (s *mysqlConf) crack(parent context.Context, addr, user, pass string) bool {
	source := fmt.Sprintf("%v:%v@tcp(%v)/%v?timeout=%dms&readTimeout=%dms",
		user, pass, addr, "mysql", s.timeout, s.timeout)

	db, err := sql.Open("mysql", source)
	if err != nil {
		s.logger.Debug("mysqlConf.open:", err.Error())
		return false
	}
	defer db.Close()
	db.SetConnMaxLifetime(time.Duration(s.timeout) * time.Millisecond)
	_, err = db.QueryContext(parent, `select 1;`)
	if err != nil {
		s.logger.Debug("mysqlConf.QueryContext", err.Error())
		return false
	}
	return true
}

func (s *mysqlConf) Identify(parent context.Context, ip, port string) bool {
	b, err := hello(parent, "tcp", net.JoinHostPort(ip, port), nil, s.timeout)
	if err != nil {
		return false
	}
	return bytes.Contains(b, []byte("mysql_native_password"))
}

func (s *mysqlConf) Attack(parent context.Context, ip, port string) error {
	return fmt.Errorf("not support")
}

func (s *mysqlConf) Output() interface{} {
	return nil
}
