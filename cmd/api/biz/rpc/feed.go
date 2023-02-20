package rpc

import (
	"context"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/edufriendchen/light-tiktok/kitex_gen/feed"
	"github.com/edufriendchen/light-tiktok/kitex_gen/feed/feedservice"
	"github.com/edufriendchen/light-tiktok/pkg/consts"
	"github.com/edufriendchen/light-tiktok/pkg/errno"
	"github.com/edufriendchen/light-tiktok/pkg/global"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	"github.com/kitex-contrib/registry-nacos/resolver"
)

var feedClient feedservice.Client

func initFeed() {
	c, err := feedservice.NewClient(
		consts.FEED_SERVICE_NAME,
		client.WithResolver(resolver.NewNacosResolver(global.NacosClient)),
		client.WithMuxConnection(1),
		client.WithSuite(tracing.NewClientSuite()),
		client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: consts.API_SERVICE_NAME}),
	)
	if err != nil {
		panic(err)
	}
	feedClient = c
}

// MGetFeedList
func MGetFeedList(ctx context.Context, req *feed.FeedRequest) (*feed.FeedResponse, error) {
	resp, err := feedClient.MGetFeedList(ctx, req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 0 {
		return nil, errno.NewErrNo(resp.StatusCode, *resp.StatusMsg)
	}
	return resp, nil
}
