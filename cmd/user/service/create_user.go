package service

import (
	"context"

	"github.com/edufriendchen/light-tiktok/cmd/user/dal"
	user "github.com/edufriendchen/light-tiktok/kitex_gen/user"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type CreateUserService struct {
	ctx     context.Context
	session neo4j.SessionWithContext
}

func NewCreateUserService(ctx context.Context, driver neo4j.DriverWithContext) *CreateUserService {
	session := driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)
	return &CreateUserService{ctx: ctx, session: session}
}

func (s *CreateUserService) CreateUserNode(req *user.CreateUserRequest) (int64, error) {
	user_id, err := dal.CreateUser(s.ctx, s.session, req)
	if err != nil {
		return 0, err
	}
	return user_id, err
}
