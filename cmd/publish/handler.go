package main

import (
	"context"

	"github.com/edufriendchen/light-tiktok/cmd/publish/service"
	"github.com/edufriendchen/light-tiktok/kitex_gen/publish"
	"github.com/edufriendchen/light-tiktok/pkg/errno"
	"github.com/edufriendchen/light-tiktok/pkg/global"
)

type PublishServiceImpl struct{}

func (s *PublishServiceImpl) ActionPulish(ctx context.Context, req *publish.ActionRequest) (resp *publish.ActionResponse, err error) {
	resp = new(publish.ActionResponse)
	if err = req.IsValid(); err != nil {
		resp = &publish.ActionResponse{StatusCode: errno.ParamErr.ErrCode, StatusMsg: &errno.ParamErr.ErrMsg}
		return resp, err
	}
	claims, err := global.Jwt.ParseToken(req.Token)
	if err != nil {
		resp = &publish.ActionResponse{StatusCode: errno.AuthorizationFailedErr.ErrCode, StatusMsg: &errno.AuthorizationFailedErr.ErrMsg}
		return resp, err
	}
	err = service.NewActionPulishService(ctx, global.Neo4jDriver).ActionPulish(req, claims.Id)
	if err != nil {
		resp = &publish.ActionResponse{StatusCode: errno.ServiceErr.ErrCode, StatusMsg: &errno.ServiceErr.ErrMsg}
		return resp, err
	}
	resp = &publish.ActionResponse{StatusCode: errno.Success.ErrCode, StatusMsg: &errno.Success.ErrMsg}
	return resp, nil
}

func (s *PublishServiceImpl) MGetPublishList(ctx context.Context, req *publish.PublishRequest) (resp *publish.PublishResponse, err error) {
	resp = new(publish.PublishResponse)
	if err = req.IsValid(); err != nil {
		resp = &publish.PublishResponse{StatusCode: errno.ParamErr.ErrCode, StatusMsg: &errno.ParamErr.ErrMsg}
		return resp, err
	}
	claims, err := global.Jwt.ParseToken(req.Token)
	if err != nil {
		resp = &publish.PublishResponse{StatusCode: errno.AuthorizationFailedErr.ErrCode, StatusMsg: &errno.AuthorizationFailedErr.ErrMsg}
		return resp, err
	}
	list, err := service.NewActionPulishService(ctx, global.Neo4jDriver).MGetPulishList(claims.Id)
	if err != nil {
		resp = &publish.PublishResponse{StatusCode: errno.ServiceErr.ErrCode, StatusMsg: &errno.ServiceErr.ErrMsg}
		return resp, err
	}
	resp = &publish.PublishResponse{StatusCode: errno.Success.ErrCode, StatusMsg: &errno.Success.ErrMsg}
	resp.VideoList = list
	return resp, nil
}
