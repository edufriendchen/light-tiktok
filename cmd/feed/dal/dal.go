package dal

import (
	"context"
	"fmt"

	favoritedal "github.com/edufriendchen/light-tiktok/cmd/favorite/dal"
	"github.com/edufriendchen/light-tiktok/kitex_gen/feed"
	"github.com/edufriendchen/light-tiktok/kitex_gen/user"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

const (
	getByLimitCYPHER = `MATCH (u:User), (v:Video) WHERE (u)-[:publish]->(v) RETURN v,u LIMIT $num`
	hasFollowCYPHER  = `MATCH (u:User) WHERE id(u) = $user_id MATCH (t:User) WHERE id(t) = $to_user_id RETURN CASE WHEN (u)-[:follow]->(t) THEN true ELSE false END AS result`
)

func MGetPublishListLimit(ctx context.Context, session neo4j.SessionWithContext, req *feed.FeedRequest, limit int64, user_id int64) ([]*feed.Video, error) {

	var videoList []*feed.Video
	_, err := neo4j.ExecuteRead[bool](ctx, session, func(tx neo4j.ManagedTransaction) (bool, error) {
		result, err := tx.Run(ctx, getByLimitCYPHER, map[string]any{
			"num": limit,
		})
		if err != nil {
			return false, err
		}
		record, err := result.Collect(ctx)
		if err != nil {
			return false, err
		}
		for i := 0; i < len(record); i++ {
			var videoItem feed.Video
			var userItem user.User
			video_value, ok := record[i].Get("v")
			if ok {
				itemNode := video_value.(neo4j.Node)
				videoItem.Id = itemNode.GetId()
				if videoItem.Title, err = neo4j.GetProperty[string](itemNode, "title"); err != nil {
					return false, err
				}
				if videoItem.CoverUrl, err = neo4j.GetProperty[string](itemNode, "cover_url"); err != nil {
					return false, err
				}
				if videoItem.PlayUrl, err = neo4j.GetProperty[string](itemNode, "play_url"); err != nil {
					return false, err
				}
				if videoItem.CommentCount, err = neo4j.GetProperty[int64](itemNode, "comment_count"); err != nil {
					return false, err
				}
				if videoItem.FavoriteCount, err = neo4j.GetProperty[int64](itemNode, "favorite_count"); err != nil {
					return false, err
				}
			}
			user_value, ok := record[i].Get("u")
			if ok {
				itemNode := user_value.(neo4j.Node)
				userItem.Id = itemNode.GetId()
				if userItem.Name, err = neo4j.GetProperty[string](itemNode, "nickname"); err != nil {
					return false, err
				}
				if userItem.Avatar, err = neo4j.GetProperty[string](itemNode, "avatar"); err != nil {
					return false, err
				}
				follow_count, err := neo4j.GetProperty[int64](itemNode, "follow_count")
				if err != nil {
					return false, err
				}
				follower_count, err := neo4j.GetProperty[int64](itemNode, "follower_count")
				if err != nil {
					return false, err
				}
				userItem.FollowCount = follow_count
				userItem.FollowerCount = follower_count
			}
			if user_id != 0 && user_id != userItem.Id {
				userItem.IsFollow, err = HasFollow(ctx, session, user_id, userItem.Id)
				if err != nil {
					fmt.Println("判断是否本人与该视频作者是否存在关注关系错误！")
				}
				videoItem.IsFavorite, err = favoritedal.IsFavorite(ctx, session, user_id, videoItem.Id)
				if err != nil {
					fmt.Println("判断是否本人与该视频是否存在点赞关系错误！")
				}
			}
			videoItem.Author = &userItem
			videoList = append(videoList, &videoItem)
		}
		return true, nil
	})
	return videoList, err
}

func HasFollow(ctx context.Context, session neo4j.SessionWithContext, user_id int64, to_user_id int64) (bool, error) {
	var params = map[string]any{
		"user_id":    user_id,
		"to_user_id": to_user_id,
	}
	hasRelation, err := neo4j.ExecuteRead[bool](ctx, session, func(tx neo4j.ManagedTransaction) (bool, error) {
		result, err := tx.Run(ctx, hasFollowCYPHER, params)
		if err != nil {
			return false, err
		}
		record, err := result.Single(ctx)
		if err != nil {
			return false, err
		}
		return record.Values[0].(bool), nil
	})
	return hasRelation, err
}
