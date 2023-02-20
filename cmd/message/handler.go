package main

import (
	"context"

	"github.com/edufriendchen/light-tiktok/cmd/message/service"
	message "github.com/edufriendchen/light-tiktok/kitex_gen/message"
	"github.com/edufriendchen/light-tiktok/pkg/errno"
	"github.com/edufriendchen/light-tiktok/pkg/global"
)

// MessageServiceImpl implements the last service interface defined in the IDL.
type MessageServiceImpl struct{}

// ChatMessage implements the MessageServiceImpl interface.
func (s *MessageServiceImpl) ChatMessage(ctx context.Context, req *message.ChatRequest) (resp *message.ChatResponse, err error) {
	// TODO: Your code here...
	resp = new(message.ChatResponse)
	if err = req.IsValid(); err != nil {
		resp = &message.ChatResponse{StatusCode: errno.ParamErr.ErrCode, StatusMsg: &errno.ParamErr.ErrMsg}
		return resp, err
	}
	claims, err := global.Jwt.ParseToken(req.Token)
	if err != nil {
		resp = &message.ChatResponse{StatusCode: errno.AuthorizationFailedErr.ErrCode, StatusMsg: &errno.AuthorizationFailedErr.ErrMsg}
		return resp, nil
	}
	chatList, err := service.NewMGetChatMsgService(ctx).MGetChatMsg(req, claims.Id)
	if err != nil {
		resp = &message.ChatResponse{StatusCode: errno.ParamErr.ErrCode, StatusMsg: &errno.ParamErr.ErrMsg}
		return resp, err
	}
	resp = &message.ChatResponse{StatusCode: errno.Success.ErrCode, StatusMsg: &errno.Success.ErrMsg}
	resp.MessageList = chatList
	return resp, nil
}

// ActionMessage implements the MessageServiceImpl interface.
func (s *MessageServiceImpl) ActionMessage(ctx context.Context, req *message.ActionRequest) (resp *message.ActionResponse, err error) {
	// TODO: Your code here...
	resp = new(message.ActionResponse)
	if err = req.IsValid(); err != nil {
		resp = &message.ActionResponse{StatusCode: errno.ParamErr.ErrCode, StatusMsg: &errno.ParamErr.ErrMsg}
		return resp, err
	}
	claims, err := global.Jwt.ParseToken(req.Token)
	if err != nil {
		resp = &message.ActionResponse{StatusCode: errno.AuthorizationFailedErr.ErrCode, StatusMsg: &errno.AuthorizationFailedErr.ErrMsg}
		return resp, nil
	}
	err = service.NewActionMsgService(ctx).ActionMsg(req, claims.Id)
	if err != nil {
		resp = &message.ActionResponse{StatusCode: errno.ParamErr.ErrCode, StatusMsg: &errno.ParamErr.ErrMsg}
		return resp, err
	}
	resp = &message.ActionResponse{StatusCode: errno.Success.ErrCode, StatusMsg: &errno.Success.ErrMsg}
	return resp, nil
}
