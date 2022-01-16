package scan

import (
	"context"
	"fmt"
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
