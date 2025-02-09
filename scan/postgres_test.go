package scan

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_postgres(t *testing.T) {
	s, err := New(Postgres, 100)
	l := logrus.New()
	l.SetLevel(logrus.DebugLevel)
	s.Setup(l, []string{"test","root"}, []string{"xx"})
	assert.Equal(t, nil, err)
	assert.Equal(t, false, s.(*postgresConf).Identify(context.Background(), "192.168.1.1", "80"))
	assert.Equal(t, true, s.(*postgresConf).Identify(context.Background(), "206.189.221.147", "5432"))
	r, err := s.(*postgresConf).Start(context.Background(), "206.189.221.147", "5432")
	fmt.Println(r, err)
}
