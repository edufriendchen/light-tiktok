package main

import (
	"context"

	service "github.com/edufriendchen/light-tiktok/cmd/relation/service"
	relation "github.com/edufriendchen/light-tiktok/kitex_gen/relation"
	"github.com/edufriendchen/light-tiktok/pkg/errno"
	"github.com/edufriendchen/light-tiktok/pkg/global"
)

// RelationServiceImpl implements the last service interface defined in the IDL.
type RelationServiceImpl struct{}

// ActionRelation implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) ActionRelation(ctx context.Context, req *relation.ActionRequest) (resp *relation.ActionResponse, err error) {
	resp = new(relation.ActionResponse)
	if err = req.IsValid(); err != nil {
		resp = &relation.ActionResponse{StatusCode: errno.ParamErr.ErrCode, StatusMsg: &errno.ParamErr.ErrMsg}
		return resp, err
	}
	claims, err := global.Jwt.ParseToken(req.Token)
	if err != nil {
		resp = &relation.ActionResponse{StatusCode: errno.AuthorizationFailedErr.ErrCode, StatusMsg: &errno.AuthorizationFailedErr.ErrMsg}
		return resp, nil
	}
	_, err = service.NewActionRelationService(ctx, global.Neo4jDriver).ActionRelation(req, claims.Id)
	if err != nil {
		resp = &relation.ActionResponse{StatusCode: errno.ServiceErr.ErrCode, StatusMsg: &errno.ServiceErr.ErrMsg}
		return resp, err
	}
	resp = &relation.ActionResponse{StatusCode: errno.Success.ErrCode, StatusMsg: &errno.Success.ErrMsg}
	return resp, nil
}

// MGetFollowList implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) MGetFollowList(ctx context.Context, req *relation.FollowRequest) (resp *relation.FollowResponse, err error) {
	resp = new(relation.FollowResponse)
	if err = req.IsValid(); err != nil {
		resp = &relation.FollowResponse{StatusCode: errno.ParamErr.ErrCode, StatusMsg: &errno.ParamErr.ErrMsg}
		return resp, err
	}
	claims, err := global.Jwt.ParseToken(req.Token)
	if err != nil {
		resp = &relation.FollowResponse{StatusCode: errno.AuthorizationFailedErr.ErrCode, StatusMsg: &errno.AuthorizationFailedErr.ErrMsg}
		return resp, nil
	}
	req.UserId = claims.Id
	list, err := service.NewGetFollowListService(ctx, global.Neo4jDriver).GetFollowList(req)
	if err != nil {
		resp = &relation.FollowResponse{StatusCode: errno.ServiceErr.ErrCode, StatusMsg: &errno.ServiceErr.ErrMsg}
		return resp, err
	}
	resp = &relation.FollowResponse{StatusCode: errno.Success.ErrCode, StatusMsg: &errno.Success.ErrMsg}
	resp.UserList = list
	return resp, nil
}

// MGetFollowerList implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) MGetFollowerList(ctx context.Context, req *relation.FollowerRequest) (resp *relation.FollowerResponse, err error) {
	resp = new(relation.FollowerResponse)
	if err = req.IsValid(); err != nil {
		resp = &relation.FollowerResponse{StatusCode: errno.ParamErr.ErrCode, StatusMsg: &errno.ParamErr.ErrMsg}
		return resp, err
	}
	claims, err := global.Jwt.ParseToken(req.Token)
	if err != nil {
		resp = &relation.FollowerResponse{StatusCode: errno.AuthorizationFailedErr.ErrCode, StatusMsg: &errno.AuthorizationFailedErr.ErrMsg}
		return resp, nil
	}
	req.UserId = claims.Id
	list, err := service.NewGetFollowerListService(ctx, global.Neo4jDriver).GetFollowerList(req)
	if err != nil {
		resp = &relation.FollowerResponse{StatusCode: errno.ServiceErr.ErrCode, StatusMsg: &errno.ServiceErr.ErrMsg}
		return resp, err
	}
	resp = &relation.FollowerResponse{StatusCode: errno.Success.ErrCode, StatusMsg: &errno.Success.ErrMsg}
	resp.UserList = list
	return resp, nil
}

// MGetFriendList implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) MGetFriendList(ctx context.Context, req *relation.FriendRequest) (resp *relation.FriendResponse, err error) {
	resp = new(relation.FriendResponse)
	if err = req.IsValid(); err != nil {
		resp = &relation.FriendResponse{StatusCode: errno.ParamErr.ErrCode, StatusMsg: &errno.ParamErr.ErrMsg}
		return resp, err
	}
	claims, err := global.Jwt.ParseToken(req.Token)
	if err != nil {
		resp = &relation.FriendResponse{StatusCode: errno.AuthorizationFailedErr.ErrCode, StatusMsg: &errno.AuthorizationFailedErr.ErrMsg}
		return resp, nil
	}
	req.UserId = claims.Id
	list, err := service.NewGetFriendListService(ctx, global.Neo4jDriver).GetFriendList(req)
	if err != nil {
		resp = &relation.FriendResponse{StatusCode: errno.ServiceErr.ErrCode, StatusMsg: &errno.ServiceErr.ErrMsg}
		return resp, err
	}
	resp = &relation.FriendResponse{StatusCode: errno.Success.ErrCode, StatusMsg: &errno.Success.ErrMsg}
	resp.UserList = list
	return resp, nil
}
