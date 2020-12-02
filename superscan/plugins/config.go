package plugins

import "sync"

type Web struct {
	Server string `json:",omitempty"`
	Title  string `json:",omitempty"`
}

type Account struct {
	Username string `json:",omitempty"`
	Password string `json:",omitempty"`
}

type MSBulletin struct {
	os string `json:",omitempty"`
	//ms17010
	Description string `json:",omitempty"`
}

type Plugins struct {
	PortOpened bool
	Web        Web        `json:",omitempty"`
	SSh        []Account  `json:",omitempty"`
	Ms17010    MSBulletin `json:",omitempty"`
	Mysql      []Account  `json:",omitempty"`

	NoTrust        bool `json:",omitempty"`
	Worker         int  `json:",omitempty"`
	TargetIp       string
	TargetPort     string
	TargetProtocol string
	TimeOut        int `json:"-"`
	DescCallback   func(string)
	highLight      bool
	locker         sync.RWMutex
}

const (
	SSHSrv = iota
	MysqlSrv
	MS17010
	WEBSrv  //放到最后
	PluginNR
)

const (
	OKNone = iota
	OKWait
	OKTimeOut
	OKDone
	OKNext
	OKStop
)
