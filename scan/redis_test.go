package scan

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func Test_redis(t *testing.T) {
	l := logrus.New()
	l.SetLevel(logrus.DebugLevel)
	s, err := New(Redis, 5000,
		[]string{"test", ""}, []string{"%user%"}, l)
	assert.Equal(t, nil, err)
	assert.Equal(t, true, s.(*redisConf).Identify(context.Background(), "vulfocus.fofa.so", "28262"))
	assert.Equal(t, false, s.(*redisConf).Identify(context.Background(), "192.168.1.1", "80"))

	r, err := s.(*redisConf).Start(context.Background(), "vulfocus.fofa.so", "28262")
	assert.Equal(t, true, strings.Contains(r.(string), "redis" ))
}
