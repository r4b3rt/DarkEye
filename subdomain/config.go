package subdomain

type ipInfo struct {
	ip         string
	Isp        string `json:"isp"`
	RegionName string `json:"regionName"`
}

type dnsInfo struct {
	domain string
	ip     []ipInfo
	cname  string
	//
	title  string
	server string
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
	dns        []dnsInfo
	Stop       int32 `json:"-"`
}

func NewConfig() SubDomain {
	return SubDomain{
		IpCheck:     true,
		Brute:       false,
		BruteLength: "3",
		BruteRate:   "50",
	}
}
