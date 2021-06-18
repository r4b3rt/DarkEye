package plugins

import (
	"encoding/json"
	"fmt"
	"github.com/orcaman/concurrent-map"
	"testing"
)

func Test_crack(t *testing.T) {
	s := new(Service)
	*s = services["ssh"]
	s.thread = 1
	s.parent = &Plugins{
		TargetIp:   "",
		TargetPort: "22",
		Result:Result{
			Output:cmap.New(),
		},
	}
	Config.UserList = []string{"root"}
	Config.PassList = []string{""}
	Config.TimeOut = 3000


	s.check(s)
	b, _ := json.MarshalIndent(&s.parent.Result, "", "	")
	fmt.Println(string(b))
}
