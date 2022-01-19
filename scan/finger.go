package scan

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/b1gcat/DarkEye/dict"
	"github.com/sirupsen/logrus"
	"github.com/twmb/murmur3"
	"regexp"
)

var (
	faviconHash = make(map[string]string, 0)
	faviconRe   = regexp.MustCompile(`href="(.*?favicon[^"]*.[i|p][c|n][o|g])"`)
)

type cmsFinger struct {
	Cms     string   `json:"cms"`
	Method  string   `json:"method"`
	Keyword []string `json:"keyword"`
}

type Finger struct {
	Cms []cmsFinger `json:"fingerprint"`
}

func init() {
	favicon, err := dict.Asset("finger.json")
	if err != nil {
		logrus.Error("faviconHash.init:", err.Error())
		return
	}

	fingers := Finger{
		Cms: make([]cmsFinger, 0),
	}
	if err = json.Unmarshal(favicon, &fingers); err != nil {
		if err != nil {
			logrus.Error("faviconHash.init:", err.Error())
			return
		}
	}

	for _, v := range fingers.Cms {
		if v.Method != "faviconhash" {
			continue
		}
		for _, k := range v.Keyword {
			faviconHash[k] = v.Cms
		}
	}
}

func getFavicon(body []byte) [][]string {
	return faviconRe.FindAllStringSubmatch(string(body), -1)
}

func standBase64(braw []byte) []byte {
	bckd := base64.StdEncoding.EncodeToString(braw)
	var buffer bytes.Buffer
	for i := 0; i < len(bckd); i++ {
		ch := bckd[i]
		buffer.WriteByte(ch)
		if (i+1)%76 == 0 {
			buffer.WriteByte('\n')
		}
	}
	buffer.WriteByte('\n')
	return buffer.Bytes()

}

func mmh3Hash32(raw []byte) (fh string) {
	var h32 = murmur3.New32()
	_, err := h32.Write(raw)
	if err == nil {
		fh = fmt.Sprintf("%d", int32(h32.Sum32()))
	}
	return
}
