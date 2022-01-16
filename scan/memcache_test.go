package scan

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_memCache(t *testing.T) {
	l := logrus.New()
	l.SetLevel(logrus.DebugLevel)
	s, err := New(Memcached, 1000,
		[]string{"test", "2"}, []string{"%user%", "1@opq"}, l)
	assert.Equal(t, nil, err)
	assert.Equal(t, true, s.(*memCacheConf).Identify(context.Background(), "212.18.225.219", "11211"))
	assert.Equal(t, false, s.(*memCacheConf).Identify(context.Background(), "192.168.1.1", "80"))
}
