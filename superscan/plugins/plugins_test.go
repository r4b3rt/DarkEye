package plugins

import (
	"encoding/json"
	"fmt"
	"testing"
)

func test_crack(t *testing.T) {
	s := new(Service)
	*s = services["rdp"]
	s.thread = 1
	s.parent = &Plugins{
		TargetIp:   "127.0.0.1",
		TargetPort: "3389",
	}
	Config.UserList = []string{"zs"}
	Config.PassList = []string{"111111"}

	s.check(s)
	b, _ := json.MarshalIndent(&s.parent.Result, "", "	")
	fmt.Println(string(b))
}

func test_redis(t *testing.T) {
	s := new(Service)
	*s = services["redis"]
	s.thread = 1
	s.parent = &Plugins{
		TargetIp:   "",
		TargetPort: "6379",
	}
	s.check(s)
	b, _ := json.MarshalIndent(&s.parent.Result, "", "	")
	fmt.Println(string(b))
}

func test_netbios(t *testing.T) {
	s := new(Service)
	*s = preServices["netbios"]
	s.thread = 1
	s.parent = &Plugins{
		TargetIp:   "192.168.1.11",
		TargetPort: "137",
	}
	s.check(s)
	b, _ := json.MarshalIndent(&s.parent.Result, "", "	")
	fmt.Println(string(b))
}
