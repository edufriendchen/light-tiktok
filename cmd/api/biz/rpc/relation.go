package rpc

import (
	"context"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/edufriendchen/light-tiktok/kitex_gen/relation"
	"github.com/edufriendchen/light-tiktok/kitex_gen/relation/relationservice"
	"github.com/edufriendchen/light-tiktok/pkg/consts"
	"github.com/edufriendchen/light-tiktok/pkg/errno"
	"github.com/edufriendchen/light-tiktok/pkg/initialize"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	"github.com/kitex-contrib/registry-nacos/resolver"
)

var relationClient relationservice.Client

func initRelation() {
	cli, err := initialize.InitNacos()
	c, err := relationservice.NewClient(
		consts.RelationServiceName,
		client.WithResolver(resolver.NewNacosResolver(cli)),
		client.WithMuxConnection(1),
		client.WithSuite(tracing.NewClientSuite()),
		client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: consts.ApiServiceName}),
	)
	if err != nil {
		panic(err)
	}
	relationClient = c
}

// ActionRelation
func ActionRelation(ctx context.Context, req *relation.ActionRequest) (*relation.ActionResponse, error) {
	resp, err := relationClient.ActionRelation(ctx, req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 0 {
		return nil, errno.NewErrNo(resp.StatusCode, *resp.StatusMsg)
	}
	return resp, nil
}

// MGetFollowList
func MGetFollowList(ctx context.Context, req *relation.FollowRequest) (*relation.FollowResponse, error) {
	resp, err := relationClient.MGetFollowList(ctx, req)
	if err != nil {
		return resp, err
	}
	if resp.StatusCode != 0 {
		return resp, errno.NewErrNo(resp.StatusCode, *resp.StatusMsg)
	}
	return resp, nil
}

// MGetFollowerList
func MGetFollowerList(ctx context.Context, req *relation.FollowerRequest) (*relation.FollowerResponse, error) {
	resp, err := relationClient.MGetFollowerList(ctx, req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 0 {
		return nil, errno.NewErrNo(resp.StatusCode, *resp.StatusMsg)
	}
	return resp, nil
}

// MGetFriendList
func MGetFriendList(ctx context.Context, req *relation.FriendRequest) (*relation.FriendResponse, error) {
	resp, err := relationClient.MGetFriendList(ctx, req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 0 {
		return nil, errno.NewErrNo(resp.StatusCode, *resp.StatusMsg)
	}
	return resp, nil
}
