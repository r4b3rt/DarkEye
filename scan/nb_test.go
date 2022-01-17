package scan

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_nb(t *testing.T) {
	s, err := New(DiscoNb, 100)
	l := logrus.New()
	l.SetLevel(logrus.DebugLevel)
	s.Setup(l)
	assert.Equal(t, nil, err)
	r, err := s.(*discovery).Start(context.Background(), "192.168.1.21","0")
	fmt.Println(r, err)
}

