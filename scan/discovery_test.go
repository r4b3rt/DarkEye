package scan

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_pingWithPrivileged(t *testing.T) {
	s, err := New(Discovery, 100, "ping")
	assert.Equal(t, nil, err)
	//assert.Equal(t, true, s.(*discovery).pingWithPrivileged(context.Background(), "127.0.0.1") == nil)
	assert.Equal(t, true, s.(*discovery).ping(context.Background(), "127.0.0.1"))
}
