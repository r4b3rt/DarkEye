package scan

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_pingWithPrivileged(t *testing.T) {
	s, err := New(DiscoPing, 100, "ping")
	assert.Equal(t, nil, err)
	//assert.Equal(t, true, s.(*discovery).pingWithPrivileged(context.Background(), "127.0.0.1") == nil)
	assert.Equal(t, true, s.(*discovery).ping(context.Background(), "127.0.0.1"))
}

func Test_http(t *testing.T) {
	s, err := New(DiscoHttp, 2000, "http")
	assert.Equal(t, nil, err)
	s.(*discovery).logger.SetLevel(logrus.DebugLevel)
	//assert.Equal(t, true, s.(*discovery).pingWithPrivileged(context.Background(), "127.0.0.1") == nil)
	r, err := s.(*discovery).Start(context.Background(), "4","443")
	fmt.Println(r, err)
}
