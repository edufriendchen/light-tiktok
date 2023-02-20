package handler

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/edufriendchen/light-tiktok/cmd/api/biz/rpc"
	"github.com/edufriendchen/light-tiktok/kitex_gen/user"
	"github.com/edufriendchen/light-tiktok/pkg/errno"
)

func Register(ctx context.Context, c *app.RequestContext) {
	var req user.CreateUserRequest
	err := c.BindAndValidate(&req)
	if err != nil {
		c.JSON(consts.StatusOK, &user.CreateUserResponse{StatusCode: errno.ParamErr.ErrCode, StatusMsg: &errno.ParamErr.ErrMsg})
		return
	}
	resp, err := rpc.CreateUser(ctx, &req)
	if err != nil {
		Err := errno.ConvertErr(err)
		c.JSON(consts.StatusOK, &user.CreateUserResponse{StatusCode: Err.ErrCode, StatusMsg: &Err.ErrMsg})
		return
	}
	SetResponse(c, resp)
	return
}

func Login(ctx context.Context, c *app.RequestContext) {
	var req user.CheckUserRequest
	err := c.BindAndValidate(&req)
	if err != nil {
		SetResponse(c, &user.CheckUserResponse{StatusCode: errno.ParamErr.ErrCode, StatusMsg: &errno.ParamErr.ErrMsg})
		return
	}
	resp, err := rpc.CheckUse(context.Background(), &req)
	if err != nil {
		SetResponse(c, &user.CheckUserResponse{StatusCode: errno.ServiceErr.ErrCode, StatusMsg: &errno.ServiceErr.ErrMsg})
		return
	}
	SetResponse(c, resp)
	return
}

func MGetUserInfo(ctx context.Context, c *app.RequestContext) {
	var req user.MGetUserRequest
	err := c.BindAndValidate(&req)
	if err != nil {
		SetResponse(c, &user.MGetUserResponse{StatusCode: errno.ParamErr.ErrCode, StatusMsg: &errno.ParamErr.ErrMsg})
		return
	}
	resp, err := rpc.MGetUserInfo(ctx, &req)
	if err != nil {
		SetResponse(c, &user.MGetUserResponse{StatusCode: errno.ServiceErr.ErrCode, StatusMsg: &errno.ServiceErr.ErrMsg})
		return
	}
	SetResponse(c, resp)
	return
}
