package handler

import (
	"context"
	"fmt"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/edufriendchen/light-tiktok/cmd/api/biz/rpc"
	"github.com/edufriendchen/light-tiktok/kitex_gen/favorite"
	"github.com/edufriendchen/light-tiktok/pkg/errno"
)

func ActionFavorite(ctx context.Context, c *app.RequestContext) {
	var req favorite.ActionRequest
	err := c.BindAndValidate(&req)
	fmt.Println("params:", req)
	if err != nil {
		c.JSON(consts.StatusOK, &favorite.FavoriteResponse{StatusCode: errno.ParamErr.ErrCode, StatusMsg: &errno.ParamErr.ErrMsg})
		return
	}
	resp, err := rpc.ActionFavorite(ctx, &req)
	if err != nil {
		c.JSON(consts.StatusOK, &favorite.FavoriteResponse{StatusCode: errno.ServiceErr.ErrCode, StatusMsg: &errno.ServiceErr.ErrMsg})
		return
	}

	SetResponse(c, resp)
	return
}

func MGetFavoriteList(ctx context.Context, c *app.RequestContext) {
	var req favorite.FavoriteRequest
	err := c.BindAndValidate(&req)
	if err != nil {
		c.JSON(consts.StatusOK, &favorite.FavoriteResponse{StatusCode: errno.ParamErr.ErrCode, StatusMsg: &errno.ParamErr.ErrMsg})
		return
	}
	resp, err := rpc.MGetFavoriteList(ctx, &req)
	if err != nil {
		Err := errno.ConvertErr(err)
		c.JSON(consts.StatusOK, &favorite.FavoriteResponse{StatusCode: Err.ErrCode, StatusMsg: &Err.ErrMsg})
		return
	}
	SetResponse(c, resp)
	return
}
