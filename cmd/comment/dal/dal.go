package dal

import (
	"context"
	"time"

	"github.com/edufriendchen/light-tiktok/cmd/user/dal"
	"github.com/edufriendchen/light-tiktok/kitex_gen/comment"
	"github.com/edufriendchen/light-tiktok/kitex_gen/user"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

const (
	createCYPHER         = `MATCH (u:User) WHERE id(u) = $user_id MATCH (v:Video) WHERE id(v) = $video_id SET v.comment_count = v.comment_count + 1 CREATE (u)-[c:comment{comment_text: $comment_text, create_date: $create_date }]->(v) RETURN id(c) as id, u`
	getByVideoIDCYPHER   = `MATCH (u:User),(v:Video), (u)-[c:comment]->(v) WHERE id(v) = $video_id  RETURN u,c`
	getByCommentIDCYPHER = `MATCH (u)-[c:comment]->(v) WHERE id(c) = $comment_id RETURN u,c`
	deleteCYPHER         = `MATCH (u:User)-[c:comment]->(v:Video) WHERE id(c) = $comment_id SET v.comment_count = v.comment_count - 1 DELETE c`
)

func CreateComment(ctx context.Context, session neo4j.SessionWithContext, req *comment.ActionRequest, user_id int64) (*comment.Comment, error) {
	var user *user.User
	nowTime := time.Now().Format("01-02")
	comment, err := neo4j.ExecuteWrite[*comment.Comment](ctx, session, func(tx neo4j.ManagedTransaction) (*comment.Comment, error) {
		result, err := tx.Run(ctx, createCYPHER, map[string]any{
			"user_id":      user_id,
			"video_id":     req.VideoId,
			"comment_text": req.CommentText,
			"create_date":  nowTime,
		})
		if err != nil {
			return nil, err
		}
		record, err := result.Single(ctx)
		if err != nil {
			return nil, err
		}
		comment_id := record.Values[0].(int64)
		user_value, ok := record.Get("u")
		if ok {
			itemNode := user_value.(neo4j.Node)
			user, err = dal.GetUserInNode(itemNode)
			if err != nil {
				return nil, err
			}
		}
		return &comment.Comment{
			Id:         comment_id,
			User:       user,
			Content:    req.CommentText,
			CreateDate: nowTime,
		}, nil
	})
	return comment, err
}

func MGetCommentListById(ctx context.Context, session neo4j.SessionWithContext, video_id int64) ([]*comment.Comment, error) {
	var commentList []*comment.Comment
	var user *user.User
	_, err := neo4j.ExecuteRead[bool](ctx, session, func(tx neo4j.ManagedTransaction) (bool, error) {
		result, err := tx.Run(ctx, getByVideoIDCYPHER, map[string]any{
			"video_id": video_id,
		})
		if err != nil {
			return false, err
		}
		record, err := result.Collect(ctx)
		if err != nil {
			return false, err
		}

		for i := 0; i < len(record); i++ {
			var commentItem comment.Comment
			video_value, ok := record[i].Get("c")
			if ok {
				itemRelationship := video_value.(neo4j.Relationship)
				commentItem.Id = itemRelationship.GetId()
				if commentItem.Content, err = neo4j.GetProperty[string](itemRelationship, "comment_text"); err != nil {
					return false, err
				}
				if commentItem.CreateDate, err = neo4j.GetProperty[string](itemRelationship, "create_date"); err != nil {
					return false, err
				}
			}
			user_value, ok := record[i].Get("u")
			if ok {
				itemNode := user_value.(neo4j.Node)
				user, err = dal.GetUserInNode(itemNode)
				if err != nil {
					return false, err
				}
			}
			commentItem.User = user
			commentList = append(commentList, &commentItem)
		}
		return true, nil
	})
	return commentList, err
}

func DeleteComment(ctx context.Context, session neo4j.SessionWithContext, comment_id int64) (*comment.Comment, error) {
	var comment comment.Comment
	var user *user.User
	_, err := neo4j.ExecuteRead[bool](ctx, session, func(tx neo4j.ManagedTransaction) (bool, error) {
		result, err := tx.Run(ctx, getByCommentIDCYPHER, map[string]any{
			"comment_id": comment_id,
		})
		if err != nil {
			return false, err
		}
		record, err := result.Single(ctx)
		if err != nil {
			return false, err
		}
		user_value, ok := record.Get("u")
		if ok {
			itemNode := user_value.(neo4j.Node)
			user, err = dal.GetUserInNode(itemNode)
			if err != nil {
				return false, err
			}
		}
		comment_value, ok := record.Get("c")
		if ok {
			itemNode := comment_value.(neo4j.Relationship)
			if comment.Content, err = neo4j.GetProperty[string](itemNode, "comment_text"); err != nil {
				return false, err
			}
			if comment.CreateDate, err = neo4j.GetProperty[string](itemNode, "create_date"); err != nil {
				return false, err
			}
		}
		return true, nil
	})
	if err != nil {
		return nil, err
	}
	comment.Id = comment_id
	comment.User = user
	summary, err := neo4j.ExecuteWrite[neo4j.ResultSummary](ctx, session, func(tx neo4j.ManagedTransaction) (neo4j.ResultSummary, error) {
		result, err := tx.Run(ctx, deleteCYPHER, map[string]any{
			"comment_id": comment_id,
		})
		if err != nil {
			return nil, err
		}
		summary, _ := result.Consume(ctx)
		return summary, nil
	})
	if err != nil && summary.Counters().PropertiesSet() == 0 {
		return nil, err
	}
	return &comment, err
}
