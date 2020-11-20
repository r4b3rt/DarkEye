package fofa

type ipNode struct {
	Ip     string
	Domain string
	Title  string
	Server string
	Finger string
	Port   string
	Alive  int
}

type Fofa struct {
	Interval    int    `json:"interval"`
	Ip          string `json:"ip"`
	FofaSession string `json:"fofa_session"`
	Pages       int    `json:"pages"`

	ipNodes    []ipNode
	ErrChannel chan string `json:"-"`
	Stop       int32       `json:"-"`
}

func NewConfig() Fofa {
	return Fofa{
		Interval:    6,
		Pages:       5,
		FofaSession: "_fofapro_ars_session=Your-Cookie",
	}
}
