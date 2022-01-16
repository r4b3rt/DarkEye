package scan

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_mongodb(t *testing.T) {
	l := logrus.New()
	l.SetLevel(logrus.DebugLevel)
	s, _ := New(Mongodb, 1000,
		[]string{"test", "admin"}, []string{"123456"}, l)

	//assert.Equal(t, nil, err)
	assert.Equal(t, true, s.(*mongodbConf).Identify(context.Background(), "47.", "27017"))
	assert.Equal(t, false, s.(*mongodbConf).Identify(context.Background(), "192.168.1.1", "80"))
	r, _ := s.(*mongodbConf).Start(context.Background(), "47.", "27017")
	fmt.Println(r)
}
