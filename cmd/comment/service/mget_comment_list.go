package service

import (
	"context"

	"github.com/edufriendchen/light-tiktok/cmd/comment/dal"
	"github.com/edufriendchen/light-tiktok/kitex_gen/comment"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type MGetCommentService struct {
	ctx     context.Context
	session neo4j.SessionWithContext
}

func NewMGetCommentService(ctx context.Context, driver neo4j.DriverWithContext) *MGetCommentService {
	session := driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)
	return &MGetCommentService{ctx: ctx, session: session}
}

func (s *MGetCommentService) MGetCommentList(req *comment.CommentRequest) ([]*comment.Comment, error) {
	resp, err := dal.MGetCommentListById(s.ctx, s.session, req.VideoId)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
