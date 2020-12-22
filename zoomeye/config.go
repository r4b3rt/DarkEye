package zoomeye

type ZoomEye struct {
	ApiKey string `json:"api-key"`
	Query  string `json:"query"`
	Pages  string `json:"-"`

	Results    Results     `json:"-"`
	ErrChannel chan string `json:"-"`
	Stop       int32       `json:"-"`
}

type Results struct {
	Error   string
	Message string
	Url     string

	Total     int          `json:"total"`
	Available int          `json:"available"`
	Matches   []TargetInfo `json:"matches"`
}

type TargetInfo struct {
	Ip           string       `json:"ip"`
	City         City         `json:"city"`
	Country      Country      `json:"country"`
	Subdivisions Subdivisions `json:"subdivisions"`
	PortInfo     PortInfo     `json:"portinfo"`
	Protocol     Protocol     `json:"protocol"`
	Honeypot     int          `json:"honeypot"`
}

type City struct {
	Names []string `json:"names"`
}

type Country struct {
	Names []string `json:"names"`
}

type Subdivisions struct {
	Names []string `json:"names"`
}

type PortInfo struct {
	Port     int    `json:"port"`
	Title    string `json:"title"`
	Banner   string `json:"banner"`
	Service  string `json:"service"`
	Hostname string `json:"hostname"`
	Device   string `json:"device"`
	Os       string `json:"os"`
	App      string `json:"app"`
}

type Protocol struct {
	Application string `json:"application"`
	Transport   string `json:"transport"`
	Probe       string `json:"probe"`
}

func New() ZoomEye {
	return ZoomEye{

	}
}
