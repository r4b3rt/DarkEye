package scan

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func Test_pingWithPrivileged(t *testing.T) {
	s, err := New(DiscoPing, 100)
	assert.Equal(t, nil, err)
	//assert.Equal(t, true, s.(*discovery).pingWithPrivileged(context.Background(), "127.0.0.1") == nil)
	assert.Equal(t, true, s.(*discovery).ping(context.Background(), "127.0.0.1"))
}

func Test_http(t *testing.T) {
	s, err := New(DiscoHttp, 2000)
	assert.Equal(t, nil, err)
	s.(*discovery).logger.SetLevel(logrus.DebugLevel)

	s.(*discovery).host = strings.Split("www.baidu.com,www.wtf.com",",")
	r, err := s.(*discovery).Start(context.Background(), "192.168.1.1", "80")
	f, _:= json.MarshalIndent(r, ""," ")
	fmt.Println(string(f))
}