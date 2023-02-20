package rpc

import (
	"context"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/edufriendchen/light-tiktok/kitex_gen/comment"
	"github.com/edufriendchen/light-tiktok/kitex_gen/comment/commentservice"
	"github.com/edufriendchen/light-tiktok/pkg/consts"
	"github.com/edufriendchen/light-tiktok/pkg/errno"
	"github.com/edufriendchen/light-tiktok/pkg/global"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	"github.com/kitex-contrib/registry-nacos/resolver"
)

var commentClient commentservice.Client

func initComment() {
	c, err := commentservice.NewClient(
		consts.COMMENT_SERVICE_NAME,
		client.WithResolver(resolver.NewNacosResolver(global.NacosClient)),
		client.WithMuxConnection(1),
		client.WithSuite(tracing.NewClientSuite()),
		client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: consts.API_SERVICE_NAME}),
	)
	if err != nil {
		panic(err)
	}
	commentClient = c
}

func ActionComment(ctx context.Context, req *comment.ActionRequest) (*comment.ActionResponse, error) {
	resp, err := commentClient.ActionComment(ctx, req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 0 {
		return nil, errno.NewErrNo(resp.StatusCode, *resp.StatusMsg)
	}
	return resp, nil
}

func MGetComment(ctx context.Context, req *comment.CommentRequest) (*comment.CommentResponse, error) {
	resp, err := commentClient.MGetCommentList(ctx, req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 0 {
		return nil, errno.NewErrNo(resp.StatusCode, *resp.StatusMsg)
	}
	return resp, nil
}
