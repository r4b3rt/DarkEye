package scan

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"strings"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
)

type mssqlConf struct {
	timeout  int
	risk
}

func NewMssql(timeout int) (Scan, error) {
	s := &mssqlConf{
		timeout:  timeout,
	}

	return s, nil
}

func (s *mssqlConf) Start(parent context.Context, ip, port string) (interface{}, error) {
	addr := net.JoinHostPort(ip, port)
	return weakPass(parent, "mssql", addr, s.username, s.password, s.crack)
}

func (s *mssqlConf) Setup(args ...interface{}) {
	setupRisk(&s.risk, args)
}

func (s *mssqlConf) crack(parent context.Context, addr, user, pass string) bool {
	if err := s.auth(parent, addr, user, pass); err != nil {
		s.logger.Debug(addr, ":", err.Error())
		return false
	} else {
		return true
	}
}

func (s *mssqlConf) auth(_ context.Context, addr, user, pass string) error {
	ip, port, _ := net.SplitHostPort(addr)
	source := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%s;encrypt=disable;timeout=%dms",
		ip, user, pass, port, time.Duration(s.timeout)*time.Millisecond)
	db, err := sql.Open("mssql", source)
	if err != nil {
		return fmt.Errorf("mssqlConf.Open:" + err.Error())
	}
	defer db.Close()
	db.SetConnMaxLifetime(time.Duration(s.timeout) * time.Millisecond)
	if err = db.Ping(); err != nil {
		return fmt.Errorf("mssqlConf.Ping:" + err.Error())
	}
	return nil
}

func (s *mssqlConf) Identify(parent context.Context, ip, port string) bool {
	err := s.auth(parent, net.JoinHostPort(ip, port), "fuck", "fuck")
	if err == nil || strings.Contains(err.Error(), "Login failed for user") {
		return true
	}
	return false
}

func (s *mssqlConf) Attack(parent context.Context, ip, port string) error {
	return fmt.Errorf("not support")
}

func (s *mssqlConf) Output() interface{} {
	return nil
}
