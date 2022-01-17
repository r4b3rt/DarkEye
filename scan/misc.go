package scan

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"strings"
)

func weakPass(parent context.Context, s, addr string, ul, pl []string,
	cb func(context.Context, string, string, string) bool) (interface{}, error) {
	ctx, _ := context.WithCancel(parent)
	for _, u := range ul {
		for _, p := range pl {
			p = strings.ReplaceAll(p, `%user%`, u)
			select {
			case <-ctx.Done():
				return nil, nil
			default:
			}
			if cb(ctx, addr, u, p) {
				return fmt.Sprintf("%s %s - %s/%s", s, addr, u, p), nil
			}
		}
	}
	return nil, nil
}

func setupRisk(r *risk, args ...interface{}) {
	for k, v := range args {
		switch v.(type) {
		case *logrus.Logger:
			r.logger = v.(*logrus.Logger)
		case []string:
			switch k {
			case 1:
				if x, ok := v.([]string); ok {
					r.username = x
				}
			case 2:
				if x, ok := v.([]string); ok {
					r.password = x
				}
			}
		}
	}
}

type risk struct {
	username []string
	password []string
	logger   *logrus.Logger
}
