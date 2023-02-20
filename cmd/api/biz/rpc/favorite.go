package rpc

import (
	"context"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/edufriendchen/light-tiktok/kitex_gen/favorite"
	"github.com/edufriendchen/light-tiktok/kitex_gen/favorite/favoriteservice"
	"github.com/edufriendchen/light-tiktok/pkg/consts"
	"github.com/edufriendchen/light-tiktok/pkg/errno"
	"github.com/edufriendchen/light-tiktok/pkg/initialize"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	"github.com/kitex-contrib/registry-nacos/resolver"
)

var favoriteClient favoriteservice.Client

func initFavorite() {
	cli, err := initialize.InitNacos()
	c, err := favoriteservice.NewClient(
		consts.FavoriteServiceName,
		client.WithResolver(resolver.NewNacosResolver(cli)),
		client.WithMuxConnection(1),
		client.WithSuite(tracing.NewClientSuite()),
		client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: consts.ApiServiceName}),
	)
	if err != nil {
		panic(err)
	}
	favoriteClient = c
}

// ActionFavorite
func ActionFavorite(ctx context.Context, req *favorite.ActionRequest) (*favorite.ActionResponse, error) {
	resp, err := favoriteClient.ActionFavorite(ctx, req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 0 {
		return nil, errno.NewErrNo(resp.StatusCode, *resp.StatusMsg)
	}
	return resp, nil
}

// MGetFavoriteList
func MGetFavoriteList(ctx context.Context, req *favorite.FavoriteRequest) (*favorite.FavoriteResponse, error) {
	resp, err := favoriteClient.MGetFavoriteList(ctx, req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 0 {
		return nil, errno.NewErrNo(resp.StatusCode, *resp.StatusMsg)
	}
	return resp, nil
}
