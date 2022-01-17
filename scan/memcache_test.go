package scan

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_memCache(t *testing.T) {
	s, err := New(Memcached, 1000)
	assert.Equal(t, nil, err)
	l := logrus.New()
	l.SetLevel(logrus.DebugLevel)
	s.Setup(l, []string{"test","root"}, []string{"xx"})
	assert.Equal(t, true, s.(*memCacheConf).Identify(context.Background(), "212.18.225.219", "11211"))
	assert.Equal(t, false, s.(*memCacheConf).Identify(context.Background(), "192.168.1.1", "80"))
}
