package scan

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func Test_ssh(t *testing.T) {
	l := logrus.New()
	l.SetLevel(logrus.DebugLevel)
	s, err := New(Ssh, 100,
		[]string{"test","kali"}, []string{"%user%"}, l)
	assert.Equal(t, nil, err)
	assert.Equal(t, false,s.(*sshConf).Identify(context.Background(), "192.168.1.1","80"))
	assert.Equal(t, true,s.(*sshConf).Identify(context.Background(), "127.0.0.1","2222"))
	r, err := s.(*sshConf).Start(context.Background(), "127.0.0.1","2222")
	assert.Equal(t, true, strings.Contains(r.(string), "127.0.0.1"))
}
