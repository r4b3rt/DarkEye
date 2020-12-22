package fofa

//IpNode add comment
type IpNode struct {
	Ip     string
	Domain string
	Title  string
	Server string
	Finger string
	Port   string
	Alive  int
}

//Fofa add comment
type Fofa struct {
	Interval    int    `json:"interval"`
	Ip          string `json:"ip"`
	FofaSession string `json:"fofa_session"`
	Pages       int    `json:"pages"`

	IpNodes    []IpNode    `json:"-"`
	ErrChannel chan string `json:"-"`
	Stop       int32       `json:"-"`
}

//NewConfig add comment
func NewConfig() Fofa {
	return Fofa{
		Interval:    3,
		Pages:       5,
		FofaSession: "_fofapro_ars_session=Your-Cookie",
	}
}
