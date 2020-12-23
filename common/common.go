package common

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"time"
)

var (
	//BaseDir add comment
	BaseDir = "."
	//UserAgents add comment
	UserAgents = []string{"Mozilla/5.0 (Windows NT 6.1; WOW64; rv:34.0) Gecko/20100101 Firefox/34.0", "Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1; en) Opera 9.50", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/534.57.2 (KHTML, like Gecko) Version/5.1.7 Safari/534.57.2", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.71 Safari/537.36", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.11 (KHTML, like Gecko) Chrome/23.0.1271.64 Safari/537.11", "Mozilla/5.0 (Windows; U; Windows NT 6.1; en-US) AppleWebKit/534.16 (KHTML, like Gecko) Chrome/10.0.648.133 Safari/534.16", "Mozilla/5.0 (Windows NT 6.1; WOW64; Trident/7.0; rv:11.0) like Gecko", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/536.11 (KHTML, like Gecko) Chrome/20.0.1132.11 TaoBrowser/2.0 Safari/536.11", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Maxthon/4.4.3.4000 Chrome/30.0.1599.101 Safari/537.36", "Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 5.1; Trident/4.0; SV1; QQDownload 732; .NET4.0C; .NET4.0E; SE 2.X MetaSr 1.0)", "Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1; SV1; QQDownload 732; .NET4.0C; .NET4.0E; LBBROWSER)", "Mozilla/5.0 (Windows; U; Windows NT 6.1; en-us) AppleWebKit/534.50 (KHTML, like Gecko) Version/5.1 Safari/534.50", "Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.1; Trident/5.0", "Opera/9.80 (Windows NT 6.1; U; en) Presto/2.8.131 Version/11.11", "Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 5.1; TencentTraveler 4.0)"}
	//ProgramVersion add comment
	ProgramVersion = "1.0." + fmt.Sprintf("%d%d%d%d%d\nhttps://github.com/zsdevX/DarkEye\n大橘Oo0\n84500316@qq.com",
		time.Now().Year(), time.Now().Month(), time.Now().Day(), time.Now().Hour(), time.Now().Minute())
	//Banner add comment
	Banner = `
██████╗  ██╗ ██████╗  ██████╗ █████╗ ████████╗
██╔══██╗███║██╔════╝ ██╔════╝██╔══██╗╚══██╔══╝
██████╔╝╚██║██║  ███╗██║     ███████║   ██║   
██╔══██╗ ██║██║   ██║██║     ██╔══██║   ██║   
██████╔╝ ██║╚██████╔╝╚██████╗██║  ██║   ██║   
╚═════╝  ╚═╝ ╚═════╝  ╚═════╝╚═╝  ╚═╝   ╚═╝
` + "\n" + ProgramVersion + "\n"
)

var (
	//PortList add comment
	PortList = "21,22,23,80-89,389,443,445,512,513,514,873,1025,111,1433,1521,2082,2083,2222,2601,2604,3128,3306,3312,3311,3389,4430-4445,5432,5900,5984,6082,6379,7001,7002,7778,8000-9090,8080,8083,8649,8888,9200,9300,10000,11211,27017,27018,28017,50000,50070,50030,554,53,110,1080"
)

var (
	//LowCaseAlpha add comment
	LowCaseAlpha = "abcedfghijklmnopqrstuvwxyz"
)

func init() {
	BaseDir, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	logFile = filepath.Join(BaseDir, logFile)
}

//GenHumanSecond add comment
func GenHumanSecond(base int) int {
	return base
}

//TrimUseless add comment
func TrimUseless(a string) string {
	//a = strings.Replace(a, " ", " ", -1)
	a = strings.Replace(a, "\n", "", -1)
	a = strings.Replace(a, "\r", "", -1)
	a = strings.Replace(a, ",", " ", -1)
	return a
}

//StopIt add comment
func StopIt(stop *int32) {
	atomic.StoreInt32(stop, 1)
}

//StartIt add comment
func StartIt(stop *int32) {
	atomic.StoreInt32(stop, 0)
}

//ShouldStop add comment
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

//ISUtf8 add comment
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
