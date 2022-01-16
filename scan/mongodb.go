package scan

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2"
	"net"
	"strings"
	"time"
)

type mongodbConf struct {
	timeout  int
	username []string
	password []string
	logger   *logrus.Logger
}

func NewMongodb(timeout int, args []interface{}) (Scan, error) {
	s := &mongodbConf{
		timeout:  timeout,
		username: args[0].([]string),
		password: args[1].([]string),
		logger:   args[2].(*logrus.Logger),
	}

	return s, nil
}

func (s *mongodbConf) Start(parent context.Context, ip, port string) (interface{}, error) {
	addr := net.JoinHostPort(ip, port)
	if unAuth, err := mgo.DialWithTimeout(
		"mongodb://"+addr+"/"+"admin",
		time.Duration(s.timeout)*time.Millisecond); err == nil {
		if _, err = unAuth.DatabaseNames(); err != nil {
			s.logger.Debug("un-auth.Ping:", err.Error())
		} else {
			return fmt.Sprintf("mongodb %s unauth", net.JoinHostPort(ip, port)), nil
		}
	}
	return weakPass(parent, "mongodb", addr, s.username, s.password, s.crack)
}

func (s *mongodbConf) crack(parent context.Context, addr, user, pass string) bool {
	_, err := mgo.DialWithTimeout(
		"mongodb://"+user+":"+pass+"@"+addr+"/"+"admin",
		time.Duration(s.timeout)*time.Millisecond)
	if err != nil {
		//mog.v2 not support mongodb 5.x
		if strings.Contains(err.Error(), "is an unknown field") {
			return true
		}
		s.logger.Debug("mongodbConf:", err.Error())
		return false
	}
	return true
}

func (s *mongodbConf) Identify(_ context.Context, ip, port string) bool {
	_, err := mgo.DialWithTimeout(
		"mongodb://"+net.JoinHostPort(ip, port)+"/"+"admin",
		time.Duration(s.timeout)*time.Millisecond)
	if err == nil ||
		strings.Contains(err.Error(), "Authentication failed") {
		return true
	}
	return false
}

func (s *mongodbConf) Attack(parent context.Context, ip, port string) error {
	return fmt.Errorf("not support")
}

func (s *mongodbConf) Output() interface{} {
	return nil
}
