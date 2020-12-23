package spider

//SensitiveInterface add comment
type SensitiveInterface struct {
	API   string
	Level int
}

//Spider add comment
type Spider struct {
	//spider
	Url               string `json:"url"`
	DisAllowedRequest string `json:"disallow_request"`
	RequestMatchRule  string `json:"request_match_rule"`
	ResponseMatchRule string `json:"response_match_rule"`
	MaxDeps           int    `json:"max_deps"`
	LocalLink         bool   `json:"local_link"`
	Cookie            string `json:"cookie"`
	//search
	Query        string `json:"query"`
	SearchAPIKey string `json:"search_api_key"`
	SearchEnable bool   `json:"search_enable"`

	//
	ErrChannel         chan string `json:"-"`
	Stop               int32       `json:"-"`
	sensitiveInterface []SensitiveInterface
}

//NewConfig add comment
func NewConfig() Spider {
	return Spider{
		LocalLink: true,
		MaxDeps:   2,
		//禁止spider去访问这些后缀的资源
		DisAllowedRequest: `(\.png|\.jpg|\.jpeg|\.bmp|\.zip|\.rar|\.gz|\.tar|\.swf|\.flv|\.mp4|\.avi|\.ico)`,
		//定义spider的链接
		RequestMatchRule: "a:href,script:src",
		//获取资源数据后匹配对应的数据
		//[\"|\'][/]?[a-zA-Z]+[/]+[a-zA-Z]+.*[\"|\'] => "login/abc"
		//(https?|ftp|file)://[-A-Za-z0-9+&@#/%?=~_|!:,.;]+[-A-Za-z0-9+&@#/%=~_|] => http://a.b
		ResponseMatchRule: `(\"|\')[-A-Za-z0-9+&@#/%?=~_|!:,.;]+[-A-Za-z0-9+&@#/%=~_|](\"|\')|(https?|ftp|file)://[-A-Za-z0-9+&@#/%?=~_|!:,.;]+[-A-Za-z0-9+&@#/%=~_|]`,
		Query:             "site:ooxx.com filetype:txt",
	}
}
