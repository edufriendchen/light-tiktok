package service

import (
	"context"
	"fmt"

	"github.com/edufriendchen/light-tiktok/cmd/favorite/dal"
	"github.com/edufriendchen/light-tiktok/kitex_gen/favorite"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type ActionFavoriteService struct {
	ctx     context.Context
	session neo4j.SessionWithContext
}

func NewActionFavoriteService(ctx context.Context, driver neo4j.DriverWithContext) *ActionFavoriteService {
	session := driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)
	return &ActionFavoriteService{ctx: ctx, session: session}
}

func (s *ActionFavoriteService) ActionFavoriteService(req *favorite.ActionRequest, user_id int64) error {
	is, err := dal.IsFavorite(s.ctx, s.session, user_id, req.VideoId)
	if err != nil {
		return err
	}
	if (req.ActionType == 1 && is) || (req.ActionType == 2 && !is) {
		return nil
	}
	if req.ActionType == 1 && !is {
		if _, err := dal.DoFavorite(s.ctx, s.session, user_id, req.VideoId); err != nil {
			fmt.Println("DoFavorite-Err:", err)
			return err
		}
	}
	if req.ActionType == 2 && is {
		fmt.Println("UnFavorite")
		if _, err := dal.UnFavorite(s.ctx, s.session, user_id, req.VideoId); err != nil {
			fmt.Println("UnFavorite-Err:", err)
			return err
		}
	}
	return nil
}
