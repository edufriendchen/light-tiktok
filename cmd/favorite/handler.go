package main

import (
	"context"
	"fmt"

	"github.com/edufriendchen/light-tiktok/cmd/favorite/service"
	"github.com/edufriendchen/light-tiktok/kitex_gen/favorite"
	"github.com/edufriendchen/light-tiktok/pkg/errno"
	"github.com/edufriendchen/light-tiktok/pkg/global"
)

// FavoriteServiceImpl implements the last service interface defined in the IDL.
type FavoriteServiceImpl struct{}

// ActionFavorite implements the FavoriteServiceImpl interface.
func (s *FavoriteServiceImpl) ActionFavorite(ctx context.Context, req *favorite.ActionRequest) (resp *favorite.ActionResponse, err error) {
	// TODO: Your code here...
	resp = new(favorite.ActionResponse)
	if err = req.IsValid(); err != nil {
		resp = &favorite.ActionResponse{StatusCode: errno.ParamErr.ErrCode, StatusMsg: &errno.ParamErr.ErrMsg}
		return resp, err
	}
	claims, err := global.Jwt.ParseToken(req.Token)
	if err != nil {
		resp = &favorite.ActionResponse{StatusCode: errno.AuthorizationFailedErr.ErrCode, StatusMsg: &errno.AuthorizationFailedErr.ErrMsg}
		return resp, err
	}
	err = service.NewActionFavoriteService(ctx, global.Neo4jDriver).ActionFavoriteService(req, claims.Id)
	if err != nil {
		resp = &favorite.ActionResponse{StatusCode: errno.ServiceErr.ErrCode, StatusMsg: &errno.ServiceErr.ErrMsg}
		return resp, err
	}
	resp = &favorite.ActionResponse{StatusCode: errno.Success.ErrCode, StatusMsg: &errno.Success.ErrMsg}
	fmt.Println("resp:", resp)
	return resp, nil
}

// MGetFavoriteList implements the FavoriteServiceImpl interface.
func (s *FavoriteServiceImpl) MGetFavoriteList(ctx context.Context, req *favorite.FavoriteRequest) (resp *favorite.FavoriteResponse, err error) {
	// TODO: Your code here...
	resp = new(favorite.FavoriteResponse)
	if err = req.IsValid(); err != nil {
		resp = &favorite.FavoriteResponse{StatusCode: errno.ParamErr.ErrCode, StatusMsg: &errno.ParamErr.ErrMsg}
		return resp, err
	}
	claims, err := global.Jwt.ParseToken(req.Token)
	if err != nil {
		resp = &favorite.FavoriteResponse{StatusCode: errno.AuthorizationFailedErr.ErrCode, StatusMsg: &errno.AuthorizationFailedErr.ErrMsg}
		return resp, err
	}
	list, err := service.NewMGetFavoriteListService(ctx, global.Neo4jDriver).MGetFavoriteList(claims.Id)
	if err != nil {
		resp = &favorite.FavoriteResponse{StatusCode: errno.ServiceErr.ErrCode, StatusMsg: &errno.ServiceErr.ErrMsg}
		return resp, err
	}
	resp = &favorite.FavoriteResponse{StatusCode: errno.Success.ErrCode, StatusMsg: &errno.Success.ErrMsg}
	resp.VideoList = list
	return resp, nil
}
