package common

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)


func Test_httpFinger(t *testing.T) {
	req := HttpRequest{
		Url:     "http://180.166.82.70:18081/main/login", //感谢fofa提供IP
		Method:  "GET",
	}
	response, err := req.Go()
	assert.Equal(t, err, nil)
	headerStr, err := json.Marshal(&response.ResponseHeaders)
	if err != nil {
		Log("whatWeb.Marshal", err.Error(), FAULT)
		return
	}
	assert.Equal(t, getFinger(headerStr, response.Body), "Shiro")
}
