package spider

import "encoding/base64"

type Link struct {
}

type Spider struct {
	Url               string `json:"url"`
	DisAllowedRequest string `json:"disallow_request"`
	RequestMatchRule  string `json:"request_match_rule"`
	ResponseMatchRule string `json:"response_match_rule"`
	ResponseFilter    string `json:"response_filter_rule"`
	MaxDeps           int    `json:"max_deps"`
	LocalLink         bool   `json:"local_link"`
	Cookie            string `json:"cookie"`

	links      []Link
	ErrChannel chan string `json:"-"`
	Stop       int32       `json:"-"`
}

func NewConfig() Spider {
	return Spider{
		LocalLink:         true,
		MaxDeps:           2,
		//禁止spider去访问这些后缀的资源
		DisAllowedRequest: base64.StdEncoding.EncodeToString([]byte(`(\.png|\.jpg|\.jpeg|\.bmp|\.zip|\.rar|\.gz|\.tar|\.swf|\.flv|\.mp4|\.avi|\.ico)`)),
		//定义spider的链接
		RequestMatchRule:  "a:href,script:src",
		//获取资源数据后匹配对应的数据
		//[\"|\'][/]?[a-zA-Z]+[/]+[a-zA-Z]+.*[\"|\'] => "login/abc"
		//(https?|ftp|file)://[-A-Za-z0-9+&@#/%?=~_|!:,.;]+[-A-Za-z0-9+&@#/%=~_|] => http://a.b
		ResponseMatchRule: base64.StdEncoding.EncodeToString([]byte(`(\"|\')[-A-Za-z0-9+&@#/%?=~_|!:,.;]+[-A-Za-z0-9+&@#/%=~_|](\"|\')|(https?|ftp|file)://[-A-Za-z0-9+&@#/%?=~_|!:,.;]+[-A-Za-z0-9+&@#/%=~_|]`)),
		//对提取的数据过滤，去掉匹配的数据
		ResponseFilter:    base64.StdEncoding.EncodeToString([]byte(`(text/css|text/javascript|text/css|\.css|text/html)`)),
	}
}
