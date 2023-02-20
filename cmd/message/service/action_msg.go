package service

import (
	"context"
	"fmt"

	"github.com/edufriendchen/light-tiktok/cmd/message/dal"
	"github.com/edufriendchen/light-tiktok/kitex_gen/message"
)

type ActionMsgService struct {
	ctx context.Context
}

func NewActionMsgService(ctx context.Context) *ActionMsgService {
	return &ActionMsgService{ctx: ctx}
}

func (s *ActionMsgService) ActionMsg(req *message.ActionRequest, user_id int64) error {
	fmt.Println("IN SERVICE")
	msglist := []*dal.Message{{
		ToUserId:   req.ToUserId,
		FromUserId: user_id,
		Content:    req.Content,
	}}
	err := dal.CreateMessage(s.ctx, msglist)
	if err != nil {
		return err
	}
	return nil
}
