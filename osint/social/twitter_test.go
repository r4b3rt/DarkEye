package social

import (
	"context"
	"testing"
)

func Test_profile(t *testing.T) {
	tw := &Twitter{}
	tw = tw.New(context.Background())
	req := Request{
		ScreenName: "iloveho70012979",
	}
	prof, err := tw.Profile(&req)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if prof.ScreenName != req.ScreenName {
		t.Fatalf("Bad ScreenName '%s' != '%s'", prof.ScreenName, req.ScreenName)
	}
	if _, err = tw.Follower(&req); err != nil {
		t.Fatalf(err.Error())
	}
	if _, err = tw.Follow(&req); err != nil {
		t.Fatalf(err.Error())
	}
}
