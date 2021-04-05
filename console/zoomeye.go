package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/zsdevX/DarkEye/common"
	"github.com/zsdevX/DarkEye/zoomeye"
	"strconv"
	"strings"
)

type zoomEyeRuntime struct {
	api    string
	search string
	page   string

	flagSet *flag.FlagSet
}

var (
	zoomEyeRuntimeRuntimeOptions = &zoomEyeRuntime{
		flagSet: flag.NewFlagSet("zoomEye", flag.ExitOnError),
	}
)

func zoomEyeInitRunTime() {
	zoomEyeRuntimeRuntimeOptions.flagSet.StringVar(&zoomEyeRuntimeRuntimeOptions.api,
		"api", "you-key", "API-KEY")
	zoomEyeRuntimeRuntimeOptions.flagSet.StringVar(&zoomEyeRuntimeRuntimeOptions.search,
		"search", "ip:8.8.8.8", "https://www.zoomeye.org/")
	zoomEyeRuntimeRuntimeOptions.flagSet.StringVar(&zoomEyeRuntimeRuntimeOptions.page,
		"page", "5", "返回查询页面数量")
}

func (a *zoomEyeRuntime) compileArgs(cmd []string) error {
	ret := make([]string, 0)
	search := []string{"-search"}
	s := ""
	for _, c := range cmd {
		switch c {
		case "-api":
			fallthrough
		case "-page":
			ret = append(ret, strings.SplitN(c, " ", 2)...)
		default:
			s += strings.ReplaceAll(c, " ", "") + " "
		}
	}
	if s == "" {
		return fmt.Errorf("搜索参数为空")
	}
	search = append(search, s)
	ret = append(ret, search...)
	if err := a.flagSet.Parse(ret); err != nil {
		return err
	}
	a.flagSet.Parsed()
	return nil
}

func (a *zoomEyeRuntime) usage() {
	fmt.Println(fmt.Sprintf("Usage of %s:", zoomEye))
	fmt.Println("Options:")
	a.flagSet.VisitAll(func(f *flag.Flag) {
		var opt = "  -" + f.Name
		fmt.Println(opt)
		fmt.Println(fmt.Sprintf("		%v (default '%v')", f.Usage, f.DefValue))
	})
}

func (a *zoomEyeRuntime) start(ctx context.Context) {
	z := zoomeye.New()
	z.Query = a.search
	z.ApiKey = a.api
	z.Pages, _ = strconv.Atoi(a.page)
	z.ErrChannel = make(chan string, 10)
	go func() {
		for {
			select {
			case <-ctx.Done():
				common.StopIt(&z.Stop)
				return
			case m, ok := <-z.ErrChannel:
				if !ok {
					return
				}
				fmt.Println(m)
			default:
			}
			msg := <-z.ErrChannel
			fmt.Println(msg)
		}
	}()
	matches := z.Run()
	for _, m := range matches {
		e := &analysisEntity{
			Ip:              m.Ip,
			Port:            strconv.Itoa(m.Port),
			Country:         m.Country,
			Service:         m.Service,
			Url:             m.Url,
			Title:           m.Title,
			WebServer:       m.App,
			WebResponseCode: int32(m.HttpCode),
			Hostname:        m.Hostname,
			Os:              m.Os,
			Device:          m.Device,
			Banner:          m.Banner,
			Version:         m.Version,
			ExtraInfo:       m.ExtraInfo,
			RDns:            m.RDns,
		}
		analysisRuntimeOptions.createOrUpdate(e)
	}

}
