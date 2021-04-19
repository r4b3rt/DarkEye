package social

import (
	"context"
	"encoding/json"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/zsdevX/DarkEye/common"
	"io/ioutil"
	"strconv"
)

var (
	twitterConfigFile = "twitter.cnf"
	twitterKeyIndex   = 0
	twitterConfig     = TwitterConfig{}
)

type Twitter struct {
	Social

	ctx       context.Context
	Prof      *Profile
	Follows   *Follow
	Followers *Follower
}

type TwitterConfig struct {
	Name string        `json:"name"`
	Keys []TwitterKeys `json:"keys"`
}

type TwitterKeys struct {
	ApiKey          string `json:"api_key"`
	ApiSecret       string `json:"api_secret"`
	ApiAccessToken  string `json:"api_access_token"`
	ApiAccessSecret string `json:"api_access_secret"`
}

func init() {
	if b, err := ioutil.ReadFile(twitterConfigFile); err == nil {
		if err = json.Unmarshal(b, &twitterConfig); err == nil {
			return
		}
	}
	common.Log("twitter.init", "读配置失败", common.ALERT)
}

func (t *Twitter) New(ctx context.Context) *Twitter {
	return &Twitter{
		ctx: ctx,
	}
}

func (t *Twitter) Profile(req *Request) (*Profile, error) {
	if t.Prof != nil {
		return t.Prof, nil
	}
	if req == nil {
		return nil, nil
	}
	client := t.client()
	id, err := strconv.ParseInt(req.IdStr, 10, 64)
	user, _, err := client.Users.Show(&twitter.UserShowParams{
		UserID:     id,
		ScreenName: req.ScreenName,
	})
	if err != nil {
		return nil, err
	}
	t.Prof = &Profile{}
	t.copy(t.Prof, user)
	return t.Prof, nil

}

func (t *Twitter) Follow(req *Request) (*Follow, error) {
	if t.Follows != nil {
		return t.Follows, nil
	}
	if req == nil || t.Prof.Follow == 0 {
		return nil, nil
	}
	client := t.client()
	id, err := strconv.ParseInt(req.IdStr, 10, 64)
	follows, _, err := client.Friends.List(&twitter.FriendListParams{
		UserID:     id,
		ScreenName: req.ScreenName,
		Count:      req.Count,
	})
	if err != nil {
		return nil, err
	}
	f := make([]Profile, 0)
	for _, v := range follows.Users {
		p := Profile{}
		t.copy(&p, &v)
		f = append(f, p)
	}
	b, err := json.Marshal(&f)
	t.Follows = &Follow{
		List{
			Social:   twitterConfig.Name,
			IDStr:    req.IdStr,
			Profiles: string(b),
		},
	}
	return t.Follows, nil
}

func (t *Twitter) Follower(req *Request) (*Follower, error) {
	if t.Followers != nil {
		return t.Followers, nil
	}
	if req == nil || t.Prof.Follower == 0 {
		return nil, nil
	}
	client := t.client()
	id, err := strconv.ParseInt(req.IdStr, 10, 64)
	followers, _, err := client.Followers.List(&twitter.FollowerListParams{
		UserID:     id,
		ScreenName: req.ScreenName,
		Count:      req.Count,
	})
	if err != nil {
		return nil, err
	}
	f := make([]Profile, 0)
	for _, v := range followers.Users {
		p := Profile{}
		t.copy(&p, &v)
		f = append(f, p)
	}
	b, err := json.Marshal(&f)
	t.Followers = &Follower{
		List{
			Social:   twitterConfig.Name,
			IDStr:    req.IdStr,
			Profiles: string(b),
		},
	}
	return t.Followers, nil
}

func (t *Twitter) client() *twitter.Client {
	config := oauth1.NewConfig(
		twitterConfig.Keys[twitterKeyIndex].ApiKey,
		twitterConfig.Keys[twitterKeyIndex].ApiSecret)
	token := oauth1.NewToken(
		twitterConfig.Keys[twitterKeyIndex].ApiAccessToken,
		twitterConfig.Keys[twitterKeyIndex].ApiAccessSecret)
	httpClient := config.Client(oauth1.NoContext, token)

	return twitter.NewClient(httpClient)
}

func (t *Twitter) copy(d *Profile, s *twitter.User) {
	d.IDStr = s.IDStr
	d.Name = s.Name
	d.ProfileBannerURL = s.ProfileBannerURL
	d.ProfileImageURL = s.ProfileImageURL
	d.Email = s.Email
	d.CreatedAt = s.CreatedAt
	d.ScreenName = s.ScreenName
	d.Description = s.Description
	d.Follow = s.FriendsCount
	d.Follower = s.FollowersCount
	d.Social = "twitter"
}
