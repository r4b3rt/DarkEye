package plugins

import (
	"context"
	"github.com/b1gcat/DarkEye/superscan/dic"
	"github.com/go-redis/redis"
	"strings"
	"time"
)

func redisCheck(s *Service) {
	s.crack()
}

func redisConn(parent context.Context, s *Service, _, pass string) (ok int) {
	ok = OKNext
	client := redis.NewClient(&redis.Options{
		Addr:        s.parent.TargetIp + ":" + s.parent.TargetPort,
		Password:    pass,
		DB:          0,
		DialTimeout: time.Millisecond * time.Duration(Config.TimeOut),
	})
	ctx, _ := context.WithCancel(parent)
	defer func() {
		if ok == OKDone || ok == OKNoAuth {
			if Config.Attack {
				redisAttack(ctx, s, client)
			}
		}
		_ = client.Close()
	}()

	pong, err := client.Ping(ctx).Result()
	if err == nil {
		if strings.Contains(pong, "PONG") {
			ok = OKDone
		} else {
			ok = OKStop
		}
	} else {
		if strings.Contains(err.Error(), " without any password configured") {
			ok = OKNoAuth
			return
		}
		//账号密码错误
		if strings.Contains(err.Error(), "username-password") ||
			strings.Contains(err.Error(), "OAUTH Authentication required") {
			return
		}
		if strings.Contains(err.Error(), "connection reset by peer") {
			//连接限制
			ok = OKWait
			return
		}
		if strings.Contains(err.Error(), "i/o timeout") {
			ok = OKTimeOut
			return
		}
		if strings.Contains(err.Error(), " protected mode") {
			//配置限制访问
			ok = OKForbidden
			return
		}
		ok = OKStop
	}
	return
}

func init() {
	services["redis"] = Service{
		name:    "redis",
		port:    "6379",
		user:    dic.DIC_USERNAME_REDIS,
		pass:    dic.DIC_PASSWORD_REDIS,
		check:   redisCheck,
		connect: redisConn,
		thread:  1,
	}
}
