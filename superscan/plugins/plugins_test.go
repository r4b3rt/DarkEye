package plugins

import (
	"encoding/json"
	"fmt"
	"testing"
)

func Test_crack(t *testing.T) {
	s := new(Service)
	*s = services["ftp"]
	s.thread = 1
	s.parent = &Plugins{
		TargetIp: "192.168.1.1",
		TargetPort: "21",
	}
	Config.UserList = []string{"varbing"}
	Config.PassList = []string{"varbing@123@woshitiancai"}

	s.check(s)
	b, _ := json.MarshalIndent(&s.parent.Result,"","	")
	fmt.Println(string(b))
}
