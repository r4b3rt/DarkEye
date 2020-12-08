package common

import (
	"bufio"
	"context"
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

type FromTo struct {
	From int
	To   int
}

func GetIPRange(ip string) (base string, start, end int, err error) {
	fromTo := strings.Split(ip, "-")
	ipStart := fromTo[0]
	err = fmt.Errorf(LogBuild("common.func", "IP格式错误(eg. 1.1.1.1-3)", FAULT))
	if strings.Contains(ip, "/") {
		return
	}
	start = 0
	tIp := strings.Split(ipStart, ".")
	end = start
	if len(fromTo) == 2 {
		end, _ = strconv.Atoi(fromTo[1])
	}
	for _, v := range tIp {
		base += fmt.Sprintf("%s.", v)
	}
	base = strings.TrimRight(base, ".")
	err = nil
	return
}

func GenIP(ipSeg string, ip int) string {
	a := make([]byte, 4)
	binary.BigEndian.PutUint32(a, uint32(ip))
	b := strings.Split(ipSeg, ".")
	bl := len(b)
	for k := range a {
		if bl == 0 {
			break
		}
		v, _ := strconv.Atoi(b[k])
		a[k] += uint8(v)
		bl --
	}
	return net.IPv4(a[0], a[1], a[2], a[3]).String()
}

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

func ImportFiles(f, cnt string) (string, error) {
	file, err := os.Open(f)
	r := ""
	if err != nil {
		return r, err
	}
	defer file.Close()
	r = cnt
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		one := scanner.Text()
		if strings.HasPrefix(one, "#") {
			continue
		}
		one = strings.TrimSpace(one)
		one = strings.Trim(one, "\r\n")
		if one == "" {
			continue
		}
		if r == "" {
			r += one
		} else {
			r += "," + one
		}
	}
	return r, nil
}

func IsAlive(ip, port string, millTimeOut int) int {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Duration(millTimeOut)*time.Millisecond)
	defer cancel()
	//start := time.Now()
	d := &net.Dialer{}
	c, err := d.DialContext(ctx, "tcp4", ip+":"+port)
	//duration := time.Now().Sub(start)
	if err != nil {
		//fmt.Println("timeout duration", duration, err.Error(), millTimeOut)
		if eo, ok := err.(net.Error); ok {
			if eo.Timeout() {
				return TimeOut
			}
		}
		return Die
	}
	_ = c.Close()
	return Alive
}
