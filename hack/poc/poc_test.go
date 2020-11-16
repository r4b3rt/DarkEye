package poc

import (
	"fmt"
	"testing"
)

func Test_CheckPoc(t *testing.T) {
	poc := Poc{
		//FileName: "/Users/mac/Desktop/tmp/xray/pocs",
		FileName: "testpoc/fuck.yml",
		Urls:     "http://wx.shgmc.net:8080/cmiims/a/login",
	}
	poc.ErrChannel = make(chan string, 10)

	go poc.Check()

	for {
		msg := <-poc.ErrChannel
		fmt.Println(msg)
	}
}

//go test -v -test.run Test_trySearch
func Test_trySearch(t *testing.T) {
	re := `<input type="hidden" name="csrftoken" value="(?P<token>.+?)"`
	body := []byte(`testtest<input type="hidden" name="csrftoken" value="this_is_a_token"`)
	res := trySearch(re, body)
	fmt.Println(res)
	if len(res) != 1 {
		t.Fail()
	}
}


