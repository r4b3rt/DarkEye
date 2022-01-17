package scan

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/sirupsen/logrus"
	"net"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

type postgresConf struct {
	timeout  int
	username []string
	password []string
	logger   *logrus.Logger
}

func NewPostgres(timeout int) (Scan, error) {
	s := &postgresConf{
		timeout:  timeout,
	}

	return s, nil
}

func (s *postgresConf) Start(parent context.Context, ip, port string) (interface{}, error) {
	addr := net.JoinHostPort(ip, port)
	return weakPass(parent, "postgres", addr, s.username, s.password, s.crack)
}

func (s *postgresConf) Setup(args ...interface{}) {
	s.username = args[0].([]string)
	s.password = args[1].([]string)
	s.logger = args[2].(*logrus.Logger)
}

func (s *postgresConf) crack(parent context.Context, addr, user, pass string) bool {
	timeOut := time.Millisecond * time.Duration(s.timeout)
	db, err := s.auth(parent, addr, user, pass)
	if err != nil {
		s.logger.Debug("postgresConf.open:", err.Error())
		return false
	}
	defer db.Close()
	db.SetConnMaxLifetime(timeOut)
	err = db.Ping()
	if err != nil {
		s.logger.Debug("postgresConf.ping:", err.Error())
		return false
	}
	return true
}

func (s *postgresConf) auth(_ context.Context, addr, user, pass string) (*sql.DB, error) {
	host, port, _ := net.SplitHostPort(addr)
	source := fmt.Sprintf("host=%s	port=%s user=%s password=%s dbname=postgres sslmode=disable connect_timeout=%d",
		host, port, user, pass, s.timeout/1000+1)
	return sql.Open("postgres", source)
}

func (s *postgresConf) Identify(parent context.Context, ip, port string) bool {
	d, err := s.auth(parent, net.JoinHostPort(ip, port), "fuck", "fuck")
	if err != nil {
		s.logger.Debug("postgresConf.Identify:", err.Error())
		return false
	}
	defer d.Close()
	d.SetConnMaxLifetime(time.Millisecond * time.Duration(s.timeout))
	err = d.Ping()
	if err != nil {
		s.logger.Debug("postgresConf.Ping:", err.Error())
		if !strings.Contains(err.Error(), "password authentication") {
			return false
		}
	}
	return true
}

func (s *postgresConf) Attack(parent context.Context, ip, port string) error {
	return fmt.Errorf("not support")
}

func (s *postgresConf) Output() interface{} {
	return nil
}
