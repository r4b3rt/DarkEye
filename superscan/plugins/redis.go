package plugins

import (
	"github.com/go-redis/redis"
	"github.com/zsdevX/DarkEye/superscan/dic"
	"strings"
	"time"
)

var (
	redisUsername = make([]string, 0)
	redisPassword = make([]string, 0)
)

func init() {
	checkFuncs[RedisSrv] = redisCheck
	redisUsername = dic.DIC_USERNAME_REDIS
	redisPassword = dic.DIC_PASSWORD_REDIS
	supportPlugin["redis"] = "redis"
}

func redisCheck(plg *Plugins) {
	if !plg.NoTrust && plg.TargetPort != "6379" {
		return
	}
	crack("redis", plg, redisUsername, redisPassword, redisConn)
}

func redisConn(plg *Plugins, _ string, pass string) (ok int) {
	ok = OKNext
	client := redis.NewClient(&redis.Options{
		Addr:     plg.TargetIp + ":" + plg.TargetPort,
		Password: pass, DB: 0, DialTimeout: time.Millisecond * time.Duration(plg.TimeOut)})
	defer client.Close()
	pong, err := client.Ping(client.Context()).Result()
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
