package service

import (
	"context"
	"strconv"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/edufriendchen/light-tiktok/cmd/user/dal"
	"github.com/edufriendchen/light-tiktok/kitex_gen/user"
	"github.com/edufriendchen/light-tiktok/pkg/cache"
	"github.com/edufriendchen/light-tiktok/pkg/consts"
	"github.com/edufriendchen/light-tiktok/pkg/global"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type GetUserService struct {
	ctx     context.Context
	session neo4j.SessionWithContext
	cache   cache.Client
}

func NewGetUserService(ctx context.Context, driver neo4j.DriverWithContext) *GetUserService {
	session := driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)
	return &GetUserService{ctx: ctx, session: session, cache: global.CacheManager}
}

func (s *GetUserService) GetUser(req *user.MGetUserRequest, user_id int64) (user *user.User, err error) {
	strconv.FormatInt(user_id, 10)
	err = s.cache.Get(s.ctx, consts.CACHE_USER+strconv.FormatInt(user_id, 10), &user)
	if err == nil {
		return user, nil
	}
	if req.UserId != user_id {
		user, err := dal.GetToUser(s.ctx, s.session, req, user_id)
		if err != nil {
			return nil, err
		}
		return user, nil
	}
	user, err = dal.GetUser(s.ctx, s.session, req, user_id)
	if err != nil {
		return nil, err
	}
	err = s.cache.Set(s.ctx, consts.CACHE_USER+strconv.FormatInt(user_id, 10), user)
	if err != nil {
		klog.Info("User info cache fail!")
	}
	return user, err
}
