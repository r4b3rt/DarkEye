package securitytrails

type ipInfo struct {
	ip         string
	Isp        string `json:"isp"`
	RegionName string `json:"regionName"`
}

type DnsHistoryValues struct {
	Ip string `json:"ip"`
}

type DnsHistory struct {
	Values        []DnsHistoryValues `json:"values"`
	Organizations []string           `json:"organizations"`
	Last_seen     string             `json:"last_seen"`
	First_seen    string             `json:"first_seen"`
}

type dnsInfo struct {
	domain string
	ip     []ipInfo
	cname  string
	history string
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
