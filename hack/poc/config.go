package poc

type Poc struct {
	FileName string
	Urls     string

	//qvn0kc.ceye.io
	ReverseUrl string
	//http://api.ceye.io/v1/records?token={token}&type={dns|http}&filter={filter}
	ReverseCheckUrl string
	ReverseUseDomain   bool

	ErrChannel chan string `json:"-"`
	Stop       int32       `json:"-"`
}
