package xraypoc

import (
	"github.com/zsdevX/DarkEye/common"
	"github.com/zsdevX/DarkEye/hack/poc/xraypoc/celtypes"
	"testing"
)

func Test_ReverseCheck(t *testing.T) {
	reverse := xraypoc_celtypes.Reverse{
		Domain:          "fuckfuck." + "qvn0kc.ceye.io",
		ReverseCheckUrl: "http://api.ceye.io/v1/records?token=066f3d242991929c823ac85bb60f4313&type=http&filter=",
	}
	req := common.HttpRequest{
		Method:  "GET",
		Url:     "http://" + reverse.Domain,
		TimeOut: 10,
	}
	req.Go()

	if !myReverseCheck(&reverse, 5) {
		t.Fail()
	}
}
