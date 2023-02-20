package main

import (
	"context"
	"log"

	"github.com/edufriendchen/light-tiktok/cmd/comment/service"
	"github.com/edufriendchen/light-tiktok/kitex_gen/comment"
	"github.com/edufriendchen/light-tiktok/pkg/errno"
	"github.com/edufriendchen/light-tiktok/pkg/global"
)

// CommentServiceImpl implements the last service interface defined in the IDL.
type CommentServiceImpl struct{}

// ActionComment implements the CommentServiceImpl interface.
func (s *CommentServiceImpl) ActionComment(ctx context.Context, req *comment.ActionRequest) (resp *comment.ActionResponse, err error) {
	// TODO: Your code here...
	resp = new(comment.ActionResponse)
	if err = req.IsValid(); err != nil {
		resp = &comment.ActionResponse{StatusCode: errno.ParamErr.ErrCode, StatusMsg: &errno.ParamErr.ErrMsg}
		return resp, err
	}
	claims, err := global.Jwt.ParseToken(req.Token)
	if err != nil {
		resp = &comment.ActionResponse{StatusCode: errno.AuthorizationFailedErr.ErrCode, StatusMsg: &errno.AuthorizationFailedErr.ErrMsg}
		return resp, err
	}
	log.Println("claims:", claims.Id)
	comment_value, err := service.NewActionCommentService(ctx, global.Neo4jDriver).ActionComment(req, claims.Id, req.CommentId)
	if err != nil {
		resp = &comment.ActionResponse{StatusCode: errno.ServiceErr.ErrCode, StatusMsg: &errno.ServiceErr.ErrMsg}
		return resp, err
	}
	resp = &comment.ActionResponse{StatusCode: errno.Success.ErrCode, StatusMsg: &errno.Success.ErrMsg}
	resp.Comment = comment_value
	return resp, nil
}

// MGetCommentList implements the CommentServiceImpl interface.
func (s *CommentServiceImpl) MGetCommentList(ctx context.Context, req *comment.CommentRequest) (resp *comment.CommentResponse, err error) {
	// TODO: Your code here...
	resp = new(comment.CommentResponse)
	if err = req.IsValid(); err != nil {
		resp = &comment.CommentResponse{StatusCode: errno.ParamErr.ErrCode, StatusMsg: &errno.ParamErr.ErrMsg}
		return resp, err
	}
	commentList, err := service.NewMGetCommentService(ctx, global.Neo4jDriver).MGetCommentList(req)
	if err != nil {
		resp = &comment.CommentResponse{StatusCode: errno.AuthorizationFailedErr.ErrCode, StatusMsg: &errno.AuthorizationFailedErr.ErrMsg}
		return resp, err
	}
	resp = &comment.CommentResponse{StatusCode: errno.Success.ErrCode, StatusMsg: &errno.Success.ErrMsg}
	resp.CommentList = commentList
	return resp, nil
}
