package subdomain

type IpInfo struct {
	Ip         string
	Isp        string `json:"isp"`
	RegionName string `json:"regionName"`
}

type DnsInfo struct {
	Domain string
	Ip     []IpInfo
	Cname  string
	//
	Title  string
	Server string
	Code   int32
}

type SubDomain struct {
	Queries     string `json:"queries"`
	ApiKey      string `json:"api_key"`
	DnsServer   string `json:"dns_server"`
	IpCheck     bool   `json:"ip_check"`
	Brute       bool
	BruteRate   string
	BruteLength string

	ErrChannel chan string `json:"-"`
	Dns        []DnsInfo   `json:"-"`
	Stop       int32       `json:"-"`
}

func NewConfig() SubDomain {
	return SubDomain{
		IpCheck:     true,
		Brute:       false,
		BruteLength: "3",
		BruteRate:   "50",
	}
}
