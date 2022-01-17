package scan

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_ssh(t *testing.T) {
	s, err := New(Ssh, 100)
	assert.Equal(t, nil, err)
	l := logrus.New()
	l.SetLevel(logrus.DebugLevel)
	s.Setup(l, []string{"test","root"}, []string{"xx"})
	assert.Equal(t, false,s.(*sshConf).Identify(context.Background(), "192.168.1.1","80"))
	assert.Equal(t, true,s.(*sshConf).Identify(context.Background(), "192.168.1.253","22"))
	r, err := s.(*sshConf).Start(context.Background(), "192.168.1.253","22")
	fmt.Println(r, err)
}
