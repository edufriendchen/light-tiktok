package handler

import (
	"bytes"
	"context"
	"io"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/edufriendchen/light-tiktok/cmd/api/biz/rpc"
	"github.com/edufriendchen/light-tiktok/kitex_gen/publish"
	"github.com/edufriendchen/light-tiktok/pkg/errno"
)

func PublishAction(ctx context.Context, c *app.RequestContext) {
	token := c.PostForm("token")
	title := c.PostForm("title")
	FormFile, err := c.Request.FormFile("data")
	if err != nil {
		c.JSON(consts.StatusOK, &publish.ActionResponse{StatusCode: errno.ParamErr.ErrCode, StatusMsg: &errno.ParamErr.ErrMsg})
		return
	}

	file, err := FormFile.Open()
	if err != nil {
		c.JSON(consts.StatusOK, &publish.ActionResponse{StatusCode: errno.ParamErr.ErrCode, StatusMsg: &errno.ParamErr.ErrMsg})
		return
	}
	defer file.Close()

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		c.JSON(consts.StatusOK, &publish.ActionResponse{StatusCode: errno.ParamErr.ErrCode, StatusMsg: &errno.ParamErr.ErrMsg})
		return
	}
	resp, err := rpc.ActionPublish(ctx, &publish.ActionRequest{
		Title: title,
		Token: token,
		Data:  buf.Bytes(),
	})
	if err != nil {
		Err := errno.ConvertErr(err)
		c.JSON(consts.StatusOK, &publish.ActionResponse{StatusCode: Err.ErrCode, StatusMsg: &Err.ErrMsg})
		return
	}
	SetResponse(c, resp)
	return
}

func MGetPublishList(ctx context.Context, c *app.RequestContext) {
	var req publish.PublishRequest
	err := c.BindAndValidate(&req)
	if err != nil {
		SetResponse(c, &publish.PublishResponse{StatusCode: errno.ParamErr.ErrCode, StatusMsg: &errno.ParamErr.ErrMsg})
		return
	}
	resp, err := rpc.MGetPublishList(ctx, &req)
	if err != nil {
		SetResponse(c, &publish.PublishResponse{StatusCode: errno.ServiceErr.ErrCode, StatusMsg: &errno.ServiceErr.ErrMsg})
		return
	}
	SetResponse(c, resp)
	return
}
