package service

import (
	"context"

	"github.com/edufriendchen/light-tiktok/cmd/feed/dal"
	"github.com/edufriendchen/light-tiktok/kitex_gen/feed"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type MGetFeedListService struct {
	ctx     context.Context
	session neo4j.SessionWithContext
}

func NewMGetFeedListService(ctx context.Context, driver neo4j.DriverWithContext) *MGetFeedListService {
	session := driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)
	return &MGetFeedListService{ctx: ctx, session: session}
}

func (s *MGetFeedListService) MGetFeedList(req *feed.FeedRequest, limit int64, user_id int64) ([]*feed.Video, error) {
	list, err := dal.MGetPublishListLimit(s.ctx, s.session, req, limit, user_id)
	if err != nil {
		return nil, err
	}
	return list, nil
}
