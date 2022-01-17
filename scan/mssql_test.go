package scan

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_mssql(t *testing.T) {
	s, err := New(Mssql, 100)
	assert.Equal(t, nil, err)
	l := logrus.New()
	l.SetLevel(logrus.DebugLevel)
	s.Setup(l, []string{"test","root"}, []string{"xx"})
	assert.Equal(t, true, s.(*mssqlConf).Identify(context.Background(), "187.123.94.116", "1433"))
	assert.Equal(t, false, s.(*mssqlConf).Identify(context.Background(), "192.168.1.1", "80"))
	_, err = s.(*mssqlConf).Start(context.Background(), "187.123.94.116", "1433")
}
