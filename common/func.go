package common

import (
	"context"
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

//FromTo add comment
type FromTo struct {
	From int
	To   int
}

func IPValid(ip string) bool {
	re := regexp.MustCompile(`\d{1,3}.\d{1,3}.\d{1,3}.\d{1,3}-\d{1,3}.\d{1,3}.\d{1,3}.\d{1,3}|\d{1,3}.\d{1,3}.\d{1,3}.\d{1,3}`)
	x := re.FindAllString(ip, -1)
	if x == nil {
		return false
	}
	if x[0] != ip {
		return false
	}

	return true
}

//GetIPRange add comment
func GetIPRange(ip string) (base string, start int, end string, err error) {
	err = fmt.Errorf("IP格式错误(eg. 1.1.1.1-1.1.1.255)")
	if !IPValid(ip) {
		return
	}

	start = 0
	fromTo := strings.Split(ip, "-")
	base = fromTo[0]

	end = base
	if len(fromTo) == 2 {
		end = fromTo[1]
	}
	err = nil
	return
}

//GenIP add comment
func GenIP(ipSeg string, ip int) string {
	a := make([]byte, 4)
	x := net.ParseIP(ipSeg).To4()
	y := binary.BigEndian.Uint32(x) + uint32(ip)
	binary.BigEndian.PutUint32(a, y)
	return net.IPv4(a[0], a[1], a[2], a[3]).String()
}

//CompareIP add comment
func CompareIP(a, b string) int64 {
	x := binary.BigEndian.Uint32(net.ParseIP(a).To4())
	y := binary.BigEndian.Uint32(net.ParseIP(b).To4())
	return int64(x) - int64(y)
}

//GetPortRange add comment
func GetPortRange(portRange string) ([]FromTo, int) {
	res := make([]FromTo, 0)
	tot := 0
	ports := strings.Split(portRange, ",")
	for _, port := range ports {
		from := 0
		to := 0
		fromTo := strings.Split(port, "-")
		from, _ = strconv.Atoi(fromTo[0])
		to = from
		if len(fromTo) == 2 {
			to, _ = strconv.Atoi(fromTo[1])
		}
		a := FromTo{
			From: from,
			To:   to,
		}
		res = append(res, a)
		tot += 1 + to - from
	}
	return res, tot
}

//IsAlive add comment
func IsAlive(parent context.Context, ip, port string, millTimeOut int) int {
	conn, err := DialCtx(parent, "tcp", net.JoinHostPort(ip, port), time.Duration(millTimeOut)*time.Millisecond)
	if err != nil {
		if eo, ok := err.(net.Error); ok {
			if eo.Timeout() {
				return TimeOut
			}
		}
		return Die
	}

	_ = conn.Close()
	return Alive
}

func DialCtx(parent context.Context, protocol, addr string, timeOut time.Duration) (net.Conn, error) {
	d := net.Dialer{Timeout: timeOut}
	ctx, _ := context.WithTimeout(parent, timeOut)
	return d.DialContext(ctx, protocol, addr)

}

func ParseFileOrVariable(name string) []string {
	if name != "" {
		if _, e := os.Stat(name); e == nil {
			return GenDicFromFile(name)
		} else {
			return strings.Split(name, ",")
		}
	} else {
		return nil
	}
}
