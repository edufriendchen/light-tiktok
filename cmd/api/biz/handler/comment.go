package handler

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/edufriendchen/light-tiktok/cmd/api/biz/rpc"
	"github.com/edufriendchen/light-tiktok/kitex_gen/comment"
	"github.com/edufriendchen/light-tiktok/kitex_gen/favorite"
	"github.com/edufriendchen/light-tiktok/kitex_gen/message"
	"github.com/edufriendchen/light-tiktok/pkg/errno"
)

func ActionComment(ctx context.Context, c *app.RequestContext) {
	var req comment.ActionRequest
	err := c.BindAndValidate(&req)
	if err != nil {
		c.JSON(consts.StatusOK, &favorite.FavoriteResponse{StatusCode: errno.ParamErr.ErrCode, StatusMsg: &errno.ParamErr.ErrMsg})
		return
	}
	resp, err := rpc.ActionComment(ctx, &req)
	if err != nil {
		Err := errno.ConvertErr(err)
		c.JSON(consts.StatusOK, &message.ChatResponse{StatusCode: Err.ErrCode, StatusMsg: &Err.ErrMsg})
		return
	}
	SetResponse(c, resp)
	return
}

func MGetCommentList(ctx context.Context, c *app.RequestContext) {
	var req comment.CommentRequest
	err := c.BindAndValidate(&req)
	if err != nil {
		c.JSON(consts.StatusOK, &favorite.FavoriteResponse{StatusCode: errno.ParamErr.ErrCode, StatusMsg: &errno.ParamErr.ErrMsg})
		return
	}
	resp, err := rpc.MGetComment(ctx, &req)
	if err != nil {
		Err := errno.ConvertErr(err)
		c.JSON(consts.StatusOK, &message.ChatResponse{StatusCode: Err.ErrCode, StatusMsg: &Err.ErrMsg})
		return
	}
	SetResponse(c, resp)
	return
}
