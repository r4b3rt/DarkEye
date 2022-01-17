package scan

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_mongodb(t *testing.T) {
	s, _ := New(Mongodb, 100)
	l := logrus.New()
	l.SetLevel(logrus.DebugLevel)
	s.Setup([]string{"test","kali"}, []string{"%user%"}, l)
	assert.Equal(t, true, s.(*mongodbConf).Identify(context.Background(), "47.", "27017"))
	assert.Equal(t, false, s.(*mongodbConf).Identify(context.Background(), "192.168.1.1", "80"))
	r, _ := s.(*mongodbConf).Start(context.Background(), "47.", "27017")
	fmt.Println(r)
}
