package superscan

import (
	"context"
	"github.com/zsdevX/DarkEye/superscan/plugins"
	_ "net/http/pprof"
	"testing"
)

func Test_Scan(t *testing.T) {
	s := New("192.168.1.1")
	s.PortRange = "80"
	s.ActivePort = "0"
	s.Parent = context.TODO()
	plugins.Config.ParentCtx = s.Parent
	s.Run()
}
