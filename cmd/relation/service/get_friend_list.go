package service

import (
	"context"
	"fmt"

	"github.com/edufriendchen/light-tiktok/kitex_gen/user"

	"github.com/edufriendchen/light-tiktok/kitex_gen/relation"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type GetFriendListService struct {
	ctx     context.Context
	session neo4j.SessionWithContext
}

func NewGetFriendListService(ctx context.Context, driver neo4j.DriverWithContext) *GetFriendListService {
	session := driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)
	return &GetFriendListService{ctx: ctx, session: session}
}

func (s *GetFriendListService) GetFriendList(req *relation.FriendRequest) ([]*user.User, error) {
	var list []*user.User
	userList, err := neo4j.ExecuteRead[[]*user.User](s.ctx, s.session, func(tx neo4j.ManagedTransaction) ([]*user.User, error) {
		result, err := tx.Run(s.ctx, "MATCH (n:User)-[:follow]->(to:User), (to:User)-[:follow]->(n:User)  WHERE id(n) = $userid  RETURN to",
			map[string]any{
				"userid": req.UserId,
			})
		if err != nil {
			return nil, err
		}
		record, err := result.Collect(s.ctx)
		if err != nil {
			return nil, err
		}
		for i := 0; i < len(record); i++ {
			value, ok := record[i].Get("v")
			fmt.Println("V:", value)
			if ok {
				itemNode := value.(neo4j.Node)
				id := itemNode.GetId()
				name, err := neo4j.GetProperty[string](itemNode, "u.nickname")
				if err != nil {
					return nil, err
				}
				follow_count, err := neo4j.GetProperty[int64](itemNode, "follow_count")
				if err != nil {
					return nil, err
				}
				follower_count, err := neo4j.GetProperty[int64](itemNode, "follower_count")
				if err != nil {
					return nil, err
				}
				list = append(list, &user.User{Id: id, Name: name, FollowCount: follow_count, FollowerCount: follower_count, IsFollow: true})
			}
		}
		return list, nil
	})
	if err != nil {
		return nil, err
	}
	return userList, nil
}
