package handler

import (
	"context"
	"fmt"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/edufriendchen/light-tiktok/cmd/api/biz/rpc"
	"github.com/edufriendchen/light-tiktok/kitex_gen/relation"
	"github.com/edufriendchen/light-tiktok/pkg/errno"
)

func RelationAction(ctx context.Context, c *app.RequestContext) {
	var req relation.ActionRequest
	err := c.BindAndValidate(&req)
	if err != nil {
		c.JSON(consts.StatusOK, &relation.ActionResponse{StatusCode: errno.ParamErr.ErrCode, StatusMsg: &errno.ParamErr.ErrMsg})
		return
	}
	resp, err := rpc.ActionRelation(ctx, &req)
	if err != nil {
		Err := errno.ConvertErr(err)
		c.JSON(consts.StatusOK, &relation.ActionResponse{StatusCode: Err.ErrCode, StatusMsg: &Err.ErrMsg})
		return
	}
	SetResponse(c, resp)
	return
}

func MGetFollowList(ctx context.Context, c *app.RequestContext) {
	var req relation.FollowRequest
	fmt.Println("req", req)
	fmt.Println("Query", c.Query("user_id"))
	err := c.BindAndValidate(&req)
	fmt.Println("req", req)
	if err != nil {
		SetResponse(c, &relation.FollowResponse{StatusCode: errno.ParamErr.ErrCode, StatusMsg: &errno.ParamErr.ErrMsg})
		return
	}
	resp, err := rpc.MGetFollowList(ctx, &req)
	if err != nil {
		SetResponse(c, &relation.FollowResponse{StatusCode: errno.ServiceErr.ErrCode, StatusMsg: &errno.ServiceErr.ErrMsg})
		return
	}
	SetResponse(c, resp)
	return
}

// MGetFollowerList 注册用户操作 的上下文至 User 服务的 RPC 客户端, 并获取相应的响应.
func MGetFollowerList(ctx context.Context, c *app.RequestContext) {
	var req relation.FollowerRequest
	err := c.BindAndValidate(&req)
	if err != nil {
		SetResponse(c, &relation.FollowerResponse{StatusCode: errno.ParamErr.ErrCode, StatusMsg: &errno.ParamErr.ErrMsg})
		return
	}
	resp, err := rpc.MGetFollowerList(ctx, &req)
	if err != nil {
		SetResponse(c, &relation.FollowerResponse{StatusCode: errno.ServiceErr.ErrCode, StatusMsg: &errno.ServiceErr.ErrMsg})
		return
	}
	SetResponse(c, resp)
	return
}

// MGetUserInfo 注册用户操作 的上下文至 User 服务的 RPC 客户端, 并获取相应的响应.
func MGetFriendList(ctx context.Context, c *app.RequestContext) {
	var req relation.FriendRequest
	err := c.BindAndValidate(&req)
	if err != nil {
		SetResponse(c, &relation.FriendResponse{StatusCode: errno.ParamErr.ErrCode, StatusMsg: &errno.ParamErr.ErrMsg})
		return
	}
	resp, err := rpc.MGetFriendList(ctx, &req)
	if err != nil {
		SetResponse(c, &relation.FriendResponse{StatusCode: errno.ServiceErr.ErrCode, StatusMsg: &errno.ServiceErr.ErrMsg})
		return
	}
	SetResponse(c, resp)
	return
}
