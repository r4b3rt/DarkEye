package scan

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_mysql(t *testing.T) {
	s, err := New(Mysql, 100)
	assert.Equal(t, nil, err)
	l := logrus.New()
	l.SetLevel(logrus.DebugLevel)
	s.Setup(l, []string{"test","root"}, []string{"xx"})
	assert.Equal(t, false,s.(*mysqlConf).Identify(context.Background(), "192.168.1.1","80"))
	assert.Equal(t, true,s.(*mysqlConf).Identify(context.Background(), "219.148.38.137","3306"))
	r, err := s.(*mysqlConf).Start(context.Background(), "219.148.38.137","3306")
	fmt.Println(r, err)
}
