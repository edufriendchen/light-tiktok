package service

import (
	"context"

	"github.com/edufriendchen/light-tiktok/cmd/user/dal"
	"github.com/edufriendchen/light-tiktok/kitex_gen/user"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type GetUserService struct {
	ctx     context.Context
	session neo4j.SessionWithContext
}

func NewGetUserService(ctx context.Context, driver neo4j.DriverWithContext) *GetUserService {
	session := driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)
	return &GetUserService{ctx: ctx, session: session}
}

func (s *GetUserService) GetUser(req *user.MGetUserRequest, user_id int64) (*user.User, error) {
	if req.UserId != user_id {
		user, err := dal.GetToUser(s.ctx, s.session, req, user_id)
		if err != nil {
			return nil, err
		}
		return user, err
	}
	user, err := dal.GetUser(s.ctx, s.session, req, user_id)
	if err != nil {
		return nil, err
	}
	return user, err
}

// func (s *GetUserService) GetUser(param *user.MGetUserRequest, user_id int64) (*user.User, error) {
// 	user, err := neo4j.ExecuteRead[*user.User](s.ctx, s.session, func(tx neo4j.ManagedTransaction) (*user.User, error) {
// 		result, err := tx.Run(s.ctx, "MATCH (n:User) WHERE id(n) = $toUserId MATCH (n1:User) WHERE id(n1) = $userId RETURN n, CASE WHEN (n1)-[:follow]->(n) THEN true ELSE false END AS result",
// 			map[string]any{
// 				"toUserId": param.UserId,
// 				"userId":   user_id,
// 			})
// 		if err != nil {
// 			return nil, err
// 		}
// 		record, err := result.Single(s.ctx)
// 		if err != nil {
// 			return nil, err
// 		}
// 		rawItemNode, found := record.Get("n")
// 		if !found {
// 			return nil, fmt.Errorf("could not find column")
// 		}
// 		itemNode := rawItemNode.(neo4j.Node)
// 		id := itemNode.GetId()
// 		if err != nil {
// 			return nil, err
// 		}
// 		name, err := neo4j.GetProperty[string](itemNode, "nickname")
// 		if err != nil {
// 			return nil, err
// 		}
// 		avatar, err := neo4j.GetProperty[string](itemNode, "avatar")
// 		if err != nil {
// 			return nil, err
// 		}
// 		follow_count, err := neo4j.GetProperty[int64](itemNode, "follow_count")
// 		if err != nil {
// 			return nil, err
// 		}
// 		follower_count, err := neo4j.GetProperty[int64](itemNode, "follower_count")
// 		if err != nil {
// 			return nil, err
// 		}

// 		IsFollow, found := record.Values[1].(bool)
// 		if err != nil {
// 			return nil, err
// 		}
// 		return &user.User{Id: id, Name: name, Avatar: avatar, FollowCount: follow_count, FollowerCount: follower_count, IsFollow: IsFollow}, nil
// 	})
// 	return user, err
// }
