package main

import (
	"context"

	"github.com/edufriendchen/light-tiktok/cmd/feed/service"
	feed "github.com/edufriendchen/light-tiktok/kitex_gen/feed"
	"github.com/edufriendchen/light-tiktok/pkg/consts"
	"github.com/edufriendchen/light-tiktok/pkg/errno"
	"github.com/edufriendchen/light-tiktok/pkg/global"
)

// FeedServiceImpl implements the last service interface defined in the IDL.
type FeedServiceImpl struct{}

// MGetFeedList implements the FeedServiceImpl interface.
func (s *FeedServiceImpl) MGetFeedList(ctx context.Context, req *feed.FeedRequest) (resp *feed.FeedResponse, err error) {
	// TODO: Your code here...
	resp = new(feed.FeedResponse)
	if err = req.IsValid(); err != nil {
		resp = &feed.FeedResponse{StatusCode: errno.ParamErr.ErrCode, StatusMsg: &errno.ParamErr.ErrMsg}
		return resp, err
	}
	var user_id int64 = 0
	if req.Token != "" {
		claims, err := global.Jwt.ParseToken(req.Token)
		if err != nil {
			resp = &feed.FeedResponse{StatusCode: errno.AuthorizationFailedErr.ErrCode, StatusMsg: &errno.AuthorizationFailedErr.ErrMsg}
			return resp, err
		}
		user_id = claims.Id
	}
	list, err := service.NewMGetFeedListService(ctx, global.Neo4jDriver).MGetFeedList(req, consts.Limit, user_id)
	if err != nil {
		resp = &feed.FeedResponse{StatusCode: errno.ServiceErr.ErrCode, StatusMsg: &errno.ServiceErr.ErrMsg}
		return resp, err
	}
	resp = &feed.FeedResponse{StatusCode: errno.Success.ErrCode, StatusMsg: &errno.Success.ErrMsg}
	resp.VideoList = list
	return resp, nil
}
