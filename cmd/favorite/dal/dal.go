package dal

import (
	"context"
	"fmt"

	"github.com/edufriendchen/light-tiktok/cmd/user/dal"
	"github.com/edufriendchen/light-tiktok/kitex_gen/feed"
	"github.com/edufriendchen/light-tiktok/kitex_gen/user"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

const (
	createFavariteCYPHER     = `MATCH (u:User) WHERE id(u) = $user_id MATCH (v:Video) WHERE id(v) = $video_id MERGE (u)-[f:favorite]->(v) SET v.favorite_count = v.favorite_count + 1, RETURN id(f) as id`
	queryFavariteByUIDCYPHER = `MATCH (u:User), (v:Video) WHERE id(u) = $user_id AND (u)-[:favorite]->(v) RETURN v,u`
)

func MGetFavoriteListById(ctx context.Context, session neo4j.SessionWithContext, user_id int64) ([]*feed.Video, error) {
	var videoList []*feed.Video
	var videoItem feed.Video
	var userItem *user.User
	_, err := neo4j.ExecuteRead[bool](ctx, session, func(tx neo4j.ManagedTransaction) (bool, error) {
		result, err := tx.Run(ctx, queryFavariteByUIDCYPHER, map[string]any{
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
				userItem, err = dal.GetUserInNode(itemNode)
			}
			videoItem.Author = userItem
			videoList = append(videoList, &videoItem)
		}
		return true, nil
	})
	return videoList, err
}

func IsFavorite(ctx context.Context, session neo4j.SessionWithContext, user_id int64, video_id int64) (bool, error) {
	is, err := neo4j.ExecuteRead[bool](ctx, session, func(tx neo4j.ManagedTransaction) (bool, error) {
		result, err := tx.Run(ctx, "MATCH (u:User) WHERE id(u) = $user_id MATCH (v:Video) WHERE id(v) = $video_id MATCH (u)-[f:favorite]->(v) WITH COUNT(f) > 0 as is_favorite RETURN is_favorite", map[string]any{
			"user_id":  user_id,
			"video_id": video_id,
		})
		if err != nil {
			return false, err
		}
		record, err := result.Single(ctx)
		if err != nil {
			return false, err
		}
		fmt.Println("is_favorite:", record.Values)
		return record.Values[0].(bool), nil
	})
	return is, err
}

func DoFavorite(ctx context.Context, session neo4j.SessionWithContext, user_id int64, video_id int64) (int64, error) {
	fmt.Println("DoFavorite")
	relationship_id, err := neo4j.ExecuteWrite[int64](ctx, session, func(tx neo4j.ManagedTransaction) (int64, error) {
		result, err := tx.Run(ctx, createFavariteCYPHER, map[string]any{
			"user_id":  user_id,
			"video_id": video_id,
		})
		if err != nil {
			return 0, err
		}
		record, err := result.Single(ctx)
		if err != nil {
			return 0, err
		}
		return record.Values[0].(int64), nil
	})
	return relationship_id, err
}

func UnFavorite(ctx context.Context, session neo4j.SessionWithContext, user_id int64, video_id int64) (int64, error) {
	relationship_id, err := neo4j.ExecuteWrite[int64](ctx, session, func(tx neo4j.ManagedTransaction) (int64, error) {
		result, err := tx.Run(ctx, "MATCH (u:User) WHERE id(u) = $user_id MATCH (v:Video) WHERE id(v) = $video_id MATCH (u)-[f:favorite]->(v) SET v.favorite_count = v.favorite_count - 1 DELETE f RETURN id(f)", map[string]any{
			"user_id":  user_id,
			"video_id": video_id,
		})
		if err != nil {
			return 0, err
		}
		record, err := result.Single(ctx)
		if err != nil {
			return 0, err
		}
		return record.Values[0].(int64), nil
	})
	return relationship_id, err
}
