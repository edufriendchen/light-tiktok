package service

import (
	"context"
	"fmt"

	"github.com/edufriendchen/light-tiktok/cmd/comment/dal"
	"github.com/edufriendchen/light-tiktok/kitex_gen/comment"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type ActionCommentService struct {
	ctx     context.Context
	session neo4j.SessionWithContext
}

func NewActionCommentService(ctx context.Context, driver neo4j.DriverWithContext) *ActionCommentService {
	session := driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)
	return &ActionCommentService{ctx: ctx, session: session}
}

func (s *ActionCommentService) ActionComment(req *comment.ActionRequest, user_id int64, comment_id int64) (*comment.Comment, error) {
	if req.ActionType == 1 {
		resp, err := dal.CreateComment(s.ctx, s.session, req, user_id)
		if err != nil {
			fmt.Println("CreateComment-Err:", err)
			return nil, err
		}
		return resp, nil
	}
	if req.ActionType == 2 {
		if _, err := dal.DeleteComment(s.ctx, s.session, comment_id); err != nil {
			fmt.Println("DeleteComment-Err:", err)
			return nil, err
		}
	}
	return nil, nil
}
