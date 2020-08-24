package securitytrails

type ipInfo struct {
	ip         string
	Isp        string `json:"isp"`
	RegionName string `json:"regionName"`
}

type dnsInfo struct {
	domain string
	ip     []ipInfo
	cname  string
}

type SecurityTrails struct {
	Queries   string `json:"queries"`
	ApiKey    string `json:"api_key"`
	DnsServer string `json:"dns_server"`
	IpCheck   bool   `json:"ip_check"`
	dns       []dnsInfo

	ErrChannel chan string `json:"-"`
	Stop       int32       `json:"-"`
}
