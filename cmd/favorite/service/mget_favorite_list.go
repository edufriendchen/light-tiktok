package service

import (
	"context"

	"github.com/edufriendchen/light-tiktok/cmd/favorite/dal"
	"github.com/edufriendchen/light-tiktok/kitex_gen/feed"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type MGetFavoriteListService struct {
	ctx     context.Context
	session neo4j.SessionWithContext
}

func NewMGetFavoriteListService(ctx context.Context, driver neo4j.DriverWithContext) *MGetFavoriteListService {
	session := driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)
	return &MGetFavoriteListService{ctx: ctx, session: session}
}

func (s *MGetFavoriteListService) MGetFavoriteList(user_id int64) ([]*feed.Video, error) {
	list, err := dal.MGetFavoriteListById(s.ctx, s.session, user_id)
	if err != nil {
		return nil, err
	}
	return list, nil
}
