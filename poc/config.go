package poc

type PocResult struct {
	PocName string
	Url     string
}

type Poc struct {
	FileName string
	Urls     string
	//
	ReverseUrl string
	//http://api.ceye.io/v1/records?token={token}&type={dns|http}&filter={filter}
	ReverseCheckUrl  string
	ReverseUseDomain bool

	ErrChannel chan string `json:"-"`
	Stop       int32       `json:"-"`
	Results    []PocResult `json:"-"`
}

func NewConfig() Poc {
	return Poc{
		ReverseUrl:      "qvn0kc.ceye.io",
		ReverseCheckUrl: "http://api.ceye.io/v1/records?token=066f3d242991929c823ac85bb60f4313&type=http&filter=",
	}
}
