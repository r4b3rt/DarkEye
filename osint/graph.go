package main

import (
	"encoding/json"
	"fmt"
	"github.com/zsdevX/DarkEye/osint/graph"
	"github.com/zsdevX/DarkEye/osint/social"
)

func (g *OsInt) updateProfile(so social.Social) error {
	profile, _ := so.Profile(nil)
	if profile == nil {
		return fmt.Errorf("无profile信息")
	}
	_, err := g.Graph.AdjustItem(graph.AddNode, map[string]interface{}{
		"id":     profile.IDStr,
		"social": profile.Social,
		"name":   profile.Name,
	})
	if err != nil {
		return err
	}
	return nil
}

func (g *OsInt) updateFollower(so social.Social) error {
	follower, _ := so.Follower(nil)
	if follower == nil {
		return nil
	}
	list := make([]social.Profile, 0)
	if err := json.Unmarshal([]byte(follower.Profiles), &list); err != nil {
		return err
	}
	for _, f := range list {
		_, err := g.Graph.AdjustItem(graph.AddNode, map[string]interface{}{
			"id":     f.IDStr,
			"social": f.Social,
			"name":   f.Name,
		})
		if err != nil {
			return err
		}
		if _, err := g.Graph.AdjustItem(graph.AddRelation, map[string]interface{}{
			"id1":       f.IDStr,
			"id2":       follower.IDStr,
			"social":    f.Social,
			"condition": `n1.social = $social AND n1.id_str = $id1 AND n2.social = $social AND n2.id_str = $id2`,
		}); err != nil {
			return err
		}
	}
	return nil
}

func (g *OsInt) updateFollow(so social.Social) error {
	follow, _ := so.Follow(nil)
	if follow == nil {
		return nil
	}
	list := make([]social.Profile, 0)
	if err := json.Unmarshal([]byte(follow.Profiles), &list); err != nil {
		return err
	}
	for _, f := range list {
		_, err := g.Graph.AdjustItem(graph.AddNode, map[string]interface{}{
			"id":     f.IDStr,
			"social": f.Social,
			"name":   f.Name,
		})
		if err != nil {
			return err
		}
		if _, err := g.Graph.AdjustItem(graph.AddRelation, map[string]interface{}{
			"id1":       follow.IDStr,
			"id2":       f.IDStr,
			"social":    f.Social,
			"condition": `n1.social = $social AND n1.id_str = $id1 AND n2.social = $social AND n2.id_str = $id2`,
		}); err != nil {
			return err
		}
	}
	return nil
}
