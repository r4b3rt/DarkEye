package common

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"time"
)

var (
	BaseDir        = "."
	UserAgents     = []string{"Mozilla/5.0 (Windows NT 6.1; WOW64; rv:34.0) Gecko/20100101 Firefox/34.0", "Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1; en) Opera 9.50", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/534.57.2 (KHTML, like Gecko) Version/5.1.7 Safari/534.57.2", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.71 Safari/537.36", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.11 (KHTML, like Gecko) Chrome/23.0.1271.64 Safari/537.11", "Mozilla/5.0 (Windows; U; Windows NT 6.1; en-US) AppleWebKit/534.16 (KHTML, like Gecko) Chrome/10.0.648.133 Safari/534.16", "Mozilla/5.0 (Windows NT 6.1; WOW64; Trident/7.0; rv:11.0) like Gecko", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/536.11 (KHTML, like Gecko) Chrome/20.0.1132.11 TaoBrowser/2.0 Safari/536.11", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Maxthon/4.4.3.4000 Chrome/30.0.1599.101 Safari/537.36", "Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 5.1; Trident/4.0; SV1; QQDownload 732; .NET4.0C; .NET4.0E; SE 2.X MetaSr 1.0)", "Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1; SV1; QQDownload 732; .NET4.0C; .NET4.0E; LBBROWSER)", "Mozilla/5.0 (Windows; U; Windows NT 6.1; en-us) AppleWebKit/534.50 (KHTML, like Gecko) Version/5.1 Safari/534.50", "Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.1; Trident/5.0", "Opera/9.80 (Windows NT 6.1; U; en) Presto/2.8.131 Version/11.11", "Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 5.1; TencentTraveler 4.0)"}
	ProgramVersion = "1.0." + fmt.Sprintf("%d%d%d%d%d\nhttps://github.com/zsdevX/DarkEye\n大橘Oo0\n84500316@qq.com",
		time.Now().Year(), time.Now().Month(), time.Now().Day(), time.Now().Hour(), time.Now().Minute())
	Banner = `___.   .__       _________         __   
\_ |__ |__| ____ \_   ___ \_____ _/  |_ 
 | __ \|  |/ ___\/    \  \/\__  \\   __\
 | \_\ \  / /_/  >     \____/ __ \|  |  
 |___  /__\___  / \______  (____  /__|  
     \/  /_____/         \/     \/`
)

var (
	LowCaseAlpha = "abcedfghijklmnopqrstuvwxyz"
)

func init() {
	BaseDir, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	logFile = filepath.Join(BaseDir, logFile)
}

func GenHumanSecond(base int) int {
	return base + (rand.Int() % 3)
}

func TrimUseless(a string) string {
	//a = strings.Replace(a, " ", " ", -1)
	a = strings.Replace(a, "\n", "", -1)
	a = strings.Replace(a, "\r", "", -1)
	a = strings.Replace(a, ",", " ", -1)
	return a
}

func StopIt(stop *int32) {
	atomic.StoreInt32(stop, 1)
}

func StartIt(stop *int32) {
	atomic.StoreInt32(stop, 0)
}

func ShouldStop(stop *int32) bool {
	return atomic.LoadInt32(stop) == 1
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

func ISUtf8(data []byte) bool {
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
