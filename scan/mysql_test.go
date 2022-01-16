package scan

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_mysql(t *testing.T) {
	l := logrus.New()
	l.SetLevel(logrus.DebugLevel)
	s, err := New(Mysql, 100,
		[]string{"test","kali"}, []string{"%user%"}, l)
	assert.Equal(t, nil, err)
	assert.Equal(t, false,s.(*mysqlConf).Identify(context.Background(), "192.168.1.1","80"))
	assert.Equal(t, true,s.(*mysqlConf).Identify(context.Background(), "219.148.38.137","3306"))
	r, err := s.(*mysqlConf).Start(context.Background(), "219.148.38.137","3306")
	fmt.Println(r, err)


}
