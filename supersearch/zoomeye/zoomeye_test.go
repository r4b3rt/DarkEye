package zoomeye

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	testAuthCode = ""
)

func Test_login(t *testing.T) {
	auth, err := login(context.Background(), testAuthCode)
	fmt.Println(auth, err)
	assert.Equal(t, nil, err)
	assert.Equal(t, true, auth != "")
}

func Test_resource(t *testing.T) {
	auth, err := login(context.Background(), testAuthCode)
	assert.Equal(t, err, nil)
	key := testAuthCode
	if auth != "" {
		key = ""
	}
	log := make(chan string, 10)
	err = resource(context.Background(), auth, key, log)
	assert.Equal(t, err, nil)
	l := <-log
	fmt.Println(l)
}