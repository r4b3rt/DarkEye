package plugins

type Web struct {
	Server string `json:",omitempty"`
	Title  string `json:",omitempty"`
}

type SSH struct {
	Username string `json:",omitempty"`
	Password string `json:",omitempty"`
}

type Plugins struct {
	PortOpened bool
	Web        Web   `json:",omitempty"`
	SSh        []SSH `json:",omitempty"`

	TargetIp       string
	TargetPort     string
	TargetProtocol string
	TimeOut        int `json:"-"`
}

const (
	WEBSrv = iota
	SSHSrv
	PluginNR
)
