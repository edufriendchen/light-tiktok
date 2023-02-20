package main

import (
	"context"

	"github.com/edufriendchen/light-tiktok/pkg/errno"
	"github.com/edufriendchen/light-tiktok/pkg/global"

	service "github.com/edufriendchen/light-tiktok/cmd/user/service"

	"github.com/edufriendchen/light-tiktok/kitex_gen/user"
	"github.com/edufriendchen/light-tiktok/pkg/jwt"
)

// UserServiceImpl implements the last service interface defined in the IDL.
type UserServiceImpl struct{}

// CreateUser implements the UserServiceImpl interface.
func (s *UserServiceImpl) CreateUser(ctx context.Context, req *user.CreateUserRequest) (resp *user.CreateUserResponse, err error) {
	resp = new(user.CreateUserResponse)
	if err = req.IsValid(); err != nil {
		resp = &user.CreateUserResponse{StatusCode: errno.ParamErr.ErrCode, StatusMsg: &errno.ParamErr.ErrMsg}
		return resp, err
	}
	userid, err := service.NewCreateUserService(ctx, global.Neo4jDriver).CreateUserNode(req)
	if err != nil {
		resp = &user.CreateUserResponse{StatusCode: errno.ServiceErr.ErrCode, StatusMsg: &errno.ServiceErr.ErrMsg}
		return resp, err
	}
	token, err := global.Jwt.CreateToken(jwt.CustomClaims{
		Id: userid,
	})
	resp = &user.CreateUserResponse{StatusCode: errno.Success.ErrCode, StatusMsg: &errno.Success.ErrMsg, UserId: userid, Token: token}
	return resp, nil
}

// CheckUser implements the UserServiceImpl interface.
func (s *UserServiceImpl) CheckUser(ctx context.Context, req *user.CheckUserRequest) (resp *user.CheckUserResponse, err error) {
	resp = new(user.CheckUserResponse)
	if err = req.IsValid(); err != nil {
		resp = &user.CheckUserResponse{StatusCode: errno.ParamErr.ErrCode, StatusMsg: &errno.ParamErr.ErrMsg}
		return resp, err
	}
	userid, err := service.NewCheckUserService(ctx, global.Neo4jDriver).CheckUser(req)
	if err != nil {
		resp = &user.CheckUserResponse{StatusCode: errno.ServiceErr.ErrCode, StatusMsg: &errno.ServiceErr.ErrMsg}
		return resp, err
	}
	token, err := global.Jwt.CreateToken(jwt.CustomClaims{
		Id: userid,
	})
	resp = &user.CheckUserResponse{StatusCode: errno.Success.ErrCode, StatusMsg: &errno.Success.ErrMsg, UserId: userid, Token: token}
	return resp, nil
}

// MGetUser implements the UserServiceImpl interface.
func (s *UserServiceImpl) MGetUser(ctx context.Context, req *user.MGetUserRequest) (resp *user.MGetUserResponse, err error) {
	resp = new(user.MGetUserResponse)
	if err = req.IsValid(); err != nil {
		resp = &user.MGetUserResponse{StatusCode: errno.ParamErr.ErrCode, StatusMsg: &errno.ParamErr.ErrMsg}
		return resp, err
	}
	claims, err := global.Jwt.ParseToken(req.Token)
	if err != nil {
		resp = &user.MGetUserResponse{StatusCode: errno.AuthorizationFailedErr.ErrCode, StatusMsg: &errno.AuthorizationFailedErr.ErrMsg}
		return resp, nil
	}
	if req.UserId == 0 {
		req.UserId = claims.Id
	}
	userInfo, err := service.NewGetUserService(ctx, global.Neo4jDriver).GetUser(req, claims.Id)
	if err != nil {
		resp = &user.MGetUserResponse{StatusCode: errno.ParamErr.ErrCode, StatusMsg: &errno.ParamErr.ErrMsg}
		return resp, err
	}
	resp = &user.MGetUserResponse{StatusCode: errno.Success.ErrCode, StatusMsg: &errno.Success.ErrMsg}
	resp.User = userInfo
	return resp, nil
}
