package social

import (
	"context"
	"fmt"
)

type Social interface {
	//用户资料
	Profile(req *Request) (*Profile,error)
	//粉丝和关注的人
	Follower(req *Request) (*Follower, error)
	//关注的人
	Follow(req *Request) (*Follow, error)
}

type Request struct {
	SocialName string `json:"social_name"`
	ScreenName string `json:"screen_name"`
	IdStr      string `json:"id_str"`
	Count      int    `json:"count"`
}

type Profile struct {
	CreatedAt        string `json:"created_at"`
	IDStr            string `json:"id_str"`
	Name             string `json:"name"`
	ScreenName       string `json:"screen_name"`
	ProfileImageURL  string `json:"profile_image_url"`
	ProfileBannerURL string `json:"profile_banner_image_url"`
	Description      string `json:"description"`
	Email            string `json:"email"`
	Social           string `json:"social"`

	//粉丝
	Follower int `json:"follower"`
	//关注的人数量
	Follow int `json:"follow"`
}

type Follower struct {
	List
}

type Follow struct {
	List
}

type List struct {
	IDStr    string `json:"id_str"`
	Social   string `json:"social"`
	Profiles string `json:"profiles"`
}

func New(ctx context.Context, request *Request) (Social, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	switch request.SocialName {
	case "twitter":
		return Social(Twitter{}.New(ctx)), nil
	}
	return nil, fmt.Errorf("不支持的社交类型:" + request.SocialName)
}
