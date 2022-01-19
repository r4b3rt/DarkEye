package scan

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"strings"
)

func weakPass(parent context.Context, s, addr string, ul, pl []string,
	cb func(context.Context, string, string, string) bool) (interface{}, error) {
	ctx, _ := context.WithCancel(parent)
	for _, u := range ul {
		for _, p := range pl {
			p = strings.ReplaceAll(p, `%user%`, u)
			select {
			case <-ctx.Done():
				return nil, nil
			default:
			}
			if cb(ctx, addr, u, p) {
				return fmt.Sprintf("%s %s - %s/%s", s, addr, u, p), nil
			}
		}
	}
	return nil, nil
}

func preNUm(data byte) int {
	var mask byte = 0x80
	var num int = 0
	//8bit中首个0bit前有多少个1bits
	for i := 0; i < 8; i++ {
		if (data & mask) == mask {
			num++
			mask = mask >> 1
		} else {
			break
		}
	}
	return num
}

func isUtf8(data []byte) bool {
	i := 0
	for i < len(data) {
		if (data[i] & 0x80) == 0x00 {
			// 0XXX_XXXX
			i++
			continue
		} else if num := preNUm(data[i]); num > 2 {
			// 110X_XXXX 10XX_XXXX
			// 1110_XXXX 10XX_XXXX 10XX_XXXX
			// 1111_0XXX 10XX_XXXX 10XX_XXXX 10XX_XXXX
			// 1111_10XX 10XX_XXXX 10XX_XXXX 10XX_XXXX 10XX_XXXX
			// 1111_110X 10XX_XXXX 10XX_XXXX 10XX_XXXX 10XX_XXXX 10XX_XXXX
			// preNUm() 返回首个字节的8个bits中首个0bit前面1bit的个数，该数量也是该字符所使用的字节数
			i++
			for j := 0; j < num-1; j++ {
				//判断后面的 num - 1 个字节是不是都是10开头
				if (data[i] & 0xc0) != 0x80 {
					return false
				}
				i++
			}
		} else {
			//其他情况说明不是utf-8
			return false
		}
	}
	return true
}

func setupRisk(r *risk, args []interface{}) {
	for k, v := range args {
		switch v.(type) {
		case *logrus.Logger:
			r.logger = v.(*logrus.Logger)
		case []string:
			switch k {
			case 1:
				x, ok := v.([]string)
				if !ok {
					break
				}
				if len(x) != 0 {
					r.logger.Debug("user dict change to:", x)
					r.username = x
				}
			case 2:
				x, ok := v.([]string)
				if !ok {
					break
				}
				if len(x) != 0 {
					r.logger.Debug("pass dict change to:", x)
					r.password = x
				}
			}
		default:
		}
	}
}

type risk struct {
	username []string
	password []string
	logger   *logrus.Logger
}
