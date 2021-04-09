package zoomeye

//ZoomEye add comment
type ZoomEye struct {
	ApiKey string `json:"api-key"`
	Query  string `json:"query"`
	Pages  int    `json:"-"`

	ErrChannel chan string `json:"-"`
}

type Match struct {
	Ip      string
	Country string
	//PortInfo
	Port      int
	Os        string
	Hostname  string
	Service   string
	Banner    string
	Title     string
	Version   string
	Device    string
	ExtraInfo string
	RDns      string
	App       string
	//
	Url      string
	HttpCode int
}

//ZoomEye add comment
func New() ZoomEye {
	return ZoomEye{}
}
