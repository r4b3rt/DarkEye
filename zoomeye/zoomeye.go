package zoomeye

import (
	"fmt"
	"github.com/zsdevX/DarkEye/common"
	"strconv"
	"strings"
	"time"
)

//Run add comment
func (z *ZoomEye) Run() []Match {
	ret := make([]Match, 0)
	i := 1
	for i <= z.Pages {
		matches := z.run(i)
		for _, match := range matches.Array() {
			m := Match{}
			m.Ip = match.Get("ip").String()
			m.Country = match.Get("geoinfo.country.names.en").String()

			m.Port = int(match.Get("portinfo.port").Int())
			m.Os = match.Get("portinfo.os").String()
			m.Hostname = match.Get("portinfo.hostname").String()
			m.Service = match.Get("portinfo.service").String()
			m.Banner = match.Get("portinfo.banner").String()[:128]
			m.Title = match.Get("portinfo.title").String()
			m.Version = match.Get("portinfo.version").String()
			m.Device = match.Get("portinfo.device").String()
			m.ExtraInfo = match.Get("portinfo.extrainfo").String()
			m.RDns = match.Get("portinfo.rdns").String()
			m.App = match.Get("portinfo.app").String()

			if m.Service == "http" || m.Service == "https" {
				m.Url = m.Service + "://" + m.Ip + ":" + strconv.Itoa(m.Port)
				if m.Banner != "" {
					x := strings.Split(m.Banner, " ")
					if len(x) == 3 {
						m.HttpCode, _ = strconv.Atoi(x[2])
					}
				}
			}
			ret = append(ret, m)
		}
		z.ErrChannel <- common.LogBuild("zoomEye",
			fmt.Sprintf("%s:获取第%d页信息", z.Query, i), common.INFO)
		time.Sleep(time.Second * 3)
	}
	return ret
}
