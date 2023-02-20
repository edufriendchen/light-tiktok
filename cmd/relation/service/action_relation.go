package service

import (
	"context"
	"fmt"

	"github.com/edufriendchen/light-tiktok/kitex_gen/relation"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type ActionRelationService struct {
	ctx     context.Context
	session neo4j.SessionWithContext
}

func NewActionRelationService(ctx context.Context, driver neo4j.DriverWithContext) *ActionRelationService {
	session := driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)
	return &ActionRelationService{ctx: ctx, session: session}
}

func (s *ActionRelationService) ActionRelation(req *relation.ActionRequest, user_id int64) (bool, error) {
	is, err := isFollow(s.ctx, s.session, user_id, req.ToUserId)
	if err != nil {
		return false, err
	}
	if (req.ActionType == 1 && is) || (req.ActionType == 2 && !is) {
		return true, nil
	}
	if req.ActionType == 1 && !is {
		if _, err := doFollow(s.ctx, s.session, user_id, req.ToUserId); err != nil {
			return false, err
		}
	}
	if req.ActionType == 2 && is {
		if _, err := unFollow(s.ctx, s.session, user_id, req.ToUserId); err != nil {
			return false, err
		}
	}
	return true, nil
}

func isFollow(ctx context.Context, session neo4j.SessionWithContext, user_id int64, to_user_id int64) (bool, error) {
	is, err := neo4j.ExecuteRead[bool](ctx, session, func(tx neo4j.ManagedTransaction) (bool, error) {
		result, err := tx.Run(ctx, "MATCH (a:User) WHERE id(a) = $user_id MATCH (b:User) WHERE id(b) = $to_user_id MATCH (a)-[f:follow]->(b) WITH COUNT(f) > 0  as is_follow RETURN is_follow", map[string]any{
			"user_id":    user_id,
			"to_user_id": to_user_id,
		})
		if err != nil {
			return false, err
		}
		record, err := result.Single(ctx)
		if err != nil {
			return false, err
		}
		fmt.Println("isFollow:", record.Values)
		return record.Values[0].(bool), nil
	})
	return is, err
}

func doFollow(ctx context.Context, session neo4j.SessionWithContext, user_id int64, to_user_id int64) (int64, error) {
	relationship_id, err := neo4j.ExecuteWrite[int64](ctx, session, func(tx neo4j.ManagedTransaction) (int64, error) {
		result, err := tx.Run(ctx, "MATCH (a:User) WHERE id(a) = $user_id MATCH (b:User) WHERE id(b) = $to_user_id MERGE (a)-[f:follow]->(b) SET a.follow_count = a.follow_count + 1, b.follower_count = b.follower_count + 1 RETURN id(f) as id", map[string]any{
			"user_id":    user_id,
			"to_user_id": to_user_id,
		})
		if err != nil {
			return 0, err
		}
		record, err := result.Single(ctx)
		if err != nil {
			return 0, err
		}
		return record.Values[0].(int64), nil
	})
	return relationship_id, err
}

func unFollow(ctx context.Context, session neo4j.SessionWithContext, user_id int64, to_user_id int64) (int64, error) {
	relationship_id, err := neo4j.ExecuteWrite[int64](ctx, session, func(tx neo4j.ManagedTransaction) (int64, error) {
		result, err := tx.Run(ctx, "MATCH (a:User) WHERE id(a) = $user_id MATCH (b:User) WHERE id(b) = $to_user_id MATCH (a)-[f:follow]->(b)  DELETE f SET a.follow_count = a.follow_count - 1, b.follower_count = b.follower_count - 1 RETURN id(f)", map[string]any{
			"user_id":    user_id,
			"to_user_id": to_user_id,
		})
		if err != nil {
			return 0, err
		}
		record, err := result.Single(ctx)
		if err != nil {
			return 0, err
		}
		return record.Values[0].(int64), nil
	})
	return relationship_id, err
}
