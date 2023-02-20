package pack

import (
	"context"

	"github.com/edufriendchen/light-tiktok/cmd/feed/dal"
	"github.com/edufriendchen/light-tiktok/kitex_gen/user"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// User pack user info
func User(ctx context.Context, session neo4j.SessionWithContext, u *user.User, fromID int64) (*user.User, error) {
	if u == nil {
		return &user.User{
			Name: "已注销用户",
		}, nil
	}
	isFollow := false
	isFollow, err := dal.HasFollow(ctx, session, fromID, u.Id)
	if err != nil {
		return u, err
	}
	u.IsFollow = isFollow
	return u, nil
}

// Users pack list of user info
func Users(ctx context.Context, session neo4j.SessionWithContext, us []*user.User, fromID int64) ([]*user.User, error) {
	users := make([]*user.User, 0)
	for _, u := range us {
		user2, err := User(ctx, session, u, fromID)
		if err != nil {
			return nil, err
		}
		if user2 != nil {
			users = append(users, user2)
		}
	}
	return users, nil
}
