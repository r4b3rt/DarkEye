package common

import (
	"fmt"
	"strconv"
	"strings"
)

func ParseNmapIP(ip string) (base string, start, end int, err error) {
	fromTo := strings.Split(ip, "-")
	ipStart := fromTo[0]
	err = fmt.Errorf(LogBuild("common.func", "IP格式错误(eg. 1.1.1.1-3)", FAULT))

	tIp := strings.Split(ipStart, ".")
	if len(tIp) != 4 {
		return
	}
	start, _ = strconv.Atoi(tIp[3])
	end = start
	if len(fromTo) == 2 {
		end, _ = strconv.Atoi(fromTo[1])
	}
	if end == 0 {
		return
	}
	base = fmt.Sprintf("%s.%s.%s", tIp[0], tIp[1], tIp[2])
	err = nil
	return
}

