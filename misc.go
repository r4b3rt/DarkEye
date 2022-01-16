package main

import (
	"encoding/binary"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net"
	"os"
	"strconv"
	"strings"
)

//readList @l：文件或string列表(a,b,c)
func readList(l string) []string {
	if _, err := os.Stat(l); err != nil {
		return strings.Split(l, ",")
	}
	lb, err := ioutil.ReadFile(l)
	if err != nil {
		return []string{}
	}
	return strings.Split(string(lb), "\n")
}

func portSplit(p string) []string {
	var list []string
	if _, err := os.Stat(p); err != nil {
		list = strings.Split(p, ",")
	} else {
		lb, err := ioutil.ReadFile(p)
		if err != nil {
			return []string{}
		}
		list = strings.Split(string(lb), "\n")
	}

	e := make([]string, 0)
	for _, l := range list {
		if !strings.Contains(l, "-") {
			e = append(e, l)
			continue
		}
		//80-81
		ls := strings.Split(l, "-")
		if len(ls) != 2 {
			logrus.Error("not support port format:", l)
			return []string{}
		}
		f, _ := strconv.Atoi(ls[0])
		t, _ := strconv.Atoi(ls[1])
		if f < 0 || f > 65535 || t <= 0 || t > 65535 || f >= t {
			logrus.Error("not support port format:", l)
			return []string{}
		}
		for f <= t {
			e = append(e, strconv.Itoa(f))
			f++
		}
	}
	return e
}

// splitIp2C says ...
// support 1.1.1.1 or 1.1.1.1/24 or 1.1-2
func splitIp2C(ip string) ([]string, error) {
	minMaskSize := 16
	//1.1.1.1/24
	if strings.Contains(ip, "/") {
		_, ipn, err := net.ParseCIDR(ip)
		if err != nil {
			return nil, err
		}
		n, _ := ipn.Mask.Size()
		if n == 0 {
			return nil, fmt.Errorf("not support ip format %v", ip)
		}

		if n < minMaskSize {
			return nil, fmt.Errorf("too small mask %v < %v(min_mask_size)", n, minMaskSize)
		}

		x := 24 - n
		if x <= 0 {
			return []string{ip}, nil
		}
		r := make([]string, 0)
		next := ipn.IP
		loop := 0
		end := pow(2, x)*256 - 1 //how many ips
		for loop < end {
			r = append(r, next.To4().String())
			next = ipInc(next, 256)
			loop += 256
		}
		return r, nil
	} else if strings.Contains(ip, "-") { //1.1.1.1-234
		dot := strings.Split(ip, ".")
		switch len(dot) {
		case 0:
		case 1:
		case 2: //1.1-2
			ft := strings.Split(dot[1], "-")
			if len(ft) == 2 {
				f, _ := strconv.Atoi(ft[0])
				t, _ := strconv.Atoi(ft[1])
				if f < 0 || f > 255 || t <= 0 || t > 255 || f >= t {
					return nil, fmt.Errorf("not support ip format %v", ip)
				}
				r := make([]string, 0)
				for f <= t {
					ipc := fmt.Sprintf("%s.%d.0.0/16", dot[0], f)
					ips, err := splitIp2C(ipc)
					if err != nil {
						return nil, err
					}
					f++
					r = append(r, ips...)
				}
				return r, nil
			}
		case 3: //1.1.1-2
			ft := strings.Split(dot[2], "-")
			if len(ft) == 2 {
				f, _ := strconv.Atoi(ft[0])
				t, _ := strconv.Atoi(ft[1])
				if f < 0 || f > 255 || t <= 0 || t > 255 || f >= t {
					return nil, fmt.Errorf("not support ip format %v", ip)
				}
				r := make([]string, 0)
				for f <= t {
					ipc := fmt.Sprintf("%s.%s.%d.0", dot[0], dot[1], f)
					f++
					r = append(r, ipc)
				}
				return r, nil
			}
		case 4: //1.1.1.1-2
			return []string{ip}, nil
		}
	} else { //1.1.1.1
		if net.ParseIP(ip) != nil {
			return []string{ip}, nil
		}
	}

	return nil, fmt.Errorf("not support ip format %v", ip)
}

//splitIpC2Ip split ip C to ips
func splitIpC2Ip(ip string) ([]net.IP, error) {
	if strings.Contains(ip, "-") { //1.1.1.1-2
		return _splitIpC2Ip(ip, "-")
	} else if strings.Contains(ip, "/") {
		return _splitIpC2Ip(ip, "/")
	} else {
		if net.ParseIP(ip) != nil {
			return []net.IP{net.ParseIP(ip).To4()}, nil
		}
	}
	return nil, fmt.Errorf("格式错误:%v", ip)
}

func _splitIpC2Ip(ip, sep string) ([]net.IP, error) {
	ips := strings.Split(ip, sep)
	if len(ips) != 2 {
		return nil, fmt.Errorf("not support format:%v", ip)
	}
	ips2 := strings.Split(ips[0], ".")
	if len(ips2) != 4 {
		return nil, fmt.Errorf("not support format:%v", ip)
	}
	f, _ := strconv.Atoi(ips2[3])
	t, _ := strconv.Atoi(ips[1])

	if sep == "/" {
		_, ipn, err := net.ParseCIDR(ip)
		if err != nil {
			return nil, fmt.Errorf("not support format:%v", ip)
		}
		n, _ := ipn.Mask.Size()
		l := pow(2, 32-n) - 2
		f = int(ipn.IP.To4()[3]) + 1
		t = f + l - 1
	}

	if f < 0 || f > 255 || t <= 0 || t > 255 || f >= t {
		return nil, fmt.Errorf("not support format:%v", ip)
	}
	r := make([]net.IP, 0)
	for f <= t {
		ipc := net.ParseIP(fmt.Sprintf("%s.%s.%s.%d", ips2[0], ips2[1], ips2[2], f)).To4()
		r = append(r, ipc)
		f++
	}
	return r, nil
}

func ipInc(ip net.IP, num int) net.IP {
	var r = make(net.IP, len(ip))
	y := binary.BigEndian.Uint32(ip.To4()) + uint32(num)
	binary.BigEndian.PutUint32(r, y)
	return r
}

func pow(x, n int) int {
	i := 0
	y := 1
	for i < n {
		y *= x
		i++
	}
	return y
}
