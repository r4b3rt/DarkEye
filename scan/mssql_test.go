package scan

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_mssql(t *testing.T) {
	l := logrus.New()
	l.SetLevel(logrus.DebugLevel)
	s, err := New(Mssql, 100,
		[]string{"test", "kali"}, []string{"%user%"}, l)
	assert.Equal(t, nil, err)
	assert.Equal(t, true, s.(*mssqlConf).Identify(context.Background(), "187.123.94.116", "1433"))
	assert.Equal(t, false, s.(*mssqlConf).Identify(context.Background(), "192.168.1.1", "80"))
	_, err = s.(*mssqlConf).Start(context.Background(), "187.123.94.116", "1433")
}
