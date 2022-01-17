package scan

import (
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"testing"
)

func Test_ftp(t *testing.T) {
	s, err := New(Ftp, 100)
	assert.Equal(t, nil, err)
	l := logrus.New()
	l.SetLevel(logrus.DebugLevel)
	s.Setup(l, []string{"test","root"}, []string{"xx"})
	assert.Equal(t, true, s.(*ftpConf).Identify(context.Background(), "154.220.26.25", "21"))
	assert.Equal(t, false, s.(*ftpConf).Identify(context.Background(), "192.168.1.1", "80"))
	_, err = s.(*ftpConf).Start(context.Background(), "154.220.26.25", "21")
}
