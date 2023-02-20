package dal

import (
	"context"
	"fmt"

	"github.com/edufriendchen/light-tiktok/kitex_gen/feed"
	"github.com/edufriendchen/light-tiktok/kitex_gen/user"
	"github.com/edufriendchen/light-tiktok/pkg/errno"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

const (
	createVideoCYPHER       = `MATCH (u:User) WHERE id(u) = $user_id MERGE (u)-[p:publish]->(v:Video { play_url: $play_url, cover_url: $cover_url, favorite_count: $favorite_count, comment_count: $comment_count, title: $title } ) SET u.work_count = u.work_count + 1 RETURN p`
	queryPublishByUidCYPHER = `MATCH (u:User), (v:Video) WHERE id(u) = $user_id AND (u)-[:publish]->(v) RETURN u,v`
)

func CreateVideo(ctx context.Context, session neo4j.SessionWithContext, user_id int64, play_url string, cover_url string, title string) (int64, error) {
	video_id, err := neo4j.ExecuteWrite[int64](ctx, session, func(tx neo4j.ManagedTransaction) (int64, error) {
		result, err := tx.Run(ctx, createVideoCYPHER, map[string]any{
			"user_id":        user_id,
			"play_url":       play_url,
			"cover_url":      cover_url,
			"favorite_count": 0,
			"comment_count":  0,
			"title":          title,
		})
		if err != nil {
			return 0, err
		}
		record, err := result.Single(ctx)
		if err != nil {
			return 0, err
		}
		rawPerson, found := record.Get("p")
		if !found {
			return 0, errno.Neo4jColumnFailedErr
		}
		itemNode, ok := rawPerson.(neo4j.Relationship)
		if !ok {
			return 0, fmt.Errorf("expected result to be a map but was %T", rawPerson)
		}
		video_id := itemNode.GetId()
		return video_id, nil
	})
	return video_id, err
}

func MGetPublishListById(ctx context.Context, session neo4j.SessionWithContext, user_id int64) ([]*feed.Video, error) {
	var videoList []*feed.Video
	var videoItem feed.Video
	var userItem user.User
	_, err := neo4j.ExecuteRead[bool](ctx, session, func(tx neo4j.ManagedTransaction) (bool, error) {
		result, err := tx.Run(ctx, queryPublishByUidCYPHER, map[string]any{
			"user_id": user_id,
		})
		if err != nil {
			return false, err
		}
		record, err := result.Collect(ctx)
		if err != nil {
			return false, err
		}
		for i := 0; i < len(record); i++ {
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
				userItem.BackgroundImage, err = neo4j.GetProperty[string](itemNode, "background_image")
				if err != nil {
					return false, err
				}
				userItem.Signature, err = neo4j.GetProperty[string](itemNode, "signature")
				if err != nil {
					return false, err
				}
				total_favorited, err := neo4j.GetProperty[int64](itemNode, "total_favorited")
				if err != nil {
					return false, err
				}
				work_count, err := neo4j.GetProperty[int64](itemNode, "work_count")
				if err != nil {
					return false, err
				}
				favorite_count, err := neo4j.GetProperty[int64](itemNode, "favorite_count")
				if err != nil {
					return false, err
				}
				userItem.FollowCount = follow_count
				userItem.FollowerCount = follower_count
				userItem.TotalFavorited = total_favorited
				userItem.WorkCount = work_count
				userItem.FavoriteCount = favorite_count
			}
			videoItem.Author = &userItem
			videoList = append(videoList, &videoItem)
		}
		return true, nil
	})
	return videoList, err
}
