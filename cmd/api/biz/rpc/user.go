package rpc

import (
	"context"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/edufriendchen/light-tiktok/kitex_gen/user"
	"github.com/edufriendchen/light-tiktok/kitex_gen/user/userservice"
	"github.com/edufriendchen/light-tiktok/pkg/consts"
	"github.com/edufriendchen/light-tiktok/pkg/errno"
	"github.com/edufriendchen/light-tiktok/pkg/global"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	"github.com/kitex-contrib/registry-nacos/resolver"
)

var userClient userservice.Client

func initUser() {
	c, err := userservice.NewClient(
		consts.USER_SERVICE_NAME,
		client.WithResolver(resolver.NewNacosResolver(global.NacosClient)),
		client.WithMuxConnection(1),
		client.WithSuite(tracing.NewClientSuite()),
		client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: consts.API_SERVICE_NAME}),
	)
	if err != nil {
		panic(err)
	}
	userClient = c
}

// CreateUser create user info
func CreateUser(ctx context.Context, req *user.CreateUserRequest) (*user.CreateUserResponse, error) {
	resp, err := userClient.CreateUser(ctx, req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 0 {
		return nil, errno.NewErrNo(resp.StatusCode, *resp.StatusMsg)
	}
	return resp, nil
}

// CheckUser check user info
func CheckUse(ctx context.Context, req *user.CheckUserRequest) (*user.CheckUserResponse, error) {
	resp, err := userClient.CheckUser(ctx, req)
	if err != nil {
		return resp, err
	}
	if resp.StatusCode != 0 {
		return resp, errno.NewErrNo(resp.StatusCode, *resp.StatusMsg)
	}
	return resp, nil
}

// Login rpc to user service GetUserInfo
func MGetUserInfo(ctx context.Context, req *user.MGetUserRequest) (*user.MGetUserResponse, error) {
	resp, err := userClient.MGetUser(ctx, req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 0 {
		return nil, errno.NewErrNo(resp.StatusCode, *resp.StatusMsg)
	}
	return resp, nil
}
